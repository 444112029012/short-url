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
	if g.create.get("203.0.113.5").count != 3 {
		t.Fatalf("over-limit consumed quota: %d", g.create.get("203.0.113.5").count)
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
	if g.redirect.get("203.0.113.5").count != 2 {
		t.Fatalf("UT-RL-02 over-limit consumed quota: %d", g.redirect.get("203.0.113.5").count)
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
	if g.create.get(emptySourceKey) == nil {
		t.Fatal("empty source was not stored on the shared key")
	}
}

func TestGuardFullTableEvictsOldestAndAdmits(t *testing.T) {
	// SCR-006: a full in-window table admits the new source by dropping the oldest.
	now := time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC)
	g := NewGuard(5, 5)
	g.maxSources = 2
	g.now = func() time.Time { return now }
	if err := g.CheckCreate("oldest"); err != nil {
		t.Fatal(err)
	}
	now = now.Add(time.Second)
	if err := g.CheckCreate("newer"); err != nil {
		t.Fatal(err)
	}
	now = now.Add(time.Second)
	if err := g.CheckCreate("admitted"); err != nil {
		t.Fatalf("full table rejected a new source: %v", err)
	}
	if g.create.order.Len() > 2 || len(g.create.byKey) > 2 {
		t.Fatalf("create table size list=%d map=%d", g.create.order.Len(), len(g.create.byKey))
	}
	if g.create.get("oldest") != nil {
		t.Fatal("oldest create source was kept")
	}
	if g.create.get("newer") == nil || g.create.get("admitted") == nil {
		t.Fatal("expected newer and admitted to remain")
	}
	if front := g.create.order.Front().Value.(*bucket).key; front != "newer" {
		t.Fatalf("front %s", front)
	}
}

func TestGuardFullRedirectTableEvictsOldestAndAdmits(t *testing.T) {
	now := time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC)
	g := NewGuard(5, 5)
	g.maxSources = 2
	g.now = func() time.Time { return now }
	if err := g.CheckRedirect("oldest"); err != nil {
		t.Fatal(err)
	}
	now = now.Add(time.Second)
	if err := g.CheckRedirect("newer"); err != nil {
		t.Fatal(err)
	}
	now = now.Add(time.Second)
	if err := g.CheckRedirect("admitted"); err != nil {
		t.Fatalf("full redirect table rejected a new source: %v", err)
	}
	if g.redirect.order.Len() > 2 || len(g.redirect.byKey) > 2 {
		t.Fatalf("redirect table size list=%d map=%d", g.redirect.order.Len(), len(g.redirect.byKey))
	}
	if g.redirect.get("oldest") != nil {
		t.Fatal("oldest redirect source was kept")
	}
	if g.redirect.get("newer") == nil || g.redirect.get("admitted") == nil {
		t.Fatal("expected newer and admitted to remain")
	}
}

func TestGuardWindowResetMovesBucketToBack(t *testing.T) {
	now := time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC)
	g := NewGuard(1, 1)
	g.maxSources = 2
	g.now = func() time.Time { return now }
	if err := g.CheckCreate("reset-me"); err != nil {
		t.Fatal(err)
	}
	now = now.Add(time.Second)
	if err := g.CheckCreate("still-open"); err != nil {
		t.Fatal(err)
	}
	// reset-me's window elapses; still-open (started 1s later) remains inside it.
	now = now.Add(time.Minute - time.Second)
	if err := g.CheckCreate("reset-me"); err != nil {
		t.Fatalf("window reset rejected: %v", err)
	}
	if front := g.create.order.Front().Value.(*bucket).key; front != "still-open" {
		t.Fatalf("reset bucket stayed at front: %s", front)
	}
	if back := g.create.order.Back().Value.(*bucket).key; back != "reset-me" {
		t.Fatalf("reset bucket not at back: %s", back)
	}
	if err := g.CheckCreate("newcomer"); err != nil {
		t.Fatal(err)
	}
	if g.create.get("still-open") != nil {
		t.Fatal("in-window older bucket was not the one evicted")
	}
	if g.create.get("reset-me") == nil || g.create.get("newcomer") == nil {
		t.Fatal("reset bucket was evicted ahead of the older window")
	}
}

func TestGuardRedirectWindowResetMovesBucketToBack(t *testing.T) {
	now := time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC)
	g := NewGuard(1, 1)
	g.maxSources = 2
	g.now = func() time.Time { return now }
	if err := g.CheckRedirect("reset-me"); err != nil {
		t.Fatal(err)
	}
	now = now.Add(time.Second)
	if err := g.CheckRedirect("still-open"); err != nil {
		t.Fatal(err)
	}
	now = now.Add(time.Minute - time.Second)
	if err := g.CheckRedirect("reset-me"); err != nil {
		t.Fatalf("window reset rejected: %v", err)
	}
	if front := g.redirect.order.Front().Value.(*bucket).key; front != "still-open" {
		t.Fatalf("reset bucket stayed at front: %s", front)
	}
	if err := g.CheckRedirect("newcomer"); err != nil {
		t.Fatal(err)
	}
	if g.redirect.get("still-open") != nil || g.redirect.get("reset-me") == nil {
		t.Fatal("redirect table evicted the reset bucket instead of the older window")
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
