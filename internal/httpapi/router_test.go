package httpapi_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/444112029012/short-url/internal/application"
	"github.com/444112029012/short-url/internal/domain"
	"github.com/444112029012/short-url/internal/httpapi"
	"github.com/444112029012/short-url/internal/infra/sqlite"
	"github.com/444112029012/short-url/internal/observability"
	"github.com/444112029012/short-url/internal/ratelimit"
)

func newServer(t *testing.T) http.Handler {
	t.Helper()
	store, err := sqlite.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if err := store.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}
	validator := domain.NewURLValidator()
	appSvc := application.NewShortURLApplicationService(validator, domain.NewShortCodeGenerator(), store, store, "http://localhost:8080")
	redirectSvc := application.NewRedirectService(validator, store, store)
	adapter := httpapi.NewHTTPAPIAdapter(appSvc, redirectSvc, ratelimit.NewGuard(30, 120), httpapi.NewErrorMapper(), observability.NewHooks(), nil)
	return httpapi.NewRouter(adapter)
}

func TestHealth(t *testing.T) {
	srv := newServer(t)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	srv.ServeHTTP(res, req)
	if res.Code != http.StatusOK || res.Body.String() != "ok" {
		t.Fatalf("health %d %q", res.Code, res.Body.String())
	}
}

func TestCreateRedirectStatsHappyPath(t *testing.T) {
	srv := newServer(t)
	longURL := "https://example.com/path?q=1"
	res := postURL(t, srv, longURL)
	if res.Code != http.StatusCreated {
		t.Fatalf("UT-HC-01 status %d body %s", res.Code, res.Body.String())
	}
	created := decodeMap(t, res.Body.Bytes())
	code, _ := created["short_code"].(string)
	shortURL, _ := created["short_url"].(string)
	if created["long_url"] != longURL || shortURL != "http://localhost:8080/"+code {
		t.Fatalf("create body %s", res.Body.String())
	}
	if len(code) != 8 {
		t.Fatalf("code %q", code)
	}
	if len(created) != 3 {
		t.Fatalf("extra create fields %v", created)
	}

	stats := getJSON(t, srv, "/api/v1/urls/"+code+"/stats")
	if stats.Code != http.StatusOK || decodeMap(t, stats.Body.Bytes())["click_count"] != float64(0) {
		t.Fatalf("UT-HS-01 initial %d %s", stats.Code, stats.Body.String())
	}

	redirect := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/"+code+"?url=https://evil.example&long_url=https://evil.example", nil)
	req.Header.Set("X-Redirect-Target", "https://evil.example")
	srv.ServeHTTP(redirect, req)
	if redirect.Code != http.StatusFound {
		t.Fatalf("UT-HR-01 status %d", redirect.Code)
	}
	if got := redirect.Header().Get("Location"); got != longURL {
		t.Fatalf("SEC-014 Location %q", got)
	}

	stats = getJSON(t, srv, "/api/v1/urls/"+code+"/stats")
	body := decodeMap(t, stats.Body.Bytes())
	if stats.Code != http.StatusOK || body["click_count"] != float64(1) || body["short_code"] != code || len(body) != 2 {
		t.Fatalf("stats after redirect %d %s", stats.Code, stats.Body.String())
	}

	second := postURL(t, srv, longURL)
	other := decodeMap(t, second.Body.Bytes())
	if second.Code != http.StatusCreated || other["short_code"] == code || other["long_url"] != longURL {
		t.Fatalf("REQ-012 second create %d %s", second.Code, second.Body.String())
	}
}

func TestCreateValidationError(t *testing.T) {
	srv := newServer(t)
	for _, raw := range []string{"", "ftp://example.com/a", "javascript:alert(1)", "not-a-url"} {
		res := postURL(t, srv, raw)
		assertError(t, res, http.StatusBadRequest, "invalid_url")
	}
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/urls", strings.NewReader(`{"url":"https://example.com","extra":1}`))
	req.Header.Set("Content-Type", "application/json")
	srv.ServeHTTP(res, req)
	assertError(t, res, http.StatusBadRequest, "invalid_url")

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/v1/urls", strings.NewReader(`{"url":"https://example.com"}`))
	req.Header.Set("Content-Type", "text/plain")
	srv.ServeHTTP(res, req)
	assertError(t, res, http.StatusBadRequest, "invalid_url")
}

func TestRedirectAndStatsErrors(t *testing.T) {
	srv := newServer(t)
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/NoSuch01", nil))
	assertError(t, res, http.StatusNotFound, "not_found")
	if res.Header().Get("Location") != "" {
		t.Fatal("not_found set Location")
	}

	res = httptest.NewRecorder()
	srv.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/bad", nil))
	assertError(t, res, http.StatusBadRequest, "invalid_url")

	res = httptest.NewRecorder()
	srv.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/v1/urls/badcode!/stats", nil))
	assertError(t, res, http.StatusBadRequest, "invalid_url")

	res = httptest.NewRecorder()
	srv.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/v1/urls/NoSuch01/stats", nil))
	assertError(t, res, http.StatusNotFound, "not_found")
}

type blockingLimiter struct{}

func (blockingLimiter) CheckCreate(string) *domain.AppError {
	return domain.RateLimited(domain.MsgTooManyRequests)
}
func (blockingLimiter) CheckRedirect(string) *domain.AppError {
	return domain.RateLimited(domain.MsgTooManyRequests)
}

func TestAdapterHonorsLimiterStubContract(t *testing.T) {
	path := filepath.Join(t.TempDir(), "limited.db")
	store, err := sqlite.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if err := store.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}
	validator := domain.NewURLValidator()
	appSvc := application.NewShortURLApplicationService(validator, domain.NewShortCodeGenerator(), store, store, "http://localhost:8080")
	redirectSvc := application.NewRedirectService(validator, store, store)
	adapter := httpapi.NewHTTPAPIAdapter(appSvc, redirectSvc, blockingLimiter{}, httpapi.NewErrorMapper(), observability.NewHooks(), nil)
	srv := httpapi.NewRouter(adapter)

	res := postURL(t, srv, "https://example.com")
	assertError(t, res, http.StatusTooManyRequests, "rate_limited")
	if n := countRows(t, path); n != 0 {
		t.Fatalf("limiter path stored %d rows", n)
	}

	redirect := httptest.NewRecorder()
	srv.ServeHTTP(redirect, httptest.NewRequest(http.MethodGet, "/Ab12Cd39", nil))
	assertError(t, redirect, http.StatusTooManyRequests, "rate_limited")
	if redirect.Header().Get("Location") != "" {
		t.Fatal("limited redirect set Location")
	}
}

func TestClosedStoreReturnsInternalError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "closed.db")
	store, err := sqlite.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	validator := domain.NewURLValidator()
	appSvc := application.NewShortURLApplicationService(validator, domain.NewShortCodeGenerator(), store, store, "http://localhost:8080")
	redirectSvc := application.NewRedirectService(validator, store, store)
	adapter := httpapi.NewHTTPAPIAdapter(appSvc, redirectSvc, ratelimit.NewGuard(30, 120), httpapi.NewErrorMapper(), observability.NewHooks(), nil)
	res := postURL(t, httpapi.NewRouter(adapter), "https://example.com/closed")
	assertError(t, res, http.StatusInternalServerError, "internal_error")
}

func countRows(t *testing.T, path string) int {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM url_mapping`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func postURL(t *testing.T, srv http.Handler, rawURL string) *httptest.ResponseRecorder {
	t.Helper()
	payload, err := json.Marshal(map[string]string{"url": rawURL})
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/urls", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	srv.ServeHTTP(res, req)
	return res
}

func getJSON(t *testing.T, srv http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, httptest.NewRequest(http.MethodGet, path, nil))
	return res
}

func decodeMap(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("json %s: %s", err, raw)
	}
	return out
}

func assertError(t *testing.T, res *httptest.ResponseRecorder, status int, typ string) {
	t.Helper()
	if res.Code != status {
		t.Fatalf("status %d body %s", res.Code, res.Body.String())
	}
	if ct := res.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("content-type %q", ct)
	}
	body := decodeMap(t, res.Body.Bytes())
	errObj, _ := body["error"].(map[string]any)
	if errObj["type"] != typ {
		t.Fatalf("type %v body %s", errObj["type"], res.Body.String())
	}
	msg, _ := errObj["message"].(string)
	for _, forbidden := range []string{"SELECT", ".go:", "sql:", "/tmp/", "panic"} {
		if strings.Contains(res.Body.String(), forbidden) {
			t.Fatalf("leaked %s in %s", forbidden, res.Body.String())
		}
	}
	if msg == "" {
		t.Fatal("empty message")
	}
}
