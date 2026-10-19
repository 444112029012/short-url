package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/444112029012/short-url/internal/application"
	"github.com/444112029012/short-url/internal/domain"
	"github.com/444112029012/short-url/internal/httpapi"
	"github.com/444112029012/short-url/internal/infra/sqlite"
	"github.com/444112029012/short-url/internal/observability"
	"github.com/444112029012/short-url/internal/ratelimit"
)

func TestIT_FLOW_01_CreateRedirectStats(t *testing.T) {
	// IT-FLOW-01 / E2E-01: create 201, redirect 302 with the stored Location, stats increments.
	srv, client, _ := newServer(t)
	longURL := "https://example.com/path?q=1"
	created := postURL(t, client, srv.URL, longURL)
	if created.StatusCode != http.StatusCreated {
		t.Fatalf("E2E-01 status %d", created.StatusCode)
	}
	body := decodeCreate(t, created)
	if body.LongURL != longURL || len(body.ShortCode) != 8 || body.ShortURL != "http://localhost:8080/"+body.ShortCode {
		t.Fatalf("create body %+v", body)
	}

	res, err := client.Get(srv.URL + "/" + body.ShortCode)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusFound {
		t.Fatalf("redirect %d", res.StatusCode)
	}
	if got := res.Header.Get("Location"); got != longURL {
		t.Fatalf("Location %q", got)
	}

	stats := getJSON(t, client, srv.URL+"/api/v1/urls/"+body.ShortCode+"/stats")
	if stats.StatusCode != http.StatusOK {
		t.Fatalf("stats %d", stats.StatusCode)
	}
	got := decodeStats(t, stats)
	if got.ShortCode != body.ShortCode || got.ClickCount != 1 {
		t.Fatalf("stats %+v", got)
	}
}

func TestE2E_02_InvalidURLDoesNotStore(t *testing.T) {
	// E2E-02 / IT-FLOW-02: empty and non-http(s) URLs are 400 and leave no row.
	srv, client, path := newServer(t)
	for _, raw := range []string{"", "ftp://example.com/a", "javascript:alert(1)"} {
		res := postURL(t, client, srv.URL, raw)
		assertError(t, res, http.StatusBadRequest, "invalid_url")
	}
	if n := countRows(t, path); n != 0 {
		t.Fatalf("IT-FLOW-02 stored %d rows", n)
	}
}

func TestE2E_02_Length2048Boundary(t *testing.T) {
	// E2E-02 / UT-VAL-04 / UT-VAL-05: 2048 runes succeed, 2049 is invalid_url and is not stored.
	srv, client, path := newServer(t)
	valid := "https://example.com/" + strings.Repeat("a", 2048-len("https://example.com/"))
	if len(valid) != 2048 {
		t.Fatalf("fixture %d", len(valid))
	}
	created := postURL(t, client, srv.URL, valid)
	if created.StatusCode != http.StatusCreated {
		t.Fatalf("2048 status %d", created.StatusCode)
	}
	body := decodeCreate(t, created)
	if body.LongURL != valid {
		t.Fatal("2048 long_url mismatch")
	}
	if n := countRows(t, path); n != 1 {
		t.Fatalf("rows after 2048: %d", n)
	}

	tooLong := valid + "a"
	rejected := postURL(t, client, srv.URL, tooLong)
	assertError(t, rejected, http.StatusBadRequest, "invalid_url")
	if n := countRows(t, path); n != 1 {
		t.Fatalf("2049 stored a row, count %d", n)
	}
}

func TestE2E_03_UnknownShortCode(t *testing.T) {
	// E2E-03: unknown code is 404 for redirect and stats, with no Location.
	srv, client, _ := newServer(t)
	res, err := client.Get(srv.URL + "/NoSuch01")
	if err != nil {
		t.Fatal(err)
	}
	if loc := res.Header.Get("Location"); loc != "" {
		t.Fatalf("404 Location %q", loc)
	}
	assertError(t, res, http.StatusNotFound, "not_found")

	stats := getJSON(t, client, srv.URL+"/api/v1/urls/NoSuch01/stats")
	assertError(t, stats, http.StatusNotFound, "not_found")
}

func TestE2E_04_LocationIgnoresQueryRewrite(t *testing.T) {
	// E2E-04: a query or header cannot replace the stored Location.
	srv, client, _ := newServer(t)
	longURL := "https://example.com/stored"
	created := postURL(t, client, srv.URL, longURL)
	body := decodeCreate(t, created)
	req, err := http.NewRequest(http.MethodGet, srv.URL+"/"+body.ShortCode+"?url=https://evil.example&long_url=https://evil.example", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Redirect-Target", "https://evil.example")
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusFound || res.Header.Get("Location") != longURL {
		t.Fatalf("E2E-04 %d Location %q", res.StatusCode, res.Header.Get("Location"))
	}
}

func TestE2E_06_MalformedShortCode(t *testing.T) {
	// E2E-06: a malformed short code is 400 and does not increment a real code.
	srv, client, path := newServer(t)
	created := postURL(t, client, srv.URL, "https://example.com/keep")
	body := decodeCreate(t, created)

	res, err := client.Get(srv.URL + "/bad")
	if err != nil {
		t.Fatal(err)
	}
	assertError(t, res, http.StatusBadRequest, "invalid_url")

	stats := getJSON(t, client, srv.URL+"/api/v1/urls/bad/stats")
	assertError(t, stats, http.StatusBadRequest, "invalid_url")

	if n := countRows(t, path); n != 1 {
		t.Fatalf("rows %d", n)
	}
	got := decodeStats(t, getJSON(t, client, srv.URL+"/api/v1/urls/"+body.ShortCode+"/stats"))
	if got.ClickCount != 0 {
		t.Fatalf("malformed request counted a click: %d", got.ClickCount)
	}
}

func TestCreateRejectsNonJSONContentType(t *testing.T) {
	srv, client, path := newServer(t)
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/urls", strings.NewReader(`{"url":"https://example.com"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/plain")
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assertError(t, res, http.StatusBadRequest, "invalid_url")
	if n := countRows(t, path); n != 0 {
		t.Fatalf("non-json stored %d rows", n)
	}
}

func TestIT_STORE_02_ConcurrentRedirectsExactCount(t *testing.T) {
	// IT-STORE-02: N concurrent successful redirects leave click_count == N.
	const n = 24
	srv, client, _ := newServer(t)
	longURL := "https://example.com/concurrent"
	created := postURL(t, client, srv.URL, longURL)
	body := decodeCreate(t, created)

	var wg sync.WaitGroup
	errCh := make(chan string, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := client.Get(srv.URL + "/" + body.ShortCode)
			if err != nil {
				errCh <- err.Error()
				return
			}
			defer res.Body.Close()
			_, _ = io.Copy(io.Discard, res.Body)
			if res.StatusCode != http.StatusFound || res.Header.Get("Location") != longURL {
				errCh <- res.Status + " " + res.Header.Get("Location")
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for msg := range errCh {
		t.Error(msg)
	}
	if t.Failed() {
		return
	}
	got := decodeStats(t, getJSON(t, client, srv.URL+"/api/v1/urls/"+body.ShortCode+"/stats"))
	if got.ClickCount != n {
		t.Fatalf("IT-STORE-02 click_count %d, want %d", got.ClickCount, n)
	}
}

func TestDatabaseFileMode0600(t *testing.T) {
	_, _, path := newServer(t)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode %o", info.Mode().Perm())
	}
}

func TestIT_STORE_01_SchemaOmitsVisitorFields(t *testing.T) {
	// IT-STORE-01: url_mapping has no visitor_ip or user_agent column.
	_, _, path := newServer(t)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(`SELECT name FROM pragma_table_info('url_mapping')`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	got := map[string]bool{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		got[name] = true
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"visitor_ip", "user_agent"} {
		if got[forbidden] {
			t.Fatalf("SEC-004 column %s", forbidden)
		}
	}
	for _, required := range []string{"short_code", "long_url", "click_count", "created_at"} {
		if !got[required] {
			t.Fatalf("missing %s in %v", required, got)
		}
	}
	if len(got) != 4 {
		t.Fatalf("columns %v", got)
	}
}

type createBody struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
	LongURL   string `json:"long_url"`
}

type statsBody struct {
	ShortCode  string `json:"short_code"`
	ClickCount int    `json:"click_count"`
}

func newServer(t *testing.T) (*httptest.Server, *http.Client, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "shorturl.db")
	store, err := sqlite.Open(path)
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
	adapter := httpapi.NewHTTPAPIAdapter(appSvc, redirectSvc, ratelimit.NewGuard(), httpapi.NewErrorMapper(), observability.NewHooks())
	srv := httptest.NewServer(httpapi.NewRouter(adapter))
	t.Cleanup(func() {
		srv.Close()
		_ = store.Close()
	})
	client := srv.Client()
	client.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return srv, client, path
}

func postURL(t *testing.T, client *http.Client, base, rawURL string) *http.Response {
	t.Helper()
	payload, err := json.Marshal(map[string]string{"url": rawURL})
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.Post(base+"/api/v1/urls", "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func getJSON(t *testing.T, client *http.Client, url string) *http.Response {
	t.Helper()
	res, err := client.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func decodeCreate(t *testing.T, res *http.Response) createBody {
	t.Helper()
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create %d %s", res.StatusCode, raw)
	}
	var body createBody
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatal(err)
	}
	return body
}

func decodeStats(t *testing.T, res *http.Response) statsBody {
	t.Helper()
	defer res.Body.Close()
	var body statsBody
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	return body
}

func assertError(t *testing.T, res *http.Response, status int, typ string) {
	t.Helper()
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != status {
		t.Fatalf("status %d body %s", res.StatusCode, raw)
	}
	var payload struct {
		Error struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("json %s: %s", err, raw)
	}
	if payload.Error.Type != typ || payload.Error.Message == "" {
		t.Fatalf("error body %s", raw)
	}
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
