package httpapi_test

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/444112029012/short-url/internal/application"
	"github.com/444112029012/short-url/internal/domain"
	"github.com/444112029012/short-url/internal/httpapi"
	"github.com/444112029012/short-url/internal/infra/sqlite"
	"github.com/444112029012/short-url/internal/observability"
	"github.com/444112029012/short-url/internal/ratelimit"
)

func TestSecurityHeadersOnSuccessRedirectAndErrors(t *testing.T) {
	srv := newServer(t)
	created := postURL(t, srv, "https://example.com/headers")
	if created.Code != http.StatusCreated {
		t.Fatalf("201 %d %s", created.Code, created.Body.String())
	}
	assertBaselineHeaders(t, created)
	assertJSONNoStore(t, created)
	if created.Header().Get("Strict-Transport-Security") != "" {
		t.Fatal("local HTTP emitted HSTS")
	}
	code, _ := decodeMap(t, created.Body.Bytes())["short_code"].(string)

	redirect := httptest.NewRecorder()
	srv.ServeHTTP(redirect, httptest.NewRequest(http.MethodGet, "/"+code, nil))
	if redirect.Code != http.StatusFound {
		t.Fatalf("302 %d", redirect.Code)
	}
	if redirect.Header().Get("Location") != "https://example.com/headers" {
		t.Fatalf("Location %q", redirect.Header().Get("Location"))
	}
	assertBaselineHeaders(t, redirect)
	if redirect.Header().Get("Cache-Control") != "" {
		t.Fatalf("302 Cache-Control %q", redirect.Header().Get("Cache-Control"))
	}

	bad := postURL(t, srv, "ftp://example.com/a")
	assertError(t, bad, http.StatusBadRequest, "invalid_url")
	assertBaselineHeaders(t, bad)
	assertJSONNoStore(t, bad)

	missing := httptest.NewRecorder()
	srv.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/NoSuch01", nil))
	assertError(t, missing, http.StatusNotFound, "not_found")
	assertBaselineHeaders(t, missing)
	assertJSONNoStore(t, missing)
	if missing.Header().Get("Location") != "" {
		t.Fatal("404 set Location")
	}

	unknown := httptest.NewRecorder()
	srv.ServeHTTP(unknown, httptest.NewRequest(http.MethodGet, "/api/v1/does-not-exist", nil))
	if unknown.Code != http.StatusNotFound {
		t.Fatalf("chi 404 %d", unknown.Code)
	}
	assertBaselineHeaders(t, unknown)

	method := httptest.NewRecorder()
	srv.ServeHTTP(method, httptest.NewRequest(http.MethodPost, "/health", nil))
	if method.Code != http.StatusMethodNotAllowed {
		t.Fatalf("405 %d", method.Code)
	}
	assertBaselineHeaders(t, method)
}

func TestSecurityHeadersOn429And500(t *testing.T) {
	path := filepath.Join(t.TempDir(), "headers.db")
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
	limited := httpapi.NewRouter(httpapi.NewHTTPAPIAdapter(appSvc, redirectSvc, blockingLimiter{}, httpapi.NewErrorMapper(), observability.NewHooks()))
	tooMany := postURL(t, limited, "https://example.com/limited")
	assertError(t, tooMany, http.StatusTooManyRequests, "rate_limited")
	assertBaselineHeaders(t, tooMany)
	assertJSONNoStore(t, tooMany)

	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	closed := httpapi.NewRouter(httpapi.NewHTTPAPIAdapter(
		application.NewShortURLApplicationService(validator, domain.NewShortCodeGenerator(), store, store, "http://localhost:8080"),
		application.NewRedirectService(validator, store, store),
		ratelimit.NewGuard(),
		httpapi.NewErrorMapper(),
		observability.NewHooks(),
	))
	internal := postURL(t, closed, "https://example.com/closed")
	assertError(t, internal, http.StatusInternalServerError, "internal_error")
	assertBaselineHeaders(t, internal)
	assertJSONNoStore(t, internal)
}

func TestHSTSOnlyWhenOptInAndTLS(t *testing.T) {
	srv := newServerWithHSTS(t, true)
	plain := postURL(t, srv, "https://example.com/plain")
	if plain.Code != http.StatusCreated {
		t.Fatalf("status %d", plain.Code)
	}
	assertBaselineHeaders(t, plain)
	if plain.Header().Get("Strict-Transport-Security") != "" {
		t.Fatal("HSTS on plain HTTP")
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.TLS = &tls.ConnectionState{}
	srv.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("health %d", res.Code)
	}
	assertBaselineHeaders(t, res)
	if got := res.Header().Get("Strict-Transport-Security"); got != "max-age=31536000" {
		t.Fatalf("HSTS %q", got)
	}

	off := newServer(t)
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/health", nil)
	req.TLS = &tls.ConnectionState{}
	off.ServeHTTP(res, req)
	if res.Header().Get("Strict-Transport-Security") != "" {
		t.Fatal("HSTS without opt-in")
	}
}

func newServerWithHSTS(t *testing.T, hsts bool) http.Handler {
	t.Helper()
	store, err := sqlite.Open(filepath.Join(t.TempDir(), "hsts.db"))
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
	adapter := httpapi.NewHTTPAPIAdapter(appSvc, redirectSvc, ratelimit.NewGuard(), httpapi.NewErrorMapper(), observability.NewHooks())
	return httpapi.NewRouterWithOptions(adapter, httpapi.RouterOptions{HSTS: hsts})
}

func assertBaselineHeaders(t *testing.T, res *httptest.ResponseRecorder) {
	t.Helper()
	if got := res.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("nosniff %q", got)
	}
	if got := res.Header().Get("Referrer-Policy"); got != "no-referrer" {
		t.Fatalf("referrer %q", got)
	}
	if got := res.Header().Get("X-Frame-Options"); got != "DENY" {
		t.Fatalf("frame %q", got)
	}
	if got := res.Header().Get("Content-Security-Policy"); got != "default-src 'none'; frame-ancestors 'none'" {
		t.Fatalf("csp %q", got)
	}
	if got := res.Header().Get("Server"); got != "" {
		t.Fatalf("server header %q", got)
	}
	if got := res.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("cors %q", got)
	}
}

func assertJSONNoStore(t *testing.T, res *httptest.ResponseRecorder) {
	t.Helper()
	if ct := res.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("content-type %q", ct)
	}
	if got := res.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("cache-control %q", got)
	}
}
