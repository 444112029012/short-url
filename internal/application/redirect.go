package application

import (
	"context"
	"errors"

	"github.com/444112029012/short-url/internal/domain"
	"github.com/444112029012/short-url/internal/ports"
)

// RedirectService implements RedirectService.redirect.
// The target is the stored long_url. Request parameters are not inputs.
type RedirectService struct {
	validator *domain.URLValidator
	repo      ports.URLRepository
	counter   ports.ClickCounter
}

func NewRedirectService(validator *domain.URLValidator, repo ports.URLRepository, counter ports.ClickCounter) *RedirectService {
	return &RedirectService{validator: validator, repo: repo, counter: counter}
}

// Redirect validates the code, loads the stored URL, then atomically increments.
// On any failure the count is left unchanged and no target is returned.
func (s *RedirectService) Redirect(ctx context.Context, shortCode string) (domain.RedirectTarget, *domain.AppError) {
	if _, appErr := s.validator.ValidateShortCode(shortCode); appErr != nil {
		return domain.RedirectTarget{}, appErr
	}
	rec, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return domain.RedirectTarget{}, domain.NotFound(domain.MsgShortCodeNotFound)
		}
		return domain.RedirectTarget{}, domain.InternalWrap(err)
	}
	if rec == nil || !safeLocation(rec.LongURL) {
		return domain.RedirectTarget{}, domain.Internal()
	}
	if _, err := s.counter.IncrementAtomic(ctx, shortCode); err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return domain.RedirectTarget{}, domain.NotFound(domain.MsgShortCodeNotFound)
		}
		return domain.RedirectTarget{}, domain.InternalWrap(err)
	}
	return domain.RedirectTarget{LongURL: rec.LongURL}, nil
}

func safeLocation(value string) bool {
	if value == "" {
		return false
	}
	for i := 0; i < len(value); i++ {
		b := value[i]
		if b < 0x20 || b == 0x7f {
			return false
		}
	}
	return true
}
