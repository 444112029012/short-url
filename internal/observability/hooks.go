package observability

// Hooks is a no-op stand-in for ObservabilityHooks.
//
// TODO(ENG-012): record create/redirect/stats success and failure with the stable
// error type and short code only. Do not log Authorization, Cookie, or client IP (SEC-005).
type Hooks struct{}

func NewHooks() *Hooks { return &Hooks{} }

func (h *Hooks) OnCreateSuccess(shortCode string)   { _ = shortCode }
func (h *Hooks) OnCreateFailure(errorType string)   { _ = errorType }
func (h *Hooks) OnRedirectSuccess(shortCode string) { _ = shortCode }
func (h *Hooks) OnRedirectFailure(errorType string) { _ = errorType }
func (h *Hooks) OnStatsSuccess(shortCode string)    { _ = shortCode }
func (h *Hooks) OnStatsFailure(errorType string)    { _ = errorType }
