package httpapi_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"unicode/utf8"

	"github.com/444112029012/short-url/internal/observability"
	"github.com/444112029012/short-url/internal/ratelimit"
)

// createBodyLimit is the DR-G4-01 create request body limit in bytes.
const createBodyLimit = 32768

// paddedCreateBody returns a valid create JSON body of exactly size bytes:
// a small object followed by trailing JSON whitespace.
func paddedCreateBody(t *testing.T, size int) []byte {
	t.Helper()
	obj := `{"url":"https://example.com/body-limit"}`
	if size < len(obj) {
		t.Fatalf("size %d too small", size)
	}
	body := []byte(obj + strings.Repeat(" ", size-len(obj)))
	if len(body) != size {
		t.Fatalf("built %d bytes, want %d", len(body), size)
	}
	return body
}

// countingReader counts bytes read from r. With a nil r it serves an endless
// stream of spaces.
type countingReader struct {
	r io.Reader
	n atomic.Int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	if c.r == nil {
		for i := range p {
			p[i] = ' '
		}
		c.n.Add(int64(len(p)))
		return len(p), nil
	}
	n, err := c.r.Read(p)
	c.n.Add(int64(n))
	return n, err
}

func postRawBody(srv http.Handler, body io.Reader, remote string) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/urls", body)
	req.Header.Set("Content-Type", "application/json")
	if remote != "" {
		req.RemoteAddr = remote
	}
	srv.ServeHTTP(res, req)
	return res
}

func TestCreateBodyExactly32768Accepted(t *testing.T) {
	srv := newServer(t)
	res := postRawBody(srv, bytes.NewReader(paddedCreateBody(t, createBodyLimit)), "")
	if res.Code != http.StatusCreated {
		t.Fatalf("32768-byte body: status %d body %s", res.Code, res.Body.String())
	}
	if got := decodeMap(t, res.Body.Bytes())["long_url"]; got != "https://example.com/body-limit" {
		t.Fatalf("long_url %v", got)
	}
}

func TestCreateBody32769RejectedInvalidURL(t *testing.T) {
	srv := newServer(t)
	res := postRawBody(srv, bytes.NewReader(paddedCreateBody(t, createBodyLimit+1)), "")
	assertError(t, res, http.StatusBadRequest, "invalid_url")
}

func TestCreateBodyReadIsBounded(t *testing.T) {
	srv := newServer(t)
	body := &countingReader{}
	res := postRawBody(srv, body, "")
	assertError(t, res, http.StatusBadRequest, "invalid_url")
	// SCR-024: the adapter reads at most limit+1 bytes; never unbounded.
	if got := body.n.Load(); got > createBodyLimit+1 {
		t.Fatalf("read %d bytes from endless body, want <= %d", got, createBodyLimit+1)
	}
}

func TestRateLimitedOversizedCreateGets429WithoutReadingBody(t *testing.T) {
	srv := newLimitedServer(t, 1, 120, "")
	const remote = "203.0.113.50:4000"
	first := postRawBody(srv, bytes.NewReader(paddedCreateBody(t, 64)), remote)
	if first.Code != http.StatusCreated {
		t.Fatalf("first create %d %s", first.Code, first.Body.String())
	}
	// SCR-023: the counter wraps the entire oversized body.
	body := &countingReader{r: bytes.NewReader(paddedCreateBody(t, createBodyLimit+1))}
	second := postRawBody(srv, body, remote)
	assertError(t, second, http.StatusTooManyRequests, "rate_limited")
	if got := body.n.Load(); got != 0 {
		t.Fatalf("rate-limited request read %d body bytes; limit check must precede body read", got)
	}
}

// TestCreateC4aSurrogatePairEscaped2048CharsAccepted covers Security opinion
// DR-G4-01 C4(a) / method spec §13.4.1 test (a) and SCR-022: the url totals
// 2048 runes, the "https://example.com/" prefix included, and every other
// rune is outside the BMP (U+1F600) and JSON-escaped as a \uXXXX\uXXXX
// surrogate pair. The body is over the old 8192 limit and within 32768, the
// create returns 201, and the stored long_url is exactly 2048 runes.
func TestCreateC4aSurrogatePairEscaped2048CharsAccepted(t *testing.T) {
	const prefix = "https://example.com/"
	n := 2048 - utf8.RuneCountInString(prefix)
	body := []byte(`{"url":"` + prefix + strings.Repeat(`\uD83D\uDE00`, n) + `"}`)
	if len(body) <= 8192 || len(body) > createBodyLimit {
		t.Fatalf("request body %d bytes, want > 8192 and <= %d", len(body), createBodyLimit)
	}
	want := prefix + strings.Repeat("\U0001F600", n)
	if got := utf8.RuneCountInString(want); got != 2048 {
		t.Fatalf("fixture url %d runes, want 2048", got)
	}

	store := openStore(t, filepath.Join(t.TempDir(), "c4a.db"))
	srv := routerWith(t, store, ratelimit.NewGuard(30, 120), observability.NewHooks())
	res := postRawBody(srv, bytes.NewReader(body), "")
	if res.Code != http.StatusCreated {
		t.Fatalf("C4(a) status %d body %s", res.Code, res.Body.String())
	}
	code, _ := decodeMap(t, res.Body.Bytes())["short_code"].(string)
	rec, err := store.FindByShortCode(context.Background(), code)
	if err != nil {
		t.Fatalf("read back %s: %v", code, err)
	}
	if got := utf8.RuneCountInString(rec.LongURL); got != 2048 {
		t.Fatalf("stored long_url %d runes, want 2048", got)
	}
	if rec.LongURL != want {
		t.Fatal("stored long_url differs from the decoded url")
	}
}
