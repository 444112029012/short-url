package ratelimit

import (
	"container/list"
	"sync"
	"time"

	"github.com/444112029012/short-url/internal/domain"
)

const (
	window = time.Minute
	// defaultMaxSources caps each counter table. When the table is full, the
	// oldest window is evicted and the new source is admitted. A full table
	// never by itself returns rate_limited (SCR-006).
	defaultMaxSources = 8192
	// emptySourceKey is the shared bucket for a missing source identity.
	// An empty sourceID must not bypass the limit: every empty call shares
	// one quota instead of receiving its own unlimited identity.
	emptySourceKey = "empty-source"
)

// Guard is the in-process RateLimitGuard (SEC-007, NFR-003).
// Counters are per source and per method, fixed one-minute windows, and are
// not shared across processes. Pass consumes one quota; a rejection does not.
// This type never logs source identities, including ones it evicts (SEC-005).
type Guard struct {
	createLimit   int
	redirectLimit int
	maxSources    int
	now           func() time.Time

	mu       sync.Mutex
	create   table
	redirect table
}

// bucket is one fixed window. The list element is ordered by windowStart:
// front is oldest, back is newest.
type bucket struct {
	key         string
	windowStart time.Time
	count       int
}

// table is a map plus a FIFO list so admitting a source evicts in O(1).
type table struct {
	byKey map[string]*list.Element
	order *list.List
}

func newTable() table {
	return table{
		byKey: make(map[string]*list.Element),
		order: list.New(),
	}
}

func (t *table) get(key string) *bucket {
	if t == nil || t.byKey == nil {
		return nil
	}
	el := t.byKey[key]
	if el == nil {
		return nil
	}
	b, _ := el.Value.(*bucket)
	return b
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
		create:        newTable(),
		redirect:      newTable(),
	}
}

// CheckCreate consumes one create quota for sourceID or returns rate_limited.
func (g *Guard) CheckCreate(sourceID string) *domain.AppError {
	if g == nil {
		return domain.RateLimited(domain.MsgTooManyRequests)
	}
	return g.allow(&g.create, g.createLimit, sourceID)
}

// CheckRedirect consumes one redirect quota for sourceID or returns rate_limited.
// The caller must invoke this before redirecting or incrementing click_count.
func (g *Guard) CheckRedirect(sourceID string) *domain.AppError {
	if g == nil {
		return domain.RateLimited(domain.MsgTooManyRequests)
	}
	return g.allow(&g.redirect, g.redirectLimit, sourceID)
}

func (g *Guard) allow(tab *table, limit int, sourceID string) *domain.AppError {
	if g == nil || tab == nil || tab.byKey == nil || tab.order == nil {
		return domain.RateLimited(domain.MsgTooManyRequests)
	}
	if sourceID == "" {
		sourceID = emptySourceKey
	}
	now := time.Now()
	if g.now != nil {
		now = g.now()
	}
	capSources := g.maxSources
	if capSources < 1 {
		capSources = defaultMaxSources
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	if el := tab.byKey[sourceID]; el != nil {
		b, _ := el.Value.(*bucket)
		if b == nil {
			tab.order.Remove(el)
			delete(tab.byKey, sourceID)
		} else if now.Sub(b.windowStart) >= window {
			b.windowStart = now
			b.count = 1
			tab.order.MoveToBack(el)
			return nil
		} else if b.count >= limit {
			return domain.RateLimited(domain.MsgTooManyRequests)
		} else {
			b.count++
			return nil
		}
	}

	g.dropExpiredFront(tab, now)
	if tab.order.Len() >= capSources {
		g.evictFront(tab)
	}
	b := &bucket{key: sourceID, windowStart: now, count: 1}
	tab.byKey[sourceID] = tab.order.PushBack(b)
	return nil
}

// dropExpiredFront removes leading windows that have already elapsed.
// Only the front is inspected, so this stays cheap while the list is ordered
// by windowStart.
func (g *Guard) dropExpiredFront(tab *table, now time.Time) {
	for tab.order.Len() > 0 {
		el := tab.order.Front()
		b, _ := el.Value.(*bucket)
		if b != nil && now.Sub(b.windowStart) < window {
			return
		}
		g.evictFront(tab)
	}
}

// evictFront drops the oldest window. The evicted identity is not logged.
func (g *Guard) evictFront(tab *table) {
	el := tab.order.Front()
	if el == nil {
		return
	}
	tab.order.Remove(el)
	if b, _ := el.Value.(*bucket); b != nil {
		delete(tab.byKey, b.key)
	}
}
