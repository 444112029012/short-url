package ratelimit

import "github.com/444112029012/short-url/internal/domain"

// Guard is an always-allow stand-in for RateLimitGuard.
//
// TODO(ENG-011): enforce per-source limits from configuration
// (RATE_LIMIT_CREATE_PER_MIN default 30, RATE_LIMIT_REDIRECT_PER_MIN default 120).
// When limited, return domain.RateLimited(domain.MsgTooManyRequests) and do not
// create a code, redirect, or increment click_count. This stub does not implement SEC-007.
type Guard struct{}

func NewGuard() *Guard { return &Guard{} }

// CheckCreate always allows. sourceID is accepted so the HTTP adapter can keep the call shape.
func (g *Guard) CheckCreate(sourceID string) *domain.AppError {
	_ = sourceID
	return nil
}

// CheckRedirect always allows.
func (g *Guard) CheckRedirect(sourceID string) *domain.AppError {
	_ = sourceID
	return nil
}
