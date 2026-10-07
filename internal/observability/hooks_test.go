package observability

import "testing"

func TestHooksDoNotPanic(t *testing.T) {
	h := NewHooks()
	h.OnCreateSuccess("Ab12Cd39")
	h.OnCreateFailure("invalid_url")
	h.OnRedirectSuccess("Ab12Cd39")
	h.OnRedirectFailure("not_found")
	h.OnStatsSuccess("Ab12Cd39")
	h.OnStatsFailure("internal_error")
}
