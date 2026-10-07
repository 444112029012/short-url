package httpapi

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/444112029012/short-url/internal/domain"
)

func TestToHTTPResponseCatalog(t *testing.T) {
	mapper := NewErrorMapper()
	cases := []struct {
		name    string
		err     *domain.AppError
		status  int
		typ     string
		message string
	}{
		{"invalid url", domain.InvalidURL(domain.MsgURLValidationFailed), 400, "invalid_url", domain.MsgURLValidationFailed},
		{"invalid code", domain.InvalidURL(domain.MsgInvalidShortCode), 400, "invalid_url", domain.MsgInvalidShortCode},
		{"not found", domain.NotFound(domain.MsgShortCodeNotFound), 404, "not_found", domain.MsgShortCodeNotFound},
		{"rate limited", domain.RateLimited(domain.MsgTooManyRequests), 429, "rate_limited", domain.MsgTooManyRequests},
		{"internal", domain.Internal(), 500, "internal_error", domain.MsgUnexpected},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp := mapper.ToHTTPResponse(tc.err)
			if resp.status != tc.status {
				t.Fatalf("UT-ERR-01 status %d", resp.status)
			}
			var body map[string]any
			if err := json.Unmarshal(resp.body, &body); err != nil {
				t.Fatal(err)
			}
			if len(body) != 1 {
				t.Fatalf("extra top-level keys: %v", body)
			}
			errObj, _ := body["error"].(map[string]any)
			if len(errObj) != 2 || errObj["type"] != tc.typ || errObj["message"] != tc.message {
				t.Fatalf("body %s", resp.body)
			}
		})
	}
}

func TestToHTTPResponseHidesInternals(t *testing.T) {
	mapper := NewErrorMapper()
	leaky := &domain.AppError{
		Type:    "sql_exception",
		Message: "panic goroutine stack /home/app/store.go:41 SELECT * FROM url_mapping",
	}
	resp := mapper.ToHTTPResponse(leaky)
	if resp.status != 500 {
		t.Fatalf("status %d", resp.status)
	}
	text := string(resp.body)
	for _, forbidden := range []string{"panic", "goroutine", "store.go", "SELECT", "/home/", "sql_exception"} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("UT-ERR-02 leaked %q in %s", forbidden, text)
		}
	}
	if !strings.Contains(text, `"type":"internal_error"`) || !strings.Contains(text, domain.MsgUnexpected) {
		t.Fatalf("body %s", text)
	}
}
