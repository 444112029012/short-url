package ratelimit

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/444112029012/short-url/internal/domain"
)

func TestCheckCreateUnderAtOverLimit(t *testing.T) {
	g := NewGuard(3, 120)
	for i := 0; i < 3; i++ {
		if err := g.CheckCreate("203.0.113.5"); err != nil {
			t.Fatalf("UT-RL-01 under/at %d: %v", i, err)
		}
	}
	err := g.CheckCreate("203.0.113.5")
	assertRateLimited(t, err)
	if g.create["203.0.113.5"].count != 3 {
		t.Fatalf("over-limit consumed quota: %d", g.create["203.0.113.5"].count)
	}
}

func TestCheckRedirectUnderAtOverLimit(t *testing.T) {
	g := NewGuard(30, 2)
	if err := g.CheckRedirect("203.0.113.5"); err != nil {
		t.Fatal(err)
	}
	if err := g.CheckRedirect("203.0.113.5"); err != nil {
		t.Fatal(err)
	}
	assertRateLimited(t, g.CheckRedirect("203.0.113.5"))
	if g.redirect["203.0.113.5"].count != 2 {
		t.Fatalf("UT-RL-02 over-limit consumed quota: %d", g.redirect["203.0.113.5"].count)
	}
}

func TestWindowResetWithFakeClock(t *testing.T) {
	now := time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC)
	g := NewGuard(2, 1)
	g.now = func() time.Time { return now }

	if err := g.CheckCreate("203.0.113.9"); err != nil {
		t.Fatal(err)
	}
	if err := g.CheckCreate("203.0.113.9"); err != nil {
		t.Fatal(err)
	}
	assertRateLimited(t, g.CheckCreate("203.0.113.9"))

	now = now.Add(time.Minute)
	if err := g.CheckCreate("203.0.113.9"); err != nil {
		t.Fatalf("window did not reset: %v", err)
	}
	if err := g.CheckCreate("203.0.113.9"); err != nil {
		t.Fatal(err)
	}
	assertRateLimited(t, g.CheckCreate("203.0.113.9"))
}

func TestWindowResetRedirect(t *testing.T) {
	now := time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC)
	g := NewGuard(30, 1)
	g.now = func() time.Time { return now }
	if err := g.CheckRedirect("2001:db8::1"); err != nil {
		t.Fatal(err)
	}
	assertRateLimited(t, g.CheckRedirect("2001:db8::1"))
	now = now.Add(time.Minute - time.Nanosecond)
	assertRateLimited(t, g.CheckRedirect("2001:db8::1"))
	now = now.Add(time.Nanosecond)
	if err := g.CheckRedirect("2001:db8::1"); err != nil {
		t.Fatalf("redirect window did not reset: %v", err)
	}
}

func TestSourcesAndMethodsAreIndependent(t *testing.T) {
	g := NewGuard(1, 1)
	if err := g.CheckCreate("203.0.113.1"); err != nil {
		t.Fatal(err)
	}
	assertRateLimited(t, g.CheckCreate("203.0.113.1"))
	if err := g.CheckCreate("203.0.113.2"); err != nil {
		t.Fatal("other source was limited")
	}
	if err := g.CheckRedirect("203.0.113.1"); err != nil {
		t.Fatal("create quota affected redirect")
	}
	assertRateLimited(t, g.CheckRedirect("203.0.113.1"))
}

func TestEmptySourceSharesOneBucket(t *testing.T) {
	g := NewGuard(1, 1)
	if err := g.CheckCreate(""); err != nil {
		t.Fatal(err)
	}
	assertRateLimited(t, g.CheckCreate(""))
	if err := g.CheckCreate("203.0.113.5"); err != nil {
		t.Fatal("empty bucket collided with a real source")
	}
	if _, ok := g.create[emptySourceKey]; !ok {
		t.Fatal("empty source was not stored on the shared key")
	}
}

func TestGuardEvictsIdleSourcesAtCapacity(t *testing.T) {
	now := time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC)
	g := NewGuard(5, 5)
	g.maxSources = 2
	g.now = func() time.Time { return now }
	if err := g.CheckCreate("a"); err != nil {
		t.Fatal(err)
	}
	if err := g.CheckCreate("b"); err != nil {
		t.Fatal(err)
	}
	assertRateLimited(t, g.CheckCreate("c"))
	now = now.Add(time.Minute)
	if err := g.CheckCreate("c"); err != nil {
		t.Fatalf("idle sources were not evicted: %v", err)
	}
	if len(g.create) != 1 {
		t.Fatalf("map size %d", len(g.create))
	}
}

func TestGuardConcurrentExactQuota(t *testing.T) {
	const limit = 80
	g := NewGuard(limit, limit)
	var okCount atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < limit*4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if g.CheckCreate("same-source") == nil {
				okCount.Add(1)
			}
			_ = g.CheckRedirect("same-source")
		}()
	}
	wg.Wait()
	if got := okCount.Load(); got != limit {
		t.Fatalf("concurrent allows %d, want %d", got, limit)
	}
}

func TestNewGuardDefaultsNonPositiveLimits(t *testing.T) {
	g := NewGuard(0, -1)
	if g.createLimit != 30 || g.redirectLimit != 120 {
		t.Fatalf("defaults %d %d", g.createLimit, g.redirectLimit)
	}
}

func assertRateLimited(t *testing.T, err *domain.AppError) {
	t.Helper()
	if err == nil || err.Type != domain.TypeRateLimited || err.Message != domain.MsgTooManyRequests {
		t.Fatalf("want rate_limited, got %v", err)
	}
}
