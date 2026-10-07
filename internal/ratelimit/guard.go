package ratelimit

import (
	"sync"
	"time"

	"github.com/444112029012/short-url/internal/domain"
)

const (
	window = time.Minute
	// defaultMaxSources caps each counter map. Idle windows are evicted before
	// a new source is inserted; if every slot is still inside its window, a
	// new source is rejected (fail closed) so the map cannot grow without bound.
	defaultMaxSources = 8192
	// emptySourceKey is the shared bucket for a missing source identity.
	// An empty sourceID must not bypass the limit: every empty call shares
	// one quota instead of receiving its own unlimited identity.
	emptySourceKey = "empty-source"
)

// Guard is the in-process RateLimitGuard (SEC-007, NFR-003).
// Counters are per source and per method, fixed one-minute windows, and are
// not shared across processes. Pass consumes one quota; a rejection does not.
// This type never logs the source id (SEC-005).
type Guard struct {
	createLimit   int
	redirectLimit int
	maxSources    int
	now           func() time.Time

	mu       sync.Mutex
	create   map[string]*bucket
	redirect map[string]*bucket
}

type bucket struct {
	windowStart time.Time
	count       int
}

// NewGuard builds a guard. Non-positive limits fall back to the documented
// defaults (30 creates/minute, 120 redirects/minute) so a zero value cannot
// disable the control. Limits come from RATE_LIMIT_CREATE_PER_MIN and
// RATE_LIMIT_REDIRECT_PER_MIN.
func NewGuard(createPerMin, redirectPerMin int) *Guard {
	if createPerMin <= 0 {
		createPerMin = 30
	}
	if redirectPerMin <= 0 {
		redirectPerMin = 120
	}
	return &Guard{
		createLimit:   createPerMin,
		redirectLimit: redirectPerMin,
		maxSources:    defaultMaxSources,
		now:           time.Now,
		create:        make(map[string]*bucket),
		redirect:      make(map[string]*bucket),
	}
}

// CheckCreate consumes one create quota for sourceID or returns rate_limited.
func (g *Guard) CheckCreate(sourceID string) *domain.AppError {
	return g.allow(g.create, g.createLimit, sourceID)
}

// CheckRedirect consumes one redirect quota for sourceID or returns rate_limited.
// The caller must invoke this before redirecting or incrementing click_count.
func (g *Guard) CheckRedirect(sourceID string) *domain.AppError {
	return g.allow(g.redirect, g.redirectLimit, sourceID)
}

func (g *Guard) allow(buckets map[string]*bucket, limit int, sourceID string) *domain.AppError {
	if g == nil {
		return domain.RateLimited(domain.MsgTooManyRequests)
	}
	if sourceID == "" {
		sourceID = emptySourceKey
	}
	now := time.Now()
	if g.now != nil {
		now = g.now()
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	b := buckets[sourceID]
	if b != nil && now.Sub(b.windowStart) >= window {
		delete(buckets, sourceID)
		b = nil
	}
	if b == nil {
		if len(buckets) >= g.maxSources {
			g.evictIdle(buckets, now)
		}
		if len(buckets) >= g.maxSources {
			return domain.RateLimited(domain.MsgTooManyRequests)
		}
		buckets[sourceID] = &bucket{windowStart: now, count: 1}
		return nil
	}
	if b.count >= limit {
		return domain.RateLimited(domain.MsgTooManyRequests)
	}
	b.count++
	return nil
}

func (g *Guard) evictIdle(buckets map[string]*bucket, now time.Time) {
	for key, b := range buckets {
		if b == nil || now.Sub(b.windowStart) >= window {
			delete(buckets, key)
		}
	}
}
