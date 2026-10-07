package observability

import (
	"context"
	"log/slog"
	"os"
	"sync/atomic"
)

// Hooks records create, redirect, and stats outcomes (NFR-006).
// Each method logs one JSON record and increments an in-process counter.
// A logging failure is swallowed: the hook never panics and never returns
// an error to the caller (method spec §11).
//
// Success records carry event and short_code. Failure records carry event
// and error_type only. Authorization, Cookie, request bodies, long URLs,
// and client IPs are not accepted by these methods and are not logged (SEC-005).
type Hooks struct {
	logger *slog.Logger

	createSuccess   atomic.Uint64
	createFailure   atomic.Uint64
	redirectSuccess atomic.Uint64
	redirectFailure atomic.Uint64
	statsSuccess    atomic.Uint64
	statsFailure    atomic.Uint64
}

// NewHooks logs JSON lines to stderr.
func NewHooks() *Hooks {
	return NewHooksWithLogger(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
}

// NewHooksWithLogger uses logger for event lines. A nil logger still counts
// events and writes nothing.
func NewHooksWithLogger(logger *slog.Logger) *Hooks {
	return &Hooks{logger: logger}
}

// OnCreateSuccess records one create success for shortCode.
func (h *Hooks) OnCreateSuccess(shortCode string) {
	h.emit(&h.createSuccess, "create_success", slog.String("short_code", shortCode))
}

// OnCreateFailure records one create failure. errorType is a stable code.
func (h *Hooks) OnCreateFailure(errorType string) {
	h.emit(&h.createFailure, "create_failure", slog.String("error_type", errorType))
}

// OnRedirectSuccess records one redirect success for shortCode.
func (h *Hooks) OnRedirectSuccess(shortCode string) {
	h.emit(&h.redirectSuccess, "redirect_success", slog.String("short_code", shortCode))
}

// OnRedirectFailure records one redirect failure. errorType is a stable code.
func (h *Hooks) OnRedirectFailure(errorType string) {
	h.emit(&h.redirectFailure, "redirect_failure", slog.String("error_type", errorType))
}

// OnStatsSuccess records one stats success for shortCode.
func (h *Hooks) OnStatsSuccess(shortCode string) {
	h.emit(&h.statsSuccess, "stats_success", slog.String("short_code", shortCode))
}

// OnStatsFailure records one stats failure. errorType is a stable code.
func (h *Hooks) OnStatsFailure(errorType string) {
	h.emit(&h.statsFailure, "stats_failure", slog.String("error_type", errorType))
}

func (h *Hooks) emit(counter *atomic.Uint64, event string, attr slog.Attr) {
	if h == nil || counter == nil {
		return
	}
	counter.Add(1)
	if h.logger == nil {
		return
	}
	defer func() { _ = recover() }()
	h.logger.LogAttrs(context.Background(), slog.LevelInfo, event,
		slog.String("event", event),
		attr,
	)
}
