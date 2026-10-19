package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestUT_OBS_01_EachHookEmitsOneRecord(t *testing.T) {
	cases := []struct {
		name    string
		call    func(*Hooks)
		event   string
		field   string
		value   string
		counter func(*Hooks) uint64
	}{
		{"create success", func(h *Hooks) { h.OnCreateSuccess("Ab12Cd39") }, "create_success", "short_code", "Ab12Cd39", func(h *Hooks) uint64 { return h.createSuccess.Load() }},
		{"redirect success", func(h *Hooks) { h.OnRedirectSuccess("Zz98Yy76") }, "redirect_success", "short_code", "Zz98Yy76", func(h *Hooks) uint64 { return h.redirectSuccess.Load() }},
		{"stats success", func(h *Hooks) { h.OnStatsSuccess("Qw12Er34") }, "stats_success", "short_code", "Qw12Er34", func(h *Hooks) uint64 { return h.statsSuccess.Load() }},
		{"create failure", func(h *Hooks) { h.OnCreateFailure("invalid_url") }, "create_failure", "error_type", "invalid_url", func(h *Hooks) uint64 { return h.createFailure.Load() }},
		{"redirect failure", func(h *Hooks) { h.OnRedirectFailure("not_found") }, "redirect_failure", "error_type", "not_found", func(h *Hooks) uint64 { return h.redirectFailure.Load() }},
		{"stats failure", func(h *Hooks) { h.OnStatsFailure("internal_error") }, "stats_failure", "error_type", "internal_error", func(h *Hooks) uint64 { return h.statsFailure.Load() }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, buf := newBufHooks()
			tc.call(h)
			rec := singleRecord(t, buf.String())
			if rec["msg"] != tc.event || rec["event"] != tc.event || rec["level"] != "INFO" {
				t.Fatalf("record %v", rec)
			}
			if rec[tc.field] != tc.value {
				t.Fatalf("field %s = %v", tc.field, rec[tc.field])
			}
			if tc.counter(h) != 1 {
				t.Fatalf("counter %d", tc.counter(h))
			}
		})
	}
}

func TestHookCountersIncrement(t *testing.T) {
	h, _ := newBufHooks()
	for i := 0; i < 3; i++ {
		h.OnCreateSuccess("Ab12Cd39")
		h.OnCreateFailure("rate_limited")
		h.OnRedirectSuccess("Ab12Cd39")
		h.OnRedirectFailure("not_found")
		h.OnStatsSuccess("Ab12Cd39")
		h.OnStatsFailure("invalid_url")
	}
	if h.createSuccess.Load() != 3 || h.createFailure.Load() != 3 ||
		h.redirectSuccess.Load() != 3 || h.redirectFailure.Load() != 3 ||
		h.statsSuccess.Load() != 3 || h.statsFailure.Load() != 3 {
		t.Fatalf("counters %+v", []uint64{
			h.createSuccess.Load(), h.createFailure.Load(),
			h.redirectSuccess.Load(), h.redirectFailure.Load(),
			h.statsSuccess.Load(), h.statsFailure.Load(),
		})
	}
}

func TestFailureLogsOnlyErrorType(t *testing.T) {
	h, buf := newBufHooks()
	h.OnCreateFailure("invalid_url")
	h.OnRedirectFailure("not_found")
	h.OnStatsFailure("rate_limited")
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("lines %d", len(lines))
	}
	for _, line := range lines {
		rec := decodeRecord(t, line)
		if _, ok := rec["short_code"]; ok {
			t.Fatalf("failure record has short_code: %s", line)
		}
		if _, ok := rec["error_type"]; !ok {
			t.Fatalf("missing error_type: %s", line)
		}
		for _, key := range []string{"long_url", "authorization", "cookie", "ip", "remote_addr", "body"} {
			if _, ok := rec[key]; ok {
				t.Fatalf("unexpected %s in %s", key, line)
			}
		}
		allowed := map[string]bool{"time": true, "level": true, "msg": true, "event": true, "error_type": true}
		for key := range rec {
			if !allowed[key] {
				t.Fatalf("extra field %s in %s", key, line)
			}
		}
	}
}

func TestHooksDoNotPanic(t *testing.T) {
	h := NewHooksWithLogger(nil)
	h.OnCreateSuccess("Ab12Cd39")
	h.OnCreateFailure("invalid_url")
	h.OnRedirectSuccess("Ab12Cd39")
	h.OnRedirectFailure("not_found")
	h.OnStatsSuccess("Ab12Cd39")
	h.OnStatsFailure("internal_error")
	if h.createSuccess.Load() != 1 || h.statsFailure.Load() != 1 {
		t.Fatal("nil logger dropped the counter")
	}

	panicking := NewHooksWithLogger(slog.New(panicHandler{}))
	panicking.OnCreateSuccess("Ab12Cd39")
	panicking.OnCreateFailure("internal_error")
	if panicking.createSuccess.Load() != 1 || panicking.createFailure.Load() != 1 {
		t.Fatal("panic in the handler aborted the counter")
	}
}

type panicHandler struct{}

func (panicHandler) Enabled(context.Context, slog.Level) bool { return true }
func (panicHandler) Handle(context.Context, slog.Record) error {
	panic("logger failed")
}
func (panicHandler) WithAttrs([]slog.Attr) slog.Handler { return panicHandler{} }
func (panicHandler) WithGroup(string) slog.Handler      { return panicHandler{} }

func newBufHooks() (*Hooks, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	return NewHooksWithLogger(slog.New(slog.NewJSONHandler(buf, nil))), buf
}

func singleRecord(t *testing.T, raw string) map[string]any {
	t.Helper()
	raw = strings.TrimSpace(raw)
	if strings.Contains(raw, "\n") {
		t.Fatalf("expected one record, got %q", raw)
	}
	return decodeRecord(t, raw)
}

func decodeRecord(t *testing.T, raw string) map[string]any {
	t.Helper()
	var rec map[string]any
	if err := json.Unmarshal([]byte(raw), &rec); err != nil {
		t.Fatalf("json %s: %s", err, raw)
	}
	return rec
}
