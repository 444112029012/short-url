package ratelimit

import "testing"

func TestGuardAlwaysAllows(t *testing.T) {
	g := NewGuard()
	if err := g.CheckCreate("203.0.113.5:1234"); err != nil {
		t.Fatal(err)
	}
	if err := g.CheckRedirect(""); err != nil {
		t.Fatal(err)
	}
}
