package sqlite

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/444112029012/short-url/internal/domain"
	"github.com/444112029012/short-url/internal/ports"
)

func openTestStore(t *testing.T) *URLStore {
	t.Helper()
	path := filepath.Join(t.TempDir(), "shorturl.db")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if err := store.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := store.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}
	return store
}

func sample(code, longURL string) domain.URLRecord {
	return domain.URLRecord{
		ShortCode:  code,
		LongURL:    longURL,
		ClickCount: 0,
		CreatedAt:  time.Date(2026, 10, 7, 1, 2, 3, 4, time.UTC),
	}
}

func TestSaveFindExistsAndSchema(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t)
	rec := sample("Ab12Cd39", "https://example.com/a' OR '1'='1")
	if err := store.Save(ctx, rec); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(ctx, sample("Zz98Yy76", "https://example.com/other")); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(ctx, rec); !errors.Is(err, ports.ErrConflict) {
		t.Fatalf("UT-REPO-01 duplicate: %v", err)
	}

	got, err := store.FindByShortCode(ctx, rec.ShortCode)
	if err != nil {
		t.Fatal(err)
	}
	if got.LongURL != rec.LongURL || got.ShortCode != rec.ShortCode || got.ClickCount != 0 {
		t.Fatalf("UT-REPO-02: %+v", got)
	}
	if !got.CreatedAt.Equal(rec.CreatedAt) {
		t.Fatalf("created_at: %s", got.CreatedAt)
	}

	exists, err := store.Exists(ctx, rec.ShortCode)
	if err != nil || !exists {
		t.Fatalf("UT-REPO-03 exists: %v %v", exists, err)
	}
	exists, err = store.Exists(ctx, "NoSuch01")
	if err != nil || exists {
		t.Fatalf("UT-REPO-03 missing: %v %v", exists, err)
	}
	if _, err := store.FindByShortCode(ctx, "NoSuch01"); !errors.Is(err, ports.ErrNotFound) {
		t.Fatalf("find missing: %v", err)
	}

	rows, err := store.db.Query(`SELECT name FROM pragma_table_info('url_mapping')`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	gotCols := map[string]bool{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		gotCols[name] = true
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"visitor_ip", "user_agent"} {
		if gotCols[forbidden] {
			t.Fatalf("SEC-004 column %s present", forbidden)
		}
	}
	for _, required := range []string{"short_code", "long_url", "click_count", "created_at"} {
		if !gotCols[required] {
			t.Fatalf("missing column %s", required)
		}
	}

	var n int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM url_mapping`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("IT-STORE-03 row count %d", n)
	}
	if strings.Contains(got.LongURL, "1=1") && got.ShortCode != rec.ShortCode {
		t.Fatal("lookup returned the wrong row")
	}
}

func TestIncrementAndGetCount(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t)
	code := "Ab12Cd39"
	if err := store.Save(ctx, sample(code, "https://example.com/path")); err != nil {
		t.Fatal(err)
	}
	if _, err := store.IncrementAtomic(ctx, "NoSuch01"); !errors.Is(err, ports.ErrNotFound) {
		t.Fatalf("increment missing: %v", err)
	}
	if _, err := store.GetCount(ctx, "NoSuch01"); !errors.Is(err, ports.ErrNotFound) {
		t.Fatalf("count missing: %v", err)
	}
	next, err := store.IncrementAtomic(ctx, code)
	if err != nil || next != 1 {
		t.Fatalf("UT-CLK-01: %d %v", next, err)
	}
	got, err := store.GetCount(ctx, code)
	if err != nil || got != 1 {
		t.Fatalf("UT-CLK-02: %d %v", got, err)
	}
}

func TestIncrementAtomicConcurrent(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t)
	code := "Ab12Cd39"
	if err := store.Save(ctx, sample(code, "https://example.com/path")); err != nil {
		t.Fatal(err)
	}
	const n = 40
	var wg sync.WaitGroup
	errCh := make(chan error, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if _, err := store.IncrementAtomic(ctx, code); err != nil {
				errCh <- err
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Error(err)
	}
	got, err := store.GetCount(ctx, code)
	if err != nil {
		t.Fatal(err)
	}
	if got != n {
		t.Fatalf("IT-STORE-02: got %d want %d", got, n)
	}
}

func TestFileMode(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mode.db")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("perm %o", info.Mode().Perm())
	}
}
