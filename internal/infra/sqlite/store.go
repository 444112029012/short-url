package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/444112029012/short-url/internal/domain"
	"github.com/444112029012/short-url/internal/ports"
	sqlitedriver "modernc.org/sqlite"
)

// URLStore is the SQLite realization of InMemoryOrDbUrlStore (ADR-002).
// Every statement below is a fixed string with bound parameters (SEC-016).
type URLStore struct {
	db *sql.DB
}

var (
	_ ports.URLRepository = (*URLStore)(nil)
	_ ports.ClickCounter  = (*URLStore)(nil)
)

const createTableSQL = `
CREATE TABLE IF NOT EXISTS url_mapping (
    short_code TEXT PRIMARY KEY,
    long_url TEXT NOT NULL,
    click_count INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    CHECK (length(long_url) <= 2048),
    CHECK (click_count >= 0)
)`

const insertSQL = `
INSERT INTO url_mapping (short_code, long_url, click_count, created_at)
VALUES (?, ?, ?, ?)`

const selectByCodeSQL = `
SELECT short_code, long_url, click_count, created_at
FROM url_mapping
WHERE short_code = ?`

const existsSQL = `
SELECT 1 FROM url_mapping WHERE short_code = ? LIMIT 1`

// Single-statement atomic increment (SEC-016). Not a read-modify-write.
const incrementSQL = `
UPDATE url_mapping
SET click_count = click_count + 1
WHERE short_code = ?
RETURNING click_count`

const countSQL = `
SELECT click_count FROM url_mapping WHERE short_code = ?`

// Open opens a SQLite file (or :memory:) and applies busy_timeout.
// The file mode is 0600. The database is not a network listener (SEC-017).
func Open(path string) (*URLStore, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("database path is empty")
	}
	dsn := path
	if path != ":memory:" {
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		if err := os.MkdirAll(filepath.Dir(abs), 0o700); err != nil {
			return nil, err
		}
		dsn = abs
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err := db.Exec(`PRAGMA busy_timeout = 5000`); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		_ = db.Close()
		return nil, err
	}
	if path != ":memory:" {
		if err := os.Chmod(dsn, 0o600); err != nil {
			_ = db.Close()
			return nil, err
		}
	}
	return &URLStore{db: db}, nil
}

func (s *URLStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

// Migrate creates url_mapping when missing. The statement has no user input.
// Columns are only short_code, long_url, click_count, created_at (SEC-004).
func (s *URLStore) Migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, createTableSQL)
	return err
}

func (s *URLStore) Save(ctx context.Context, record domain.URLRecord) error {
	if record.ClickCount < 0 {
		return errors.New("click_count must be >= 0")
	}
	created := record.CreatedAt.UTC().Format(time.RFC3339Nano)
	_, err := s.db.ExecContext(ctx, insertSQL, record.ShortCode, record.LongURL, record.ClickCount, created)
	if err != nil {
		if isUniqueViolation(err) {
			return ports.ErrConflict
		}
		return err
	}
	return nil
}

func (s *URLStore) FindByShortCode(ctx context.Context, code string) (*domain.URLRecord, error) {
	var rec domain.URLRecord
	var created string
	err := s.db.QueryRowContext(ctx, selectByCodeSQL, code).Scan(
		&rec.ShortCode, &rec.LongURL, &rec.ClickCount, &created,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ports.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	parsed, err := parseTimestamp(created)
	if err != nil {
		return nil, err
	}
	rec.CreatedAt = parsed
	return &rec, nil
}

func (s *URLStore) Exists(ctx context.Context, code string) (bool, error) {
	var one int
	err := s.db.QueryRowContext(ctx, existsSQL, code).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *URLStore) IncrementAtomic(ctx context.Context, shortCode string) (int, error) {
	var next int64
	err := s.db.QueryRowContext(ctx, incrementSQL, shortCode).Scan(&next)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ports.ErrNotFound
	}
	if err != nil {
		return 0, err
	}
	if next < 0 {
		return 0, errors.New("click_count became negative")
	}
	return int(next), nil
}

func (s *URLStore) GetCount(ctx context.Context, shortCode string) (int, error) {
	var n int64
	err := s.db.QueryRowContext(ctx, countSQL, shortCode).Scan(&n)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ports.ErrNotFound
	}
	if err != nil {
		return 0, err
	}
	if n < 0 {
		return 0, errors.New("click_count is negative")
	}
	return int(n), nil
}

func parseTimestamp(value string) (time.Time, error) {
	if ts, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return ts.UTC(), nil
	}
	ts, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, err
	}
	return ts.UTC(), nil
}

func isUniqueViolation(err error) bool {
	var se *sqlitedriver.Error
	if errors.As(err, &se) {
		code := se.Code()
		// SQLITE_CONSTRAINT=19, SQLITE_CONSTRAINT_PRIMARYKEY=1555, SQLITE_CONSTRAINT_UNIQUE=2067.
		if code == 1555 || code == 2067 {
			return true
		}
	}
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint failed") || strings.Contains(msg, "PRIMARY KEY constraint failed")
}
