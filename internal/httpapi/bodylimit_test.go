package httpapi_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
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

// countingReader serves an endless stream of spaces and counts bytes read.
type countingReader struct{ n atomic.Int64 }

func (c *countingReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = ' '
	}
	c.n.Add(int64(len(p)))
	return len(p), nil
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
	// The adapter must stop reading just past the limit; never read unbounded.
	// Allow one read buffer of slack beyond limit+1 for the reader's chunking.
	if got := body.n.Load(); got > createBodyLimit+1+64*1024 {
		t.Fatalf("read %d bytes from endless body; read is not bounded", got)
	}
}

func TestRateLimitedOversizedCreateGets429WithoutReadingBody(t *testing.T) {
	srv := newLimitedServer(t, 1, 120, "")
	const remote = "203.0.113.50:4000"
	first := postRawBody(srv, bytes.NewReader(paddedCreateBody(t, 64)), remote)
	if first.Code != http.StatusCreated {
		t.Fatalf("first create %d %s", first.Code, first.Body.String())
	}
	body := &countingReader{}
	second := postRawBody(srv, io.MultiReader(bytes.NewReader(paddedCreateBody(t, createBodyLimit+1)), body), remote)
	assertError(t, second, http.StatusTooManyRequests, "rate_limited")
	if got := body.n.Load(); got != 0 {
		t.Fatalf("rate-limited request read %d body bytes; limit check must precede body read", got)
	}
}
