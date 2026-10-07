package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/444112029012/short-url/internal/application"
	"github.com/444112029012/short-url/internal/domain"
	"github.com/444112029012/short-url/internal/httpapi"
	"github.com/444112029012/short-url/internal/infra/sqlite"
	"github.com/444112029012/short-url/internal/observability"
	"github.com/444112029012/short-url/internal/ratelimit"
)

type recHooks struct {
	mu     sync.Mutex
	events []string
}

func (h *recHooks) add(event string) {
	h.mu.Lock()
	h.events = append(h.events, event)
	h.mu.Unlock()
}

func (h *recHooks) OnCreateSuccess(shortCode string)   { h.add("create_success:" + shortCode) }
func (h *recHooks) OnCreateFailure(errorType string)   { h.add("create_failure:" + errorType) }
func (h *recHooks) OnRedirectSuccess(shortCode string) { h.add("redirect_success:" + shortCode) }
func (h *recHooks) OnRedirectFailure(errorType string) { h.add("redirect_failure:" + errorType) }
func (h *recHooks) OnStatsSuccess(shortCode string)    { h.add("stats_success:" + shortCode) }
func (h *recHooks) OnStatsFailure(errorType string)    { h.add("stats_failure:" + errorType) }

func (h *recHooks) snapshot() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]string, len(h.events))
	copy(out, h.events)
	return out
}

func TestUT_OBS_01_AdapterCallsHooks(t *testing.T) {
	path := filepath.Join(t.TempDir(), "obs.db")
	store := openStore(t, path)
	hooks := &recHooks{}
	srv := routerWith(t, store, ratelimit.NewGuard(), hooks)

	created := postURL(t, srv, "https://example.com/obs")
	if created.Code != http.StatusCreated {
		t.Fatalf("create %d %s", created.Code, created.Body.String())
	}
	code, _ := decodeMap(t, created.Body.Bytes())["short_code"].(string)

	bad := postURL(t, srv, "ftp://example.com/a")
	assertError(t, bad, http.StatusBadRequest, "invalid_url")

	redirect := httptest.NewRecorder()
	srv.ServeHTTP(redirect, httptest.NewRequest(http.MethodGet, "/"+code, nil))
	if redirect.Code != http.StatusFound {
		t.Fatalf("redirect %d", redirect.Code)
	}

	missing := httptest.NewRecorder()
	srv.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/NoSuch01", nil))
	assertError(t, missing, http.StatusNotFound, "not_found")

	malformed := httptest.NewRecorder()
	srv.ServeHTTP(malformed, httptest.NewRequest(http.MethodGet, "/bad", nil))
	assertError(t, malformed, http.StatusBadRequest, "invalid_url")

	stats := getJSON(t, srv, "/api/v1/urls/"+code+"/stats")
	if stats.Code != http.StatusOK {
		t.Fatalf("stats %d", stats.Code)
	}
	statsMissing := getJSON(t, srv, "/api/v1/urls/NoSuch01/stats")
	assertError(t, statsMissing, http.StatusNotFound, "not_found")
	statsBad := getJSON(t, srv, "/api/v1/urls/bad!code/stats")
	assertError(t, statsBad, http.StatusBadRequest, "invalid_url")

	limited := routerWith(t, store, blockingLimiter{}, hooks)
	tooMany := postURL(t, limited, "https://example.com/limited")
	assertError(t, tooMany, http.StatusTooManyRequests, "rate_limited")
	redirLimit := httptest.NewRecorder()
	limited.ServeHTTP(redirLimit, httptest.NewRequest(http.MethodGet, "/"+code, nil))
	assertError(t, redirLimit, http.StatusTooManyRequests, "rate_limited")

	got := hooks.snapshot()
	want := []string{
		"create_success:" + code,
		"create_failure:invalid_url",
		"redirect_success:" + code,
		"redirect_failure:not_found",
		"redirect_failure:invalid_url",
		"stats_success:" + code,
		"stats_failure:not_found",
		"stats_failure:invalid_url",
		"create_failure:rate_limited",
		"redirect_failure:rate_limited",
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("hooks\n got %v\nwant %v", got, want)
	}
}

func TestUT_OBS_01_LogsOmitSensitiveValues(t *testing.T) {
	const (
		longURL = "https://sensitive.example/private-path?token=body-secret"
		authz   = "Bearer super-secret-token"
		cookie  = "session=secret-cookie"
		ip      = "203.0.113.55"
	)
	path := filepath.Join(t.TempDir(), "sensitive.db")
	store := openStore(t, path)
	var buf bytes.Buffer
	hooks := observability.NewHooksWithLogger(slog.New(slog.NewJSONHandler(&buf, nil)))
	srv := routerWith(t, store, ratelimit.NewGuard(), hooks)

	payload, err := json.Marshal(map[string]string{"url": longURL})
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/urls", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authz)
	req.Header.Set("Cookie", cookie)
	req.RemoteAddr = ip + ":443"
	req.Header.Set("X-Forwarded-For", "198.51.100.77")
	srv.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("create %d %s", res.Code, res.Body.String())
	}
	code, _ := decodeMap(t, res.Body.Bytes())["short_code"].(string)

	redirect := httptest.NewRecorder()
	rreq := httptest.NewRequest(http.MethodGet, "/"+code, nil)
	rreq.Header.Set("Authorization", authz)
	rreq.Header.Set("Cookie", cookie)
	rreq.RemoteAddr = ip + ":443"
	srv.ServeHTTP(redirect, rreq)

	stats := httptest.NewRecorder()
	sreq := httptest.NewRequest(http.MethodGet, "/api/v1/urls/"+code+"/stats", nil)
	sreq.Header.Set("Authorization", authz)
	sreq.RemoteAddr = ip + ":443"
	srv.ServeHTTP(stats, sreq)

	bad := httptest.NewRecorder()
	breq := httptest.NewRequest(http.MethodPost, "/api/v1/urls", strings.NewReader(`{"url":"ftp://secret.example/hidden"}`))
	breq.Header.Set("Content-Type", "application/json")
	breq.Header.Set("Authorization", authz)
	breq.RemoteAddr = ip + ":9"
	srv.ServeHTTP(bad, breq)

	out := buf.String()
	for _, forbidden := range []string{
		longURL, "sensitive.example", "private-path", "body-secret",
		authz, "super-secret-token", "Bearer",
		cookie, "secret-cookie", "session=",
		"Authorization", "Cookie",
		ip, "198.51.100.77",
		"ftp://secret.example/hidden", "secret.example", "hidden",
	} {
		if strings.Contains(out, forbidden) {
			t.Fatalf("log contains %q\n%s", forbidden, out)
		}
	}
	if !strings.Contains(out, code) || !strings.Contains(out, "invalid_url") {
		t.Fatalf("expected short code and error type in %s", out)
	}
}

func openStore(t *testing.T, path string) *sqlite.URLStore {
	t.Helper()
	store, err := sqlite.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if err := store.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}
	return store
}

func routerWith(t *testing.T, store *sqlite.URLStore, limiter httpapi.RateLimitChecker, hooks httpapi.EventHooks) http.Handler {
	t.Helper()
	validator := domain.NewURLValidator()
	appSvc := application.NewShortURLApplicationService(validator, domain.NewShortCodeGenerator(), store, store, "http://localhost:8080")
	redirectSvc := application.NewRedirectService(validator, store, store)
	adapter := httpapi.NewHTTPAPIAdapter(appSvc, redirectSvc, limiter, httpapi.NewErrorMapper(), hooks)
	return httpapi.NewRouter(adapter)
}
