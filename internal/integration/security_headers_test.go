package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/444112029012/short-url/internal/application"
	"github.com/444112029012/short-url/internal/config"
	"github.com/444112029012/short-url/internal/domain"
	"github.com/444112029012/short-url/internal/httpapi"
	"github.com/444112029012/short-url/internal/infra/sqlite"
	"github.com/444112029012/short-url/internal/observability"
	"github.com/444112029012/short-url/internal/ratelimit"
)

// ENG-016 / SEC-008: end-to-end response header checks against a real
// httptest server wired like cmd/shorturl (router, middleware, adapter,
// temp SQLite). Expected values mirror internal/httpapi/security_headers.go.

const (
	wantHSTS            = "max-age=31536000"
	shortCodeLongURL    = "https://example.com/eng-016"
	headerHSTS          = "Strict-Transport-Security"
	headerCacheControl  = "Cache-Control"
	headerContentType   = "Content-Type"
	contentTypeJSONOnly = "application/json"
)

var baselineSecurityHeaders = []struct{ name, value string }{
	{"X-Content-Type-Options", "nosniff"},
	{"Referrer-Policy", "no-referrer"},
	{"X-Frame-Options", "DENY"},
	{"Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'"},
}

// cacheRule is what #8 defines for Cache-Control on a response kind:
// application/json responses get exactly "no-store"; the 302 redirect and
// non-JSON router responses get no Cache-Control at all.
type cacheRule int

const (
	cacheNoStoreJSON cacheRule = iota
	cacheAbsent
)

type observed struct {
	name    string
	status  int
	errType string
	cache   cacheRule
	res     *http.Response
	body    []byte
}

// newHeaderServer wires the same stack as main: config-driven RouterOptions,
// NewRouterWithOptions, the 6-arg adapter with main's default limits
// NewGuard(30, 120) and nil trusted proxies, and a temp SQLite store.
// useTLS selects httptest.NewTLSServer.
func newHeaderServer(t *testing.T, opts httpapi.RouterOptions, useTLS bool) (*httptest.Server, *http.Client) {
	t.Helper()
	store, err := sqlite.Open(filepath.Join(t.TempDir(), "headers.db"))
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Migrate(context.Background()); err != nil {
		_ = store.Close()
		t.Fatal(err)
	}
	validator := domain.NewURLValidator()
	appSvc := application.NewShortURLApplicationService(validator, domain.NewShortCodeGenerator(), store, store, "http://localhost:8080")
	redirectSvc := application.NewRedirectService(validator, store, store)
	adapter := httpapi.NewHTTPAPIAdapter(appSvc, redirectSvc, ratelimit.NewGuard(30, 120), httpapi.NewErrorMapper(), observability.NewHooks(), nil)
	handler := httpapi.NewRouterWithOptions(adapter, opts)
	var srv *httptest.Server
	if useTLS {
		srv = httptest.NewTLSServer(handler)
	} else {
		srv = httptest.NewServer(handler)
	}
	t.Cleanup(func() {
		srv.Close()
		_ = store.Close()
	})
	client := srv.Client()
	client.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return srv, client
}

// routerOptionsFromEnv builds RouterOptions the way cmd/shorturl does:
// config.Load reads HSTS_ENABLED and main passes cfg.HSTSEnabled through.
func routerOptionsFromEnv(t *testing.T, hstsEnabled string) httpapi.RouterOptions {
	t.Helper()
	t.Setenv("HSTS_ENABLED", hstsEnabled)
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	return httpapi.RouterOptions{HSTS: cfg.HSTSEnabled}
}

func doRequest(t *testing.T, client *http.Client, method, url, contentType, body string) (*http.Response, []byte) {
	t.Helper()
	var rd io.Reader
	if body != "" {
		rd = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, rd)
	if err != nil {
		t.Fatal(err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	return res, raw
}

// exerciseAllResponseKinds drives every response class the API produces:
// 201 create, 302 redirect, 200 stats, 400 invalid_url (bad scheme,
// non-JSON Content-Type, malformed code), 404 not_found (redirect, stats),
// router 404 and 405. The create handler answers a wrong Content-Type with
// 400 invalid_url, not 415; that case is covered as such.
func exerciseAllResponseKinds(t *testing.T, srv *httptest.Server, client *http.Client) []observed {
	t.Helper()
	var out []observed
	add := func(name string, status int, errType string, cache cacheRule, res *http.Response, body []byte) {
		out = append(out, observed{name: name, status: status, errType: errType, cache: cache, res: res, body: body})
	}

	payload, err := json.Marshal(map[string]string{"url": shortCodeLongURL})
	if err != nil {
		t.Fatal(err)
	}
	res, body := doRequest(t, client, http.MethodPost, srv.URL+"/api/v1/urls", contentTypeJSONOnly, string(payload))
	add("POST create 201", http.StatusCreated, "", cacheNoStoreJSON, res, body)
	var created struct {
		ShortCode string `json:"short_code"`
	}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&created); err != nil || len(created.ShortCode) != 8 {
		t.Fatalf("create body %s (%v)", body, err)
	}

	res, body = doRequest(t, client, http.MethodGet, srv.URL+"/"+created.ShortCode, "", "")
	add("GET redirect 302", http.StatusFound, "", cacheAbsent, res, body)
	if got := res.Header.Get("Location"); got != shortCodeLongURL {
		t.Fatalf("302 Location %q, want %q", got, shortCodeLongURL)
	}

	res, body = doRequest(t, client, http.MethodGet, srv.URL+"/api/v1/urls/"+created.ShortCode+"/stats", "", "")
	add("GET stats 200", http.StatusOK, "", cacheNoStoreJSON, res, body)

	res, body = doRequest(t, client, http.MethodPost, srv.URL+"/api/v1/urls", contentTypeJSONOnly, `{"url":"ftp://example.com/a"}`)
	add("POST create 400 invalid_url scheme", http.StatusBadRequest, "invalid_url", cacheNoStoreJSON, res, body)

	res, body = doRequest(t, client, http.MethodPost, srv.URL+"/api/v1/urls", "text/plain", string(payload))
	add("POST create 400 invalid_url non-JSON Content-Type", http.StatusBadRequest, "invalid_url", cacheNoStoreJSON, res, body)

	res, body = doRequest(t, client, http.MethodGet, srv.URL+"/bad", "", "")
	add("GET redirect 400 invalid_url malformed code", http.StatusBadRequest, "invalid_url", cacheNoStoreJSON, res, body)

	res, body = doRequest(t, client, http.MethodGet, srv.URL+"/NoSuch01", "", "")
	add("GET redirect 404 not_found", http.StatusNotFound, "not_found", cacheNoStoreJSON, res, body)
	if loc := res.Header.Get("Location"); loc != "" {
		t.Fatalf("404 Location %q", loc)
	}

	res, body = doRequest(t, client, http.MethodGet, srv.URL+"/api/v1/urls/NoSuch01/stats", "", "")
	add("GET stats 404 not_found", http.StatusNotFound, "not_found", cacheNoStoreJSON, res, body)

	res, body = doRequest(t, client, http.MethodGet, srv.URL+"/api/v1/does-not-exist", "", "")
	add("GET router 404", http.StatusNotFound, "", cacheAbsent, res, body)

	res, body = doRequest(t, client, http.MethodDelete, srv.URL+"/api/v1/urls", "", "")
	add("DELETE create route 405", http.StatusMethodNotAllowed, "", cacheAbsent, res, body)

	return out
}

func assertHeaderOnce(t *testing.T, name string, h http.Header, key, want string) {
	t.Helper()
	got := h.Values(key)
	if len(got) != 1 {
		t.Fatalf("%s: %s appears %d times %q, want exactly once", name, key, len(got), got)
	}
	if got[0] != want {
		t.Fatalf("%s: %s = %q, want %q", name, key, got[0], want)
	}
}

func assertObservedResponse(t *testing.T, o observed) {
	t.Helper()
	if o.res.StatusCode != o.status {
		t.Fatalf("%s: status %d, want %d (body %s)", o.name, o.res.StatusCode, o.status, o.body)
	}
	for _, hdr := range baselineSecurityHeaders {
		assertHeaderOnce(t, o.name, o.res.Header, hdr.name, hdr.value)
	}
	switch o.cache {
	case cacheNoStoreJSON:
		assertHeaderOnce(t, o.name, o.res.Header, headerContentType, contentTypeJSONOnly)
		assertHeaderOnce(t, o.name, o.res.Header, headerCacheControl, "no-store")
	case cacheAbsent:
		if got := o.res.Header.Values(headerCacheControl); len(got) != 0 {
			t.Fatalf("%s: Cache-Control %q, want absent", o.name, got)
		}
	}
	if o.errType != "" {
		var payload struct {
			Error struct {
				Type string `json:"type"`
			} `json:"error"`
		}
		if err := json.Unmarshal(o.body, &payload); err != nil || payload.Error.Type != o.errType {
			t.Fatalf("%s: error body %s, want type %s (%v)", o.name, o.body, o.errType, err)
		}
	}
}

func TestENG016_SecurityHeadersOnEveryResponse(t *testing.T) {
	// SEC-008: the four headers are present exactly once on success,
	// redirect, and every error class; JSON gets Cache-Control: no-store,
	// the 302 gets none; default (HSTS off) plain HTTP has no HSTS.
	srv, client := newHeaderServer(t, httpapi.RouterOptions{}, false)
	for _, o := range exerciseAllResponseKinds(t, srv, client) {
		t.Run(o.name, func(t *testing.T) {
			assertObservedResponse(t, o)
			if got := o.res.Header.Values(headerHSTS); len(got) != 0 {
				t.Fatalf("HSTS %q on plain HTTP", got)
			}
		})
	}
}

func TestENG016_HSTSAbsentOnPlainHTTPWhenEnabled(t *testing.T) {
	// HSTS_ENABLED=true must not advertise HSTS over plain HTTP.
	opts := routerOptionsFromEnv(t, "true")
	if !opts.HSTS {
		t.Fatal("HSTS_ENABLED=true did not enable RouterOptions.HSTS")
	}
	srv, client := newHeaderServer(t, opts, false)
	for _, o := range exerciseAllResponseKinds(t, srv, client) {
		t.Run(o.name, func(t *testing.T) {
			if o.res.TLS != nil {
				t.Fatal("expected a plain HTTP response")
			}
			assertObservedResponse(t, o)
			if got := o.res.Header.Values(headerHSTS); len(got) != 0 {
				t.Fatalf("HSTS %q on plain HTTP with HSTS enabled", got)
			}
		})
	}
}

func TestENG016_HSTSAbsentOnTLSWhenDisabled(t *testing.T) {
	// HSTS_ENABLED=false: a real TLS connection still gets no HSTS.
	opts := routerOptionsFromEnv(t, "false")
	if opts.HSTS {
		t.Fatal("HSTS_ENABLED=false enabled RouterOptions.HSTS")
	}
	srv, client := newHeaderServer(t, opts, true)
	for _, o := range exerciseAllResponseKinds(t, srv, client) {
		t.Run(o.name, func(t *testing.T) {
			if o.res.TLS == nil {
				t.Fatal("expected a TLS response")
			}
			assertObservedResponse(t, o)
			if got := o.res.Header.Values(headerHSTS); len(got) != 0 {
				t.Fatalf("HSTS %q on TLS with HSTS disabled", got)
			}
		})
	}
}

func TestENG016_HSTSPresentOnTLSWhenEnabled(t *testing.T) {
	// HSTS_ENABLED=true over a real TLS connection: every response kind
	// carries Strict-Transport-Security: max-age=31536000 exactly once.
	opts := routerOptionsFromEnv(t, "true")
	srv, client := newHeaderServer(t, opts, true)
	for _, o := range exerciseAllResponseKinds(t, srv, client) {
		t.Run(o.name, func(t *testing.T) {
			if o.res.TLS == nil {
				t.Fatal("expected a TLS response")
			}
			assertObservedResponse(t, o)
			assertHeaderOnce(t, o.name, o.res.Header, headerHSTS, wantHSTS)
		})
	}
}
