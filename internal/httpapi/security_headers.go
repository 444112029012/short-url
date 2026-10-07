package httpapi

import (
	"net/http"
	"strings"
)

const (
	headerNosniff    = "nosniff"
	headerNoReferrer = "no-referrer"
	headerFrameDeny  = "DENY"
	// API responses are not documents. frame-ancestors closes clickjacking
	// even for clients that ignore X-Frame-Options.
	headerCSP = "default-src 'none'; frame-ancestors 'none'"
	// SEC-008: max-age at least 31536000 when HTTPS is actually in use.
	headerHSTS = "max-age=31536000"
)

// securityHeaders applies the SEC-008 subset plus the API framing and
// referrer headers to every response, including 404 and 405.
// HSTS is emitted only when enabled and the request is TLS. Plain HTTP,
// including local development, never receives Strict-Transport-Security.
// Location and the status code are left untouched. Cache-Control: no-store
// is added only for application/json so a 302 is not given a cache policy.
func securityHeaders(hsts bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", headerNosniff)
		h.Set("Referrer-Policy", headerNoReferrer)
		h.Set("X-Frame-Options", headerFrameDeny)
		h.Set("Content-Security-Policy", headerCSP)
		if hsts && requestIsTLS(r) {
			h.Set("Strict-Transport-Security", headerHSTS)
		}
		next.ServeHTTP(&cacheControlWriter{ResponseWriter: w}, r)
	})
}

func requestIsTLS(r *http.Request) bool {
	return r != nil && r.TLS != nil
}

type cacheControlWriter struct {
	http.ResponseWriter
	wrote bool
}

func (w *cacheControlWriter) WriteHeader(status int) {
	if !w.wrote {
		ct := w.Header().Get("Content-Type")
		media := strings.TrimSpace(strings.Split(ct, ";")[0])
		if strings.EqualFold(media, "application/json") && w.Header().Get("Cache-Control") == "" {
			w.Header().Set("Cache-Control", "no-store")
		}
		w.wrote = true
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *cacheControlWriter) Write(b []byte) (int, error) {
	if !w.wrote {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}
