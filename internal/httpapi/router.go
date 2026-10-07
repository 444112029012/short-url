package httpapi

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/444112029012/short-url/internal/health"
)

// RouterOptions controls optional response headers.
// HSTS is off unless the process opts in; the middleware still sends it
// only for TLS requests.
type RouterOptions struct {
	HSTS bool
}

// NewRouter registers health, create, stats, and redirect with the SEC-008
// header subset. HSTS stays off.
func NewRouter(adapter *HTTPAPIAdapter) http.Handler {
	return NewRouterWithOptions(adapter, RouterOptions{})
}

// NewRouterWithOptions is NewRouter with an explicit HSTS opt-in.
func NewRouterWithOptions(adapter *HTTPAPIAdapter, opts RouterOptions) http.Handler {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, health.Status())
	})
	r.Post("/api/v1/urls", func(w http.ResponseWriter, req *http.Request) {
		adapter.HandleCreate(req).write(w)
	})
	r.Get("/api/v1/urls/{shortCode}/stats", func(w http.ResponseWriter, req *http.Request) {
		adapter.HandleStats(req).write(w)
	})
	r.Get("/{shortCode}", func(w http.ResponseWriter, req *http.Request) {
		adapter.HandleRedirect(req).write(w)
	})
	return securityHeaders(opts.HSTS, r)
}
