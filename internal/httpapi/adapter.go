package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/444112029012/short-url/internal/application"
	"github.com/444112029012/short-url/internal/domain"
	"github.com/444112029012/short-url/internal/ratelimit"
)

const maxCreateBody = 8192

// RateLimitChecker is the RateLimitGuard port used by the adapter.
type RateLimitChecker interface {
	CheckCreate(sourceID string) *domain.AppError
	CheckRedirect(sourceID string) *domain.AppError
}

// EventHooks is the ObservabilityHooks port used by the adapter.
type EventHooks interface {
	OnCreateSuccess(shortCode string)
	OnCreateFailure(errorType string)
	OnRedirectSuccess(shortCode string)
	OnRedirectFailure(errorType string)
	OnStatsSuccess(shortCode string)
	OnStatsFailure(errorType string)
}

// HTTPAPIAdapter implements HttpApiAdapter handleCreate, handleRedirect, and handleStats.
type HTTPAPIAdapter struct {
	app      *application.ShortURLApplicationService
	redirect *application.RedirectService
	limiter  RateLimitChecker
	errors   ErrorMapper
	obs      EventHooks
	trusted  []*net.IPNet
}

func NewHTTPAPIAdapter(
	app *application.ShortURLApplicationService,
	redirect *application.RedirectService,
	limiter RateLimitChecker,
	errors ErrorMapper,
	obs EventHooks,
	trusted []*net.IPNet,
) *HTTPAPIAdapter {
	return &HTTPAPIAdapter{
		app:      app,
		redirect: redirect,
		limiter:  limiter,
		errors:   errors,
		obs:      obs,
		trusted:  trusted,
	}
}

type createRequest struct {
	URL string `json:"url"`
}

type createResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
	LongURL   string `json:"long_url"`
}

type statsResponse struct {
	ShortCode  string `json:"short_code"`
	ClickCount int    `json:"click_count"`
}

// HandleCreate serves POST /api/v1/urls.
func (a *HTTPAPIAdapter) HandleCreate(r *http.Request) response {
	if appErr := a.limiter.CheckCreate(a.sourceID(r)); appErr != nil {
		a.obs.OnCreateFailure(string(appErr.Type))
		return a.errors.ToHTTPResponse(appErr)
	}
	longURL, appErr := decodeCreate(r)
	if appErr != nil {
		a.obs.OnCreateFailure(string(appErr.Type))
		return a.errors.ToHTTPResponse(appErr)
	}
	created, appErr := a.app.CreateShortURL(r.Context(), longURL)
	if appErr != nil {
		a.obs.OnCreateFailure(string(appErr.Type))
		return a.errors.ToHTTPResponse(appErr)
	}
	raw, err := json.Marshal(createResponse{
		ShortCode: created.ShortCode,
		ShortURL:  created.ShortURL,
		LongURL:   created.LongURL,
	})
	if err != nil {
		appErr = domain.Internal()
		a.obs.OnCreateFailure(string(appErr.Type))
		return a.errors.ToHTTPResponse(appErr)
	}
	a.obs.OnCreateSuccess(created.ShortCode)
	return response{status: http.StatusCreated, header: jsonHeader(), body: raw}
}

// HandleRedirect serves GET /{shortCode}.
// Location is only the stored long_url. Query and header values are ignored (SEC-014).
func (a *HTTPAPIAdapter) HandleRedirect(r *http.Request) response {
	if appErr := a.limiter.CheckRedirect(a.sourceID(r)); appErr != nil {
		a.obs.OnRedirectFailure(string(appErr.Type))
		return a.errors.ToHTTPResponse(appErr)
	}
	target, appErr := a.redirect.Redirect(r.Context(), chi.URLParam(r, "shortCode"))
	if appErr != nil {
		a.obs.OnRedirectFailure(string(appErr.Type))
		resp := a.errors.ToHTTPResponse(appErr)
		resp.header.Del("Location")
		return resp
	}
	a.obs.OnRedirectSuccess(chi.URLParam(r, "shortCode"))
	header := make(http.Header)
	header.Set("Location", target.LongURL)
	return response{status: http.StatusFound, header: header}
}

// HandleStats serves GET /api/v1/urls/{shortCode}/stats.
func (a *HTTPAPIAdapter) HandleStats(r *http.Request) response {
	stats, appErr := a.app.GetStats(r.Context(), chi.URLParam(r, "shortCode"))
	if appErr != nil {
		a.obs.OnStatsFailure(string(appErr.Type))
		return a.errors.ToHTTPResponse(appErr)
	}
	raw, err := json.Marshal(statsResponse{ShortCode: stats.ShortCode, ClickCount: stats.ClickCount})
	if err != nil {
		appErr = domain.Internal()
		a.obs.OnStatsFailure(string(appErr.Type))
		return a.errors.ToHTTPResponse(appErr)
	}
	a.obs.OnStatsSuccess(stats.ShortCode)
	return response{status: http.StatusOK, header: jsonHeader(), body: raw}
}

func (a *HTTPAPIAdapter) sourceID(r *http.Request) string {
	var trusted []*net.IPNet
	if a != nil {
		trusted = a.trusted
	}
	return ratelimit.ClientSource(r, trusted)
}

func decodeCreate(r *http.Request) (string, *domain.AppError) {
	if r == nil || !jsonContentType(r.Header.Get("Content-Type")) {
		return "", domain.InvalidURL(domain.MsgURLValidationFailed)
	}
	body := r.Body
	if body == nil {
		body = io.NopCloser(bytes.NewReader(nil))
	}
	raw, err := io.ReadAll(io.LimitReader(body, maxCreateBody+1))
	if err != nil || len(raw) == 0 || len(raw) > maxCreateBody {
		return "", domain.InvalidURL(domain.MsgURLValidationFailed)
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	var req createRequest
	if err := dec.Decode(&req); err != nil {
		return "", domain.InvalidURL(domain.MsgURLValidationFailed)
	}
	var extra json.RawMessage
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return "", domain.InvalidURL(domain.MsgURLValidationFailed)
	}
	return req.URL, nil
}

func jsonContentType(value string) bool {
	mediaType, _, err := mime.ParseMediaType(strings.TrimSpace(value))
	return err == nil && mediaType == "application/json"
}
