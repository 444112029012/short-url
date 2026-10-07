package httpapi

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/444112029012/short-url/internal/health"
)

// NewRouter registers health, create, stats, and redirect.
//
// TODO(ENG-016): add the SEC-008 header subset (nosniff; HSTS only when the
// listener is HTTPS). Not part of this batch.
func NewRouter(adapter *HTTPAPIAdapter) http.Handler {
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
	return r
}
