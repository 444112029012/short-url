package application_test

import (
	"context"
	"testing"

	"github.com/444112029012/short-url/internal/application"
	"github.com/444112029012/short-url/internal/domain"
)

func TestRedirectSuccessIncrementsOnce(t *testing.T) {
	repo := newMem()
	repo.rows["Ab12Cd39"] = domain.URLRecord{ShortCode: "Ab12Cd39", LongURL: "https://example.com/a%2Fb?x=1&y=2", ClickCount: 3}
	svc := application.NewRedirectService(domain.NewURLValidator(), repo, repo)
	target, err := svc.Redirect(context.Background(), "Ab12Cd39")
	if err != nil {
		t.Fatal(err)
	}
	if target.LongURL != "https://example.com/a%2Fb?x=1&y=2" {
		t.Fatalf("UT-RED-01 location %s", target.LongURL)
	}
	if repo.rows["Ab12Cd39"].ClickCount != 4 {
		t.Fatalf("count %d", repo.rows["Ab12Cd39"].ClickCount)
	}
}

func TestRedirectMissingDoesNotIncrement(t *testing.T) {
	repo := newMem()
	repo.rows["Ab12Cd39"] = domain.URLRecord{ShortCode: "Ab12Cd39", LongURL: "https://example.com", ClickCount: 2}
	svc := application.NewRedirectService(domain.NewURLValidator(), repo, repo)
	_, err := svc.Redirect(context.Background(), "Zz98Yy76")
	if err == nil || err.Type != domain.TypeNotFound {
		t.Fatalf("UT-RED-02: %v", err)
	}
	if repo.rows["Ab12Cd39"].ClickCount != 2 {
		t.Fatalf("count changed: %d", repo.rows["Ab12Cd39"].ClickCount)
	}
	_, err = svc.Redirect(context.Background(), "bad code")
	if err == nil || err.Type != domain.TypeInvalidURL {
		t.Fatalf("format: %v", err)
	}
	if repo.rows["Ab12Cd39"].ClickCount != 2 {
		t.Fatal("invalid code incremented")
	}
}

func TestRedirectSkipsUnsafeStoredURL(t *testing.T) {
	repo := newMem()
	repo.rows["Ab12Cd39"] = domain.URLRecord{ShortCode: "Ab12Cd39", LongURL: "https://example.com/\r\nX: y", ClickCount: 1}
	svc := application.NewRedirectService(domain.NewURLValidator(), repo, repo)
	_, err := svc.Redirect(context.Background(), "Ab12Cd39")
	if err == nil || err.Type != domain.TypeInternal {
		t.Fatalf("%v", err)
	}
	if repo.rows["Ab12Cd39"].ClickCount != 1 {
		t.Fatalf("count changed: %d", repo.rows["Ab12Cd39"].ClickCount)
	}
}
