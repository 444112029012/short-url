package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/444112029012/short-url/internal/application"
	"github.com/444112029012/short-url/internal/domain"
)

// countingGen records every Generate call, including failures.
type countingGen struct {
	codes []string
	n     int
	err   *domain.AppError
}

func (g *countingGen) Generate() (string, *domain.AppError) {
	g.n++
	if g.err != nil {
		return "", g.err
	}
	if g.n > len(g.codes) {
		return "", domain.Internal()
	}
	return g.codes[g.n-1], nil
}

func TestCreateShortURLExhaustsMixedCollisions(t *testing.T) {
	codes := []string{
		"AAAAAAAA", "BBBBBBBB", "CCCCCCCC", "DDDDDDDD",
		"EEEEEEEE", "FFFFFFFF", "GGGGGGGG", "HHHHHHHH",
	}
	repo := newMem()
	for _, code := range codes[:4] {
		repo.rows[code] = domain.URLRecord{ShortCode: code, LongURL: "https://example.com/existing"}
	}
	repo.conflicts = 4
	gen := &countingGen{codes: codes}
	svc := application.NewShortURLApplicationService(
		domain.NewURLValidator(),
		gen,
		repo,
		repo,
		"http://localhost:8080",
	)
	_, err := svc.CreateShortURL(context.Background(), "https://example.com/new")
	if err == nil || err.Type != domain.TypeInternal || err.Message != domain.MsgUnexpected {
		t.Fatalf("got %v", err)
	}
	if gen.n != 8 {
		t.Fatalf("generate attempts %d, want 8", gen.n)
	}
	if repo.saves != 4 {
		t.Fatalf("saves %d, want 4 conflict writes", repo.saves)
	}
	for _, code := range codes[4:] {
		if _, ok := repo.rows[code]; ok {
			t.Fatalf("conflict stored %s", code)
		}
	}
	if len(repo.rows) != 4 {
		t.Fatalf("rows %d", len(repo.rows))
	}
}

func TestCreateShortURLPropagatesGeneratorFailure(t *testing.T) {
	cause := errors.New("csprng failed")
	gen := &countingGen{err: domain.InternalWrap(cause)}
	repo := newMem()
	svc := application.NewShortURLApplicationService(
		domain.NewURLValidator(),
		gen,
		repo,
		repo,
		"http://localhost:8080",
	)
	_, err := svc.CreateShortURL(context.Background(), "https://example.com/gen")
	if err == nil || err.Type != domain.TypeInternal || err.Message != domain.MsgUnexpected {
		t.Fatalf("got %v", err)
	}
	if !errors.Is(err, cause) {
		t.Fatal("generator failure was not propagated")
	}
	if gen.n != 1 || repo.saves != 0 || len(repo.rows) != 0 {
		t.Fatalf("attempts %d saves %d rows %d", gen.n, repo.saves, len(repo.rows))
	}
}
