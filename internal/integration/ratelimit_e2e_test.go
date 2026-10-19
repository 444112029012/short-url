package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/444112029012/short-url/internal/application"
	"github.com/444112029012/short-url/internal/config"
	"github.com/444112029012/short-url/internal/domain"
	"github.com/444112029012/short-url/internal/httpapi"
	"github.com/444112029012/short-url/internal/infra/sqlite"
	"github.com/444112029012/short-url/internal/observability"
	"github.com/444112029012/short-url/internal/ratelimit"
)

// ENG-011 / SEC-007 / NFR-003: end-to-end rate-limit checks against a real
// httptest server wired like cmd/shorturl. Limits and trusted proxies come
// from the real environment variables through config.Load, then
// ParseTrustedCIDRs, NewGuard, and the 6-arg adapter, exactly as main does.
//
// Contract under test (method spec v0.2 §10, OpenAPI 429, API notes §3):
//   - over the per-source, per-minute threshold: 429 with
//     {"error":{"type":"rate_limited","message":"Too many requests"}}
//   - create does not store a code; redirect does not send Location and does
//     not increment click_count
//   - create and redirect quotas are independent
//   - forwarding headers are honored only when the TCP peer is in
//     TRUSTED_PROXY_CIDRS
//
// The implementation sends no Retry-After header and the spec does not
// require one; rlAssertLimited pins that so a change is a visible decision.
//
// Addresses below are documentation ranges (RFC 5737 / RFC 3849). Tests
// never print them; assertion messages use labels only.

const (
	rlCreatePath = "/api/v1/urls"
	rlWantBody   = `{"error":{"type":"rate_limited","message":"Too many requests"}}`
)

var rlSecurityHeaders = []struct{ name, value string }{
	{"X-Content-Type-Options", "nosniff"},
	{"Referrer-Policy", "no-referrer"},
	{"X-Frame-Options", "DENY"},
	{"Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'"},
}

type rlEnv struct {
	createPerMin   string
	redirectPerMin string
	trustedCIDRs   string
}

type rlServer struct {
	srv    *httptest.Server
	dbPath string
	logs   *rlLogBuffer
}

// rlLogBuffer is a goroutine-safe sink for the observability JSON logger.
type rlLogBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *rlLogBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *rlLogBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// newRateLimitServer builds the cmd/shorturl stack from env. The observability
// hooks write the same JSON lines main writes to stderr, but into a buffer so
// the test can inspect them. The server's own error log is captured too.
func newRateLimitServer(t *testing.T, env rlEnv) *rlServer {
	t.Helper()
	t.Setenv("RATE_LIMIT_CREATE_PER_MIN", env.createPerMin)
	t.Setenv("RATE_LIMIT_REDIRECT_PER_MIN", env.redirectPerMin)
	t.Setenv("TRUSTED_PROXY_CIDRS", env.trustedCIDRs)
	t.Setenv("HSTS_ENABLED", "")
	dbPath := filepath.Join(t.TempDir(), "ratelimit.db")
	t.Setenv("DATABASE_PATH", dbPath)
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	trusted, err := ratelimit.ParseTrustedCIDRs(cfg.TrustedProxyCIDRs)
	if err != nil {
		t.Fatalf("ParseTrustedCIDRs: %v", err)
	}
	store, err := sqlite.Open(cfg.DatabasePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Migrate(context.Background()); err != nil {
		_ = store.Close()
		t.Fatal(err)
	}
	logs := &rlLogBuffer{}
	validator := domain.NewURLValidator()
	appSvc := application.NewShortURLApplicationService(validator, domain.NewShortCodeGenerator(), store, store, cfg.BaseURL)
	redirectSvc := application.NewRedirectService(validator, store, store)
	adapter := httpapi.NewHTTPAPIAdapter(
		appSvc,
		redirectSvc,
		ratelimit.NewGuard(cfg.RateLimitCreatePerMin, cfg.RateLimitRedirectPerMin),
		httpapi.NewErrorMapper(),
		observability.NewHooksWithLogger(slog.New(slog.NewJSONHandler(logs, nil))),
		trusted,
	)
	srv := httptest.NewUnstartedServer(httpapi.NewRouterWithOptions(adapter, httpapi.RouterOptions{HSTS: cfg.HSTSEnabled}))
	srv.Config.ErrorLog = slog.NewLogLogger(slog.NewTextHandler(logs, nil), slog.LevelError)
	srv.Start()
	t.Cleanup(func() {
		srv.Close()
		_ = store.Close()
	})
	return &rlServer{srv: srv, dbPath: dbPath, logs: logs}
}

// client returns an HTTP client whose TCP connections originate from
// localIP (nil means the default loopback source). On Linux every 127/8
// address is local, so distinct loopback sources give distinct real peers.
func (s *rlServer) client(t *testing.T, localIP net.IP) *http.Client {
	t.Helper()
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	if localIP != nil {
		dialer.LocalAddr = &net.TCPAddr{IP: localIP}
	}
	transport := &http.Transport{DialContext: dialer.DialContext}
	t.Cleanup(transport.CloseIdleConnections)
	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

type rlResult struct {
	status int
	header http.Header
	body   []byte
}

func rlDo(t *testing.T, client *http.Client, method, url string, body []byte, header http.Header) rlResult {
	t.Helper()
	var rd io.Reader
	if body != nil {
		rd = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, url, rd)
	if err != nil {
		t.Fatal(err)
	}
	for k, vs := range header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
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
	return rlResult{status: res.StatusCode, header: res.Header, body: raw}
}

func (s *rlServer) create(t *testing.T, client *http.Client, longURL string, header http.Header) rlResult {
	t.Helper()
	payload, err := json.Marshal(map[string]string{"url": longURL})
	if err != nil {
		t.Fatal(err)
	}
	return rlDo(t, client, http.MethodPost, s.srv.URL+rlCreatePath, payload, header)
}

func (s *rlServer) redirect(t *testing.T, client *http.Client, code string, header http.Header) rlResult {
	t.Helper()
	return rlDo(t, client, http.MethodGet, s.srv.URL+"/"+code, nil, header)
}

func (s *rlServer) clicks(t *testing.T, client *http.Client, code string) int {
	t.Helper()
	res := rlDo(t, client, http.MethodGet, s.srv.URL+rlCreatePath+"/"+code+"/stats", nil, nil)
	if res.status != http.StatusOK {
		t.Fatalf("stats status %d body %s", res.status, res.body)
	}
	var body struct {
		ClickCount int `json:"click_count"`
	}
	if err := json.Unmarshal(res.body, &body); err != nil {
		t.Fatal(err)
	}
	return body.ClickCount
}

func rlShortCode(t *testing.T, res rlResult) string {
	t.Helper()
	if res.status != http.StatusCreated {
		t.Fatalf("create status %d body %s", res.status, res.body)
	}
	var body struct {
		ShortCode string `json:"short_code"`
	}
	if err := json.Unmarshal(res.body, &body); err != nil || body.ShortCode == "" {
		t.Fatalf("create body %s", res.body)
	}
	return body.ShortCode
}

func xff(values ...string) http.Header {
	h := make(http.Header)
	for _, v := range values {
		h.Add("X-Forwarded-For", v)
	}
	return h
}

// rlAssertLimited checks the full 429 contract: status, exact error body,
// JSON content type, SEC-008 headers, Cache-Control: no-store, no Location,
// no HSTS on plain HTTP, and no Retry-After (not specified, not sent).
func rlAssertLimited(t *testing.T, label string, res rlResult) {
	t.Helper()
	if res.status != http.StatusTooManyRequests {
		t.Fatalf("%s: status %d body %s, want 429", label, res.status, res.body)
	}
	if got := strings.TrimSpace(string(res.body)); got != rlWantBody {
		t.Fatalf("%s: body %s, want %s", label, got, rlWantBody)
	}
	var env map[string]map[string]any
	if err := json.Unmarshal(res.body, &env); err != nil || len(env) != 1 || len(env["error"]) != 2 {
		t.Fatalf("%s: error envelope must be exactly {error:{type,message}}: %s", label, res.body)
	}
	if got := res.header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("%s: Content-Type %q", label, got)
	}
	for _, h := range rlSecurityHeaders {
		if got := res.header.Values(h.name); len(got) != 1 || got[0] != h.value {
			t.Fatalf("%s: %s = %q, want exactly [%q]", label, h.name, got, h.value)
		}
	}
	if got := res.header.Values("Cache-Control"); len(got) != 1 || got[0] != "no-store" {
		t.Fatalf("%s: Cache-Control %q, want [no-store]", label, got)
	}
	if got := res.header.Get("Location"); got != "" {
		t.Fatalf("%s: 429 must not carry Location, got %q", label, got)
	}
	if got := res.header.Get("Strict-Transport-Security"); got != "" {
		t.Fatalf("%s: plain HTTP must not carry HSTS, got %q", label, got)
	}
	if got := res.header.Get("Retry-After"); got != "" {
		t.Fatalf("%s: Retry-After %q is not part of the current contract", label, got)
	}
}

func TestENG011_CreateQuotaExhaustedReturns429(t *testing.T) {
	s := newRateLimitServer(t, rlEnv{createPerMin: "3", redirectPerMin: "100"})
	c := s.client(t, nil)
	for i := 0; i < 3; i++ {
		rlShortCode(t, s.create(t, c, fmt.Sprintf("https://example.com/create-%d", i), nil))
	}
	if n := countRows(t, s.dbPath); n != 3 {
		t.Fatalf("rows after quota %d, want 3", n)
	}
	for i := 0; i < 3; i++ {
		rlAssertLimited(t, fmt.Sprintf("create over quota #%d", i+1),
			s.create(t, c, fmt.Sprintf("https://example.com/over-%d", i), nil))
	}
	// Rejections store nothing and do not reset the window.
	if n := countRows(t, s.dbPath); n != 3 {
		t.Fatalf("rows after rejections %d, want 3", n)
	}
	// The guard runs before body decoding: an invalid body over quota is
	// still rate_limited, not invalid_url (method spec: 先 RateLimitGuard).
	rlAssertLimited(t, "create over quota with invalid body",
		rlDo(t, c, http.MethodPost, s.srv.URL+rlCreatePath, []byte(`not json`), nil))
}

func TestENG011_DefaultLimitsFromConfig(t *testing.T) {
	// Empty env means the documented defaults: 30 creates/min/source.
	s := newRateLimitServer(t, rlEnv{})
	c := s.client(t, nil)
	for i := 0; i < 30; i++ {
		rlShortCode(t, s.create(t, c, fmt.Sprintf("https://example.com/default-%d", i), nil))
	}
	rlAssertLimited(t, "create #31 at default limit", s.create(t, c, "https://example.com/default-31", nil))
}

func TestENG011_RedirectQuotaExhaustedReturns429WithoutCounting(t *testing.T) {
	s := newRateLimitServer(t, rlEnv{createPerMin: "5", redirectPerMin: "2"})
	c := s.client(t, nil)
	code := rlShortCode(t, s.create(t, c, "https://example.com/redirect-target", nil))
	for i := 0; i < 2; i++ {
		res := s.redirect(t, c, code, nil)
		if res.status != http.StatusFound || res.header.Get("Location") != "https://example.com/redirect-target" {
			t.Fatalf("redirect #%d status %d Location %q", i+1, res.status, res.header.Get("Location"))
		}
	}
	for i := 0; i < 3; i++ {
		rlAssertLimited(t, fmt.Sprintf("redirect over quota #%d", i+1), s.redirect(t, c, code, nil))
	}
	// Over quota, even an unknown code is rate_limited (guard runs first).
	rlAssertLimited(t, "redirect over quota unknown code", s.redirect(t, c, "Zz9Zz9Zz", nil))
	// Stats is not rate limited and shows only the two successful redirects.
	if got := s.clicks(t, c, code); got != 2 {
		t.Fatalf("click_count %d, want 2", got)
	}
}

func TestENG011_CreateAndRedirectQuotasAreIndependent(t *testing.T) {
	s := newRateLimitServer(t, rlEnv{createPerMin: "2", redirectPerMin: "2"})
	c := s.client(t, nil)
	code := rlShortCode(t, s.create(t, c, "https://example.com/independent", nil))
	rlShortCode(t, s.create(t, c, "https://example.com/independent-2", nil))
	rlAssertLimited(t, "create exhausted", s.create(t, c, "https://example.com/independent-3", nil))

	// Exhausted create quota does not block redirects.
	for i := 0; i < 2; i++ {
		if res := s.redirect(t, c, code, nil); res.status != http.StatusFound {
			t.Fatalf("redirect #%d after create exhausted: status %d", i+1, res.status)
		}
	}
	rlAssertLimited(t, "redirect exhausted", s.redirect(t, c, code, nil))

	// And a fresh server confirms the reverse: exhausted redirects leave create open.
	s2 := newRateLimitServer(t, rlEnv{createPerMin: "3", redirectPerMin: "1"})
	c2 := s2.client(t, nil)
	code2 := rlShortCode(t, s2.create(t, c2, "https://example.com/reverse", nil))
	if res := s2.redirect(t, c2, code2, nil); res.status != http.StatusFound {
		t.Fatalf("first redirect status %d", res.status)
	}
	rlAssertLimited(t, "redirect exhausted (reverse)", s2.redirect(t, c2, code2, nil))
	rlShortCode(t, s2.create(t, c2, "https://example.com/reverse-2", nil))
}

// TestENG011_PerSourceIsolationByPeerAddress uses two real TCP source
// addresses with no trusted proxies: exhausting one peer leaves the other
// untouched.
func TestENG011_PerSourceIsolationByPeerAddress(t *testing.T) {
	s := newRateLimitServer(t, rlEnv{createPerMin: "2", redirectPerMin: "2"})
	peerA := s.client(t, net.ParseIP("127.0.0.2"))
	peerB := s.client(t, net.ParseIP("127.0.0.3"))
	if _, err := peerA.Get(s.srv.URL + "/health"); err != nil {
		t.Skipf("cannot bind a second loopback source address on this host: %v", err)
	}
	code := rlShortCode(t, s.create(t, peerA, "https://example.com/peer-a-1", nil))
	rlShortCode(t, s.create(t, peerA, "https://example.com/peer-a-2", nil))
	rlAssertLimited(t, "peer A create exhausted", s.create(t, peerA, "https://example.com/peer-a-3", nil))
	for i := 0; i < 2; i++ {
		rlShortCode(t, s.create(t, peerB, fmt.Sprintf("https://example.com/peer-b-%d", i), nil))
	}
	rlAssertLimited(t, "peer B create exhausted", s.create(t, peerB, "https://example.com/peer-b-x", nil))

	for i := 0; i < 2; i++ {
		if res := s.redirect(t, peerA, code, nil); res.status != http.StatusFound {
			t.Fatalf("peer A redirect #%d status %d", i+1, res.status)
		}
	}
	rlAssertLimited(t, "peer A redirect exhausted", s.redirect(t, peerA, code, nil))
	if res := s.redirect(t, peerB, code, nil); res.status != http.StatusFound {
		t.Fatalf("peer B redirect status %d, want 302", res.status)
	}
}

// TestENG011_UntrustedPeerCannotSpoofForwardingHeaders: with
// TRUSTED_PROXY_CIDRS empty (default) or not covering the peer, a client
// rotating X-Forwarded-For / X-Real-IP still shares its peer bucket.
func TestENG011_UntrustedPeerCannotSpoofForwardingHeaders(t *testing.T) {
	for _, tc := range []struct {
		label   string
		trusted string
	}{
		{"no trusted proxies", ""},
		{"trusted list excludes peer", "10.0.0.0/8, 192.0.2.0/24"},
	} {
		t.Run(tc.label, func(t *testing.T) {
			s := newRateLimitServer(t, rlEnv{createPerMin: "3", redirectPerMin: "3", trustedCIDRs: tc.trusted})
			c := s.client(t, nil)
			var code string
			for i := 0; i < 3; i++ {
				h := xff(fmt.Sprintf("198.51.100.%d", i+1))
				h.Set("X-Real-IP", fmt.Sprintf("203.0.113.%d", i+1))
				code = rlShortCode(t, s.create(t, c, fmt.Sprintf("https://example.com/spoof-%d", i), h))
			}
			for i := 0; i < 3; i++ {
				h := xff(fmt.Sprintf("198.51.100.%d", 50+i), "192.0.2.10")
				h.Set("X-Real-IP", fmt.Sprintf("203.0.113.%d", 50+i))
				rlAssertLimited(t, fmt.Sprintf("spoofed create #%d", i+1),
					s.create(t, c, fmt.Sprintf("https://example.com/spoof-over-%d", i), h))
			}
			for i := 0; i < 3; i++ {
				if res := s.redirect(t, c, code, xff(fmt.Sprintf("198.51.100.%d", 100+i))); res.status != http.StatusFound {
					t.Fatalf("redirect #%d status %d", i+1, res.status)
				}
			}
			rlAssertLimited(t, "spoofed redirect", s.redirect(t, c, code, xff("198.51.100.200")))
			if n := countRows(t, s.dbPath); n != 3 {
				t.Fatalf("rows %d, want 3", n)
			}
		})
	}
}

// TestENG011_TrustedProxyForwardedForIsHonored: the loopback peer is a
// trusted proxy, so the client address comes from X-Forwarded-For and two
// forwarded clients get separate quotas behind one proxy connection.
func TestENG011_TrustedProxyForwardedForIsHonored(t *testing.T) {
	s := newRateLimitServer(t, rlEnv{createPerMin: "2", redirectPerMin: "2", trustedCIDRs: "127.0.0.0/8"})
	proxy := s.client(t, nil)
	clientA := xff("198.51.100.7")
	clientB := xff("198.51.100.8")

	code := rlShortCode(t, s.create(t, proxy, "https://example.com/fwd-a-1", clientA))
	rlShortCode(t, s.create(t, proxy, "https://example.com/fwd-a-2", clientA))
	rlAssertLimited(t, "forwarded client A create exhausted", s.create(t, proxy, "https://example.com/fwd-a-3", clientA))
	rlShortCode(t, s.create(t, proxy, "https://example.com/fwd-b-1", clientB))

	// A request without forwarding headers is keyed to the proxy itself,
	// which has its own untouched quota.
	rlShortCode(t, s.create(t, proxy, "https://example.com/fwd-direct", nil))

	// X-Real-IP is used when the trusted peer sends no X-Forwarded-For.
	realA := http.Header{}
	realA.Set("X-Real-IP", "198.51.100.7")
	rlAssertLimited(t, "X-Real-IP maps to client A", s.create(t, proxy, "https://example.com/fwd-real", realA))

	for i := 0; i < 2; i++ {
		if res := s.redirect(t, proxy, code, clientA); res.status != http.StatusFound {
			t.Fatalf("client A redirect #%d status %d", i+1, res.status)
		}
	}
	rlAssertLimited(t, "forwarded client A redirect exhausted", s.redirect(t, proxy, code, clientA))
	if res := s.redirect(t, proxy, code, clientB); res.status != http.StatusFound {
		t.Fatalf("client B redirect status %d, want 302", res.status)
	}
}

// TestENG011_TrustedProxyIgnoresClientPrependedEntries: behind a trusted
// proxy the rightmost untrusted X-Forwarded-For entry wins, so a client that
// prepends random values (in the same line or as an earlier header line)
// cannot mint new identities. Trusted hops on the right are skipped.
func TestENG011_TrustedProxyIgnoresClientPrependedEntries(t *testing.T) {
	s := newRateLimitServer(t, rlEnv{createPerMin: "2", redirectPerMin: "100", trustedCIDRs: "127.0.0.0/8,10.0.0.0/8"})
	proxy := s.client(t, nil)
	headers := []http.Header{
		xff("198.51.100.1, 198.51.100.9"),
		xff("198.51.100.2, 198.51.100.9, 10.1.2.3"),
		xff("198.51.100.3", "198.51.100.9"),
		xff("garbage, 198.51.100.4", "198.51.100.9, 10.0.0.1"),
	}
	for i, h := range headers {
		res := s.create(t, proxy, fmt.Sprintf("https://example.com/prepend-%d", i), h)
		if i < 2 {
			rlShortCode(t, res)
			continue
		}
		rlAssertLimited(t, fmt.Sprintf("prepended variant #%d", i+1), res)
	}
}

// TestENG011_IPv6ForwardedClientsShareSlash64 checks SCR-007 aggregation
// through the full stack: addresses in one /64 share a bucket, another /64
// does not, and an IPv4-mapped IPv6 address shares the IPv4 bucket.
func TestENG011_IPv6ForwardedClientsShareSlash64(t *testing.T) {
	s := newRateLimitServer(t, rlEnv{createPerMin: "2", redirectPerMin: "100", trustedCIDRs: "127.0.0.0/8"})
	proxy := s.client(t, nil)
	rlShortCode(t, s.create(t, proxy, "https://example.com/v6-1", xff("2001:db8:0:1::1")))
	rlShortCode(t, s.create(t, proxy, "https://example.com/v6-2", xff("2001:db8:0:1:ffff::2")))
	rlAssertLimited(t, "same /64", s.create(t, proxy, "https://example.com/v6-3", xff("[2001:db8:0:1::3]:443")))
	rlShortCode(t, s.create(t, proxy, "https://example.com/v6-other", xff("2001:db8:0:2::1")))

	rlShortCode(t, s.create(t, proxy, "https://example.com/v4-1", xff("198.51.100.20")))
	rlShortCode(t, s.create(t, proxy, "https://example.com/v4-2", xff("::ffff:198.51.100.20")))
	rlAssertLimited(t, "IPv4-mapped shares IPv4 bucket", s.create(t, proxy, "https://example.com/v4-3", xff("198.51.100.20")))
}

// TestENG011_ConcurrentCreatesAdmitExactlyLimit fires parallel creates from
// one source: exactly the limit succeed, the rest are 429, and the store
// holds exactly the admitted rows.
func TestENG011_ConcurrentCreatesAdmitExactlyLimit(t *testing.T) {
	const limit, total = 5, 40
	s := newRateLimitServer(t, rlEnv{createPerMin: fmt.Sprint(limit), redirectPerMin: "100"})
	c := s.client(t, nil)
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []rlResult
	)
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			payload := []byte(fmt.Sprintf(`{"url":"https://example.com/concurrent-%d"}`, i))
			req, err := http.NewRequest(http.MethodPost, s.srv.URL+rlCreatePath, bytes.NewReader(payload))
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			res, err := c.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			defer res.Body.Close()
			raw, _ := io.ReadAll(res.Body)
			mu.Lock()
			results = append(results, rlResult{status: res.StatusCode, header: res.Header, body: raw})
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	if t.Failed() {
		return
	}
	created, limited := 0, 0
	for i, r := range results {
		switch r.status {
		case http.StatusCreated:
			created++
		case http.StatusTooManyRequests:
			limited++
			rlAssertLimited(t, fmt.Sprintf("concurrent result #%d", i), r)
		default:
			t.Fatalf("unexpected status %d body %s", r.status, r.body)
		}
	}
	if created != limit || limited != total-limit {
		t.Fatalf("created %d limited %d, want %d and %d", created, limited, limit, total-limit)
	}
	if n := countRows(t, s.dbPath); n != limit {
		t.Fatalf("rows %d, want %d", n, limit)
	}
}

// TestENG011_LogsDoNotContainClientAddresses drives allowed and rejected
// traffic through direct, spoofed, and trusted-forwarded paths, then checks
// the captured observability and server logs: rate_limited failures are
// recorded by error_type, and no peer or forwarded address appears.
func TestENG011_LogsDoNotContainClientAddresses(t *testing.T) {
	s := newRateLimitServer(t, rlEnv{createPerMin: "1", redirectPerMin: "1", trustedCIDRs: "127.0.0.1/32"})
	proxy := s.client(t, nil)
	fwd := xff("198.51.100.77, 203.0.113.77")
	code := rlShortCode(t, s.create(t, proxy, "https://example.com/logs", fwd))
	rlAssertLimited(t, "create over quota", s.create(t, proxy, "https://example.com/logs-2", fwd))
	if res := s.redirect(t, proxy, code, fwd); res.status != http.StatusFound {
		t.Fatalf("redirect status %d", res.status)
	}
	rlAssertLimited(t, "redirect over quota", s.redirect(t, proxy, code, fwd))
	v6 := xff("2001:db8:77::77")
	rlShortCode(t, s.create(t, proxy, "https://example.com/logs-v6", v6))
	rlAssertLimited(t, "v6 create over quota", s.create(t, proxy, "https://example.com/logs-v6-2", v6))
	direct := http.Header{}
	direct.Set("X-Real-IP", "192.0.2.77")
	rlShortCode(t, s.create(t, proxy, "https://example.com/logs-real", direct))

	out := s.logs.String()
	var createLimited, redirectLimited int
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line == "" {
			continue
		}
		var rec map[string]any
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			continue
		}
		if rec["error_type"] == "rate_limited" {
			switch rec["event"] {
			case "create_failure":
				createLimited++
			case "redirect_failure":
				redirectLimited++
			}
		}
	}
	if createLimited != 2 || redirectLimited != 1 {
		t.Fatalf("rate_limited log events create=%d redirect=%d, want 2 and 1", createLimited, redirectLimited)
	}
	for _, needle := range []string{
		"127.0.0.1", "198.51.100.", "203.0.113.", "192.0.2.", "2001:db8",
		"X-Forwarded-For", "X-Real-Ip", "remote_addr",
	} {
		if strings.Contains(out, needle) {
			t.Fatalf("captured logs contain a forbidden value (%s)", needleLabel(needle))
		}
	}
}

// needleLabel keeps assertion output free of the address itself.
func needleLabel(needle string) string {
	switch {
	case strings.Contains(needle, "."), strings.Contains(needle, ":"):
		return "address"
	default:
		return needle
	}
}
