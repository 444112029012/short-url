package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

// UT-HR-405-ALLOW (SEC-010 ENG-020): a registered path requested with a
// method it does not register must answer 405 with an Allow header that
// lists exactly the registered methods, each once. chi v5.3.0 appended the
// same method again when a static route and the /{shortCode} param route
// both matched (e.g. /health), yielding "Allow: GET" twice.
func TestUTHR405Allow(t *testing.T) {
	srv := newServer(t)
	cases := []struct {
		name, method, path string
		allow              []string
	}{
		{"health static+param overlap", http.MethodPost, "/health", []string{http.MethodGet}},
		{"health PUT", http.MethodPut, "/health", []string{http.MethodGet}},
		{"create route DELETE", http.MethodDelete, "/api/v1/urls", []string{http.MethodPost}},
		{"create route GET", http.MethodGet, "/api/v1/urls", []string{http.MethodPost}},
		{"stats route POST", http.MethodPost, "/api/v1/urls/Abc12345/stats", []string{http.MethodGet}},
		{"redirect route DELETE", http.MethodDelete, "/Abc12345", []string{http.MethodGet}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			srv.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
			if rec.Code != http.StatusMethodNotAllowed {
				t.Fatalf("%s %s status %d, want 405", tc.method, tc.path, rec.Code)
			}
			var got []string
			for _, line := range rec.Header().Values("Allow") {
				for _, m := range strings.Split(line, ",") {
					if m = strings.TrimSpace(m); m != "" {
						got = append(got, m)
					}
				}
			}
			sorted := slices.Clone(got)
			slices.Sort(sorted)
			if len(slices.Compact(slices.Clone(sorted))) != len(got) {
				t.Errorf("%s %s Allow has duplicates: %q", tc.method, tc.path, got)
			}
			want := slices.Clone(tc.allow)
			slices.Sort(want)
			if !slices.Equal(sorted, want) {
				t.Errorf("%s %s Allow %q (count %d), want %q (count %d)", tc.method, tc.path, got, len(got), tc.allow, len(tc.allow))
			}
			assertBaselineHeaders(t, rec)
		})
	}
}
