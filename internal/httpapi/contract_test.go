package httpapi_test

// 契約測試 CT-*（ENG-013；關閉 DEF-001／G4-QA-R1）。
// 依安全決策不引入第三方 schema 驗證器：以 stdlib 手寫斷言，
// 期望值硬編為常數並註明 docs/04-api/openapi.yaml 行號。
// 若 openapi.yaml 變更，須同步更新本檔常數（殘餘風險：常數漂移）。

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/444112029012/short-url/internal/observability"
	"github.com/444112029012/short-url/internal/ratelimit"
)

const (
	// openapi.yaml L63, L72, L82: 成功與錯誤回應 media type。
	ctJSONMediaType = "application/json"
	// openapi.yaml L11: X-Content-Type-Options: nosniff（SEC-008）。
	ctNosniff = "nosniff"
	// openapi.yaml L69: createShortUrl 成功 201。
	ctStatusCreated = http.StatusCreated
	// openapi.yaml L79: createShortUrl 驗證失敗 400。
	ctStatusBadRequest = http.StatusBadRequest
	// openapi.yaml L123: redirectShortCode 成功 302。
	ctStatusFound = http.StatusFound
	// openapi.yaml L142: 短碼不存在 404。
	ctStatusNotFound = http.StatusNotFound
	// openapi.yaml L152: 超過速率限制 429。
	ctStatusTooMany = http.StatusTooManyRequests
	// openapi.yaml L180: getUrlStats 成功 200。
	ctStatusOK = http.StatusOK
	// openapi.yaml L101: 內部錯誤 500。
	ctStatusInternal = http.StatusInternalServerError
)

// openapi.yaml L228, L253, L270: short_code pattern。
var ctShortCodePattern = regexp.MustCompile(`^[A-Za-z0-9]{8}$`)

// openapi.yaml L288: ErrorResponse.error.type enum 只有四值。
var ctErrorTypeEnum = map[string]bool{
	"invalid_url":    true,
	"not_found":      true,
	"rate_limited":   true,
	"internal_error": true,
}

// openapi.yaml L248-L261: CreateUrlResponse required＋additionalProperties:false（L249）。
var ctCreateFields = []string{"long_url", "short_code", "short_url"}

// openapi.yaml L265-L272: UrlStatsResponse required＋additionalProperties:false（L266）。
var ctStatsFields = []string{"click_count", "short_code"}

// openapi.yaml L278-L284: ErrorResponse {error:{type,message}}，兩層 additionalProperties:false。
var (
	ctErrorTopFields   = []string{"error"}
	ctErrorInnerFields = []string{"message", "type"}
)

// openapi.yaml L259-L261: long_url maxLength 2048。
const ctLongURLMaxLen = 2048

func ctKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func ctAssertExactFields(t *testing.T, what string, m map[string]any, want []string) {
	t.Helper()
	got := ctKeys(m)
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("%s fields = %v, want exactly %v (additionalProperties:false)", what, got, want)
	}
}

func ctAssertJSONResponse(t *testing.T, res *httptest.ResponseRecorder, status int) map[string]any {
	t.Helper()
	if res.Code != status {
		t.Fatalf("status = %d, want %d; body %s", res.Code, status, res.Body.String())
	}
	ct := res.Header().Get("Content-Type")
	if media := strings.TrimSpace(strings.Split(ct, ";")[0]); media != ctJSONMediaType {
		t.Fatalf("Content-Type = %q, want %s", ct, ctJSONMediaType)
	}
	if got := res.Header().Get("X-Content-Type-Options"); got != ctNosniff {
		t.Fatalf("X-Content-Type-Options = %q, want %q", got, ctNosniff)
	}
	return decodeMap(t, res.Body.Bytes())
}

func ctAssertString(t *testing.T, m map[string]any, key string) string {
	t.Helper()
	v, ok := m[key].(string)
	if !ok {
		t.Fatalf("%s type = %T, want string", key, m[key])
	}
	return v
}

func ctAssertURI(t *testing.T, key, v string) {
	t.Helper()
	u, err := url.Parse(v)
	if err != nil || u.Scheme == "" || u.Host == "" {
		t.Fatalf("%s = %q is not an absolute uri (format: uri)", key, v)
	}
}

// ctAssertErrorResponse 驗證 ErrorResponse schema（openapi.yaml L276-L292）。
func ctAssertErrorResponse(t *testing.T, res *httptest.ResponseRecorder, status int, wantType string) {
	t.Helper()
	body := ctAssertJSONResponse(t, res, status)
	ctAssertExactFields(t, "ErrorResponse", body, ctErrorTopFields)
	inner, ok := body["error"].(map[string]any)
	if !ok {
		t.Fatalf("error type = %T, want object", body["error"])
	}
	ctAssertExactFields(t, "ErrorResponse.error", inner, ctErrorInnerFields)
	typ := ctAssertString(t, inner, "type")
	if !ctErrorTypeEnum[typ] {
		t.Fatalf("error.type = %q not in enum (openapi.yaml L288)", typ)
	}
	if typ != wantType {
		t.Fatalf("error.type = %q, want %q", typ, wantType)
	}
	if msg := ctAssertString(t, inner, "message"); msg == "" {
		t.Fatal("error.message empty")
	}
	if res.Header().Get("Location") != "" {
		t.Fatal("error response must not set Location")
	}
}

func ctCreate(t *testing.T, srv http.Handler) string {
	t.Helper()
	res := postURL(t, srv, "https://example.com/path?q=1")
	if res.Code != ctStatusCreated {
		t.Fatalf("setup create %d %s", res.Code, res.Body.String())
	}
	return decodeMap(t, res.Body.Bytes())["short_code"].(string)
}

func TestCT_CREATE_01_CreateReturns201CreateUrlResponse(t *testing.T) {
	srv := newServer(t)
	longURL := "https://example.com/path?q=1"
	res := postURL(t, srv, longURL)
	body := ctAssertJSONResponse(t, res, ctStatusCreated)
	ctAssertExactFields(t, "CreateUrlResponse", body, ctCreateFields)
	code := ctAssertString(t, body, "short_code")
	if !ctShortCodePattern.MatchString(code) {
		t.Fatalf("short_code %q does not match pattern (openapi.yaml L253)", code)
	}
	shortURL := ctAssertString(t, body, "short_url")
	ctAssertURI(t, "short_url", shortURL)
	if !strings.HasSuffix(shortURL, "/"+code) {
		t.Fatalf("short_url %q must be base URL + short_code (openapi.yaml L257)", shortURL)
	}
	gotLong := ctAssertString(t, body, "long_url")
	ctAssertURI(t, "long_url", gotLong)
	if gotLong != longURL || len([]rune(gotLong)) > ctLongURLMaxLen {
		t.Fatalf("long_url = %q, want %q", gotLong, longURL)
	}
}

func TestCT_CREATE_02_InvalidURLReturns400InvalidURL(t *testing.T) {
	srv := newServer(t)
	for _, bad := range []string{"", "ftp://example.com/x", "not a url", "https://example.com/" + strings.Repeat("a", ctLongURLMaxLen)} {
		res := postURL(t, srv, bad)
		ctAssertErrorResponse(t, res, ctStatusBadRequest, "invalid_url")
	}
}

func TestCT_REDIR_01_RedirectReturns302WithLocation(t *testing.T) {
	srv := newServer(t)
	code := ctCreate(t, srv)
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/"+code+"?url=https://evil.example", nil))
	if res.Code != ctStatusFound {
		t.Fatalf("status = %d, want %d", res.Code, ctStatusFound)
	}
	loc := res.Header().Get("Location") // openapi.yaml L126-L130: Location string/uri
	ctAssertURI(t, "Location", loc)
	if loc != "https://example.com/path?q=1" {
		t.Fatalf("Location = %q, want stored long_url (SEC-014, openapi.yaml L13)", loc)
	}
	if got := res.Header().Get("X-Content-Type-Options"); got != ctNosniff {
		t.Fatalf("X-Content-Type-Options = %q", got)
	}
}

func TestCT_REDIR_02_UnknownCodeReturns404NotFound(t *testing.T) {
	srv := newServer(t)
	res := getJSON(t, srv, "/Zz99Yy88")
	ctAssertErrorResponse(t, res, ctStatusNotFound, "not_found")
}

func TestCT_STATS_01_StatsReturns200UrlStatsResponse(t *testing.T) {
	srv := newServer(t)
	code := ctCreate(t, srv)
	redirect := httptest.NewRecorder()
	srv.ServeHTTP(redirect, httptest.NewRequest(http.MethodGet, "/"+code, nil))
	res := getJSON(t, srv, "/api/v1/urls/"+code+"/stats")
	body := ctAssertJSONResponse(t, res, ctStatusOK)
	ctAssertExactFields(t, "UrlStatsResponse", body, ctStatsFields)
	if got := ctAssertString(t, body, "short_code"); got != code || !ctShortCodePattern.MatchString(got) {
		t.Fatalf("short_code = %q, want %q (openapi.yaml L270)", got, code)
	}
	n, ok := body["click_count"].(float64)
	if !ok || n != float64(int64(n)) || n < 0 { // openapi.yaml L272-L273: integer, minimum 0
		t.Fatalf("click_count = %v (%T), want integer >= 0", body["click_count"], body["click_count"])
	}
	if n != 1 {
		t.Fatalf("click_count = %v, want 1 after one redirect", n)
	}
}

func TestCT_ERR_01_AllErrorsUseEnumAndNoExtraFields(t *testing.T) {
	t.Run("invalid_url/create", func(t *testing.T) {
		ctAssertErrorResponse(t, postURL(t, newServer(t), "ftp://x"), ctStatusBadRequest, "invalid_url")
	})
	t.Run("invalid_url/redirect_bad_code", func(t *testing.T) {
		ctAssertErrorResponse(t, getJSON(t, newServer(t), "/bad!"), ctStatusBadRequest, "invalid_url")
	})
	t.Run("invalid_url/stats_bad_code", func(t *testing.T) {
		ctAssertErrorResponse(t, getJSON(t, newServer(t), "/api/v1/urls/short/stats"), ctStatusBadRequest, "invalid_url")
	})
	t.Run("not_found/redirect", func(t *testing.T) {
		ctAssertErrorResponse(t, getJSON(t, newServer(t), "/Ab12Cd39"), ctStatusNotFound, "not_found")
	})
	t.Run("not_found/stats", func(t *testing.T) {
		ctAssertErrorResponse(t, getJSON(t, newServer(t), "/api/v1/urls/Ab12Cd39/stats"), ctStatusNotFound, "not_found")
	})
	t.Run("rate_limited/create", func(t *testing.T) {
		srv := routerWith(t, openStore(t, filepath.Join(t.TempDir(), "rl.db")), blockingLimiter{}, observability.NewHooks())
		ctAssertErrorResponse(t, postURL(t, srv, "https://example.com"), ctStatusTooMany, "rate_limited")
	})
	t.Run("rate_limited/redirect", func(t *testing.T) {
		srv := routerWith(t, openStore(t, filepath.Join(t.TempDir(), "rl.db")), blockingLimiter{}, observability.NewHooks())
		ctAssertErrorResponse(t, getJSON(t, srv, "/Ab12Cd39"), ctStatusTooMany, "rate_limited")
	})
	t.Run("internal_error/create", func(t *testing.T) {
		store := openStore(t, filepath.Join(t.TempDir(), "closed.db"))
		if err := store.Close(); err != nil {
			t.Fatal(err)
		}
		srv := routerWith(t, store, ratelimit.NewGuard(30, 120), observability.NewHooks())
		ctAssertErrorResponse(t, postURL(t, srv, "https://example.com/closed"), ctStatusInternal, "internal_error")
	})
}
