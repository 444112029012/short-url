package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/444112029012/short-url/internal/application"
	"github.com/444112029012/short-url/internal/domain"
	"github.com/444112029012/short-url/internal/ports"
)

type seqGen struct {
	codes []string
	i     int
	err   *domain.AppError
}

func (g *seqGen) Generate() (string, *domain.AppError) {
	if g.err != nil {
		return "", g.err
	}
	if g.i >= len(g.codes) {
		return "", domain.Internal()
	}
	code := g.codes[g.i]
	g.i++
	return code, nil
}

type memRepo struct {
	rows      map[string]domain.URLRecord
	saves     int
	conflicts int
	failSave  error
}

func newMem() *memRepo {
	return &memRepo{rows: map[string]domain.URLRecord{}}
}

func (m *memRepo) Save(_ context.Context, record domain.URLRecord) error {
	m.saves++
	if m.failSave != nil {
		return m.failSave
	}
	if m.conflicts > 0 {
		m.conflicts--
		return ports.ErrConflict
	}
	if _, ok := m.rows[record.ShortCode]; ok {
		return ports.ErrConflict
	}
	m.rows[record.ShortCode] = record
	return nil
}

func (m *memRepo) FindByShortCode(_ context.Context, code string) (*domain.URLRecord, error) {
	rec, ok := m.rows[code]
	if !ok {
		return nil, ports.ErrNotFound
	}
	copy := rec
	return &copy, nil
}

func (m *memRepo) Exists(_ context.Context, code string) (bool, error) {
	_, ok := m.rows[code]
	return ok, nil
}

func (m *memRepo) IncrementAtomic(_ context.Context, shortCode string) (int, error) {
	rec, ok := m.rows[shortCode]
	if !ok {
		return 0, ports.ErrNotFound
	}
	rec.ClickCount++
	m.rows[shortCode] = rec
	return rec.ClickCount, nil
}

func (m *memRepo) GetCount(_ context.Context, shortCode string) (int, error) {
	rec, ok := m.rows[shortCode]
	if !ok {
		return 0, ports.ErrNotFound
	}
	return rec.ClickCount, nil
}

func TestCreateShortURLSuccess(t *testing.T) {
	repo := newMem()
	svc := application.NewShortURLApplicationService(
		domain.NewURLValidator(),
		&seqGen{codes: []string{"Ab12Cd39"}},
		repo,
		repo,
		"http://localhost:8080",
	)
	created, err := svc.CreateShortURL(context.Background(), "https://example.com/path?q=1")
	if err != nil {
		t.Fatal(err)
	}
	if created.ShortCode != "Ab12Cd39" || created.LongURL != "https://example.com/path?q=1" {
		t.Fatalf("UT-APP-01: %+v", created)
	}
	if created.ShortURL != "http://localhost:8080/Ab12Cd39" {
		t.Fatalf("short url %s", created.ShortURL)
	}
	if repo.rows["Ab12Cd39"].ClickCount != 0 || repo.saves != 1 {
		t.Fatalf("row %+v saves %d", repo.rows["Ab12Cd39"], repo.saves)
	}
}

func TestCreateShortURLInvalidDoesNotSave(t *testing.T) {
	repo := newMem()
	svc := application.NewShortURLApplicationService(
		domain.NewURLValidator(),
		&seqGen{codes: []string{"Ab12Cd39"}},
		repo,
		repo,
		"http://localhost:8080",
	)
	_, err := svc.CreateShortURL(context.Background(), "ftp://example.com")
	if err == nil || err.Type != domain.TypeInvalidURL {
		t.Fatalf("UT-APP-02: %v", err)
	}
	if repo.saves != 0 || len(repo.rows) != 0 {
		t.Fatalf("saved on failure: %+v", repo.rows)
	}
}

func TestCreateShortURLRetriesConflict(t *testing.T) {
	repo := newMem()
	repo.conflicts = 1
	svc := application.NewShortURLApplicationService(
		domain.NewURLValidator(),
		&seqGen{codes: []string{"Ab12Cd39", "Zz98Yy76"}},
		repo,
		repo,
		"http://localhost:8080/",
	)
	created, err := svc.CreateShortURL(context.Background(), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if created.ShortCode != "Zz98Yy76" || created.ShortURL != "http://localhost:8080/Zz98Yy76" {
		t.Fatalf("UT-APP-03: %+v", created)
	}
}

func TestCreateShortURLAllowsDuplicateLongURL(t *testing.T) {
	repo := newMem()
	svc := application.NewShortURLApplicationService(
		domain.NewURLValidator(),
		&seqGen{codes: []string{"Ab12Cd39", "Zz98Yy76"}},
		repo,
		repo,
		"http://localhost:8080",
	)
	first, err := svc.CreateShortURL(context.Background(), "https://example.com/same")
	if err != nil {
		t.Fatal(err)
	}
	second, err := svc.CreateShortURL(context.Background(), "https://example.com/same")
	if err != nil {
		t.Fatal(err)
	}
	if first.ShortCode == second.ShortCode || first.LongURL != second.LongURL {
		t.Fatalf("REQ-012: %+v %+v", first, second)
	}
}

func TestCreateShortURLStoreFailureIsInternal(t *testing.T) {
	repo := newMem()
	repo.failSave = errors.New("SQL syntax near SELECT /var/lib/shorturl.db")
	svc := application.NewShortURLApplicationService(
		domain.NewURLValidator(),
		&seqGen{codes: []string{"Ab12Cd39"}},
		repo,
		repo,
		"http://localhost:8080",
	)
	_, err := svc.CreateShortURL(context.Background(), "https://example.com")
	if err == nil || err.Type != domain.TypeInternal || err.Message != domain.MsgUnexpected {
		t.Fatalf("%v", err)
	}
	if err.Error() != "internal_error" {
		t.Fatalf("leaked %q", err.Error())
	}
}

func TestGetStats(t *testing.T) {
	repo := newMem()
	repo.rows["Ab12Cd39"] = domain.URLRecord{ShortCode: "Ab12Cd39", LongURL: "https://example.com", ClickCount: 4}
	svc := application.NewShortURLApplicationService(domain.NewURLValidator(), &seqGen{}, repo, repo, "http://localhost:8080")
	stats, err := svc.GetStats(context.Background(), "Ab12Cd39")
	if err != nil || stats.ClickCount != 4 || stats.ShortCode != "Ab12Cd39" {
		t.Fatalf("UT-APP-04: %+v %v", stats, err)
	}
	_, err = svc.GetStats(context.Background(), "Zz98Yy76")
	if err == nil || err.Type != domain.TypeNotFound {
		t.Fatalf("UT-APP-05: %v", err)
	}
	_, err = svc.GetStats(context.Background(), "bad")
	if err == nil || err.Type != domain.TypeInvalidURL {
		t.Fatalf("invalid code: %v", err)
	}
}
