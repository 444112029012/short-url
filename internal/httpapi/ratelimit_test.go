package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestAdapterRateLimit429AndRedirectDoesNotCount(t *testing.T) {
	srv := newLimitedServer(t, 2, 1, "")
	first := postURLFrom(t, srv, "https://example.com/rl", "203.0.113.10:1000", "")
	if first.Code != http.StatusCreated {
		t.Fatalf("create %d %s", first.Code, first.Body.String())
	}
	code, _ := decodeMap(t, first.Body.Bytes())["short_code"].(string)

	second := postURLFrom(t, srv, "https://example.com/rl2", "203.0.113.10:1001", "198.51.100.1")
	if second.Code != http.StatusCreated {
		t.Fatalf("at-limit create %d %s", second.Code, second.Body.String())
	}
	third := postURLFrom(t, srv, "https://example.com/rl3", "203.0.113.10:1002", "198.51.100.2")
	assertError(t, third, http.StatusTooManyRequests, "rate_limited")
	body := decodeMap(t, third.Body.Bytes())
	errObj, _ := body["error"].(map[string]any)
	if errObj["message"] != domain.MsgTooManyRequests {
		t.Fatalf("429 message %v", errObj["message"])
	}

	redirect := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/"+code, nil)
	req.RemoteAddr = "203.0.113.10:2000"
	req.Header.Set("X-Forwarded-For", "198.51.100.9")
	srv.ServeHTTP(redirect, req)
	if redirect.Code != http.StatusFound || redirect.Header().Get("Location") != "https://example.com/rl" {
		t.Fatalf("redirect %d location %q", redirect.Code, redirect.Header().Get("Location"))
	}

	over := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/"+code, nil)
	req.RemoteAddr = "203.0.113.10:2001"
	req.Header.Set("X-Forwarded-For", "198.51.100.8")
	srv.ServeHTTP(over, req)
	assertError(t, over, http.StatusTooManyRequests, "rate_limited")
	if over.Header().Get("Location") != "" {
		t.Fatal("over-limit redirect set Location")
	}

	stats := getJSON(t, srv, "/api/v1/urls/"+code+"/stats")
	if stats.Code != http.StatusOK || decodeMap(t, stats.Body.Bytes())["click_count"] != float64(1) {
		t.Fatalf("click_count after limited redirect %d %s", stats.Code, stats.Body.String())
	}

	for i := 0; i < 5; i++ {
		again := getJSON(t, srv, "/api/v1/urls/"+code+"/stats")
		if again.Code != http.StatusOK {
			t.Fatalf("stats rate limited %d", again.Code)
		}
	}
}

func TestSpoofedXFFDoesNotBypassCreateLimit(t *testing.T) {
	srv := newLimitedServer(t, 1, 120, "")
	ok := postURLFrom(t, srv, "https://example.com/a", "203.0.113.20:1", "198.51.100.1")
	if ok.Code != http.StatusCreated {
		t.Fatalf("first %d %s", ok.Code, ok.Body.String())
	}
	spoofed := postURLFrom(t, srv, "https://example.com/b", "203.0.113.20:2", "198.51.100.2")
	assertError(t, spoofed, http.StatusTooManyRequests, "rate_limited")

	other := postURLFrom(t, srv, "https://example.com/c", "203.0.113.21:1", "198.51.100.1")
	if other.Code != http.StatusCreated {
		t.Fatalf("independent peer %d %s", other.Code, other.Body.String())
	}
}

func TestTrustedProxyUsesForwardedClient(t *testing.T) {
	srv := newLimitedServer(t, 1, 1, "10.0.0.0/8")
	a := postURLFrom(t, srv, "https://example.com/a", "10.1.0.2:4000", "198.51.100.20")
	if a.Code != http.StatusCreated {
		t.Fatalf("trusted client a %d %s", a.Code, a.Body.String())
	}
	b := postURLFrom(t, srv, "https://example.com/b", "10.1.0.2:4001", "198.51.100.21")
	if b.Code != http.StatusCreated {
		t.Fatalf("trusted client b %d %s", b.Code, b.Body.String())
	}
	a2 := postURLFrom(t, srv, "https://example.com/c", "10.1.0.9:4002", "198.51.100.20")
	assertError(t, a2, http.StatusTooManyRequests, "rate_limited")

	direct := postURLFrom(t, srv, "https://example.com/d", "203.0.113.30:9", "198.51.100.20")
	if direct.Code != http.StatusCreated {
		t.Fatalf("untrusted peer must ignore XFF %d %s", direct.Code, direct.Body.String())
	}
}

func newLimitedServer(t *testing.T, createLimit, redirectLimit int, trustedCIDRs string) http.Handler {
	t.Helper()
	store, err := sqlite.Open(filepath.Join(t.TempDir(), "rl.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if err := store.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}
	trusted, err := ratelimit.ParseTrustedCIDRs(trustedCIDRs)
	if err != nil {
		t.Fatal(err)
	}
	validator := domain.NewURLValidator()
	appSvc := application.NewShortURLApplicationService(validator, domain.NewShortCodeGenerator(), store, store, "http://localhost:8080")
	redirectSvc := application.NewRedirectService(validator, store, store)
	adapter := httpapi.NewHTTPAPIAdapter(
		appSvc,
		redirectSvc,
		ratelimit.NewGuard(createLimit, redirectLimit),
		httpapi.NewErrorMapper(),
		observability.NewHooks(),
		trusted,
	)
	return httpapi.NewRouter(adapter)
}

func postURLFrom(t *testing.T, srv http.Handler, rawURL, remote, xff string) *httptest.ResponseRecorder {
	t.Helper()
	payload, err := json.Marshal(map[string]string{"url": rawURL})
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/urls", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = remote
	if xff != "" {
		req.Header.Set("X-Forwarded-For", xff)
	}
	srv.ServeHTTP(res, req)
	return res
}
