package ports

import (
	"context"
	"errors"

	"github.com/444112029012/short-url/internal/domain"
)

var (
	// ErrConflict is a unique short_code collision. Callers may retry generation.
	ErrConflict = errors.New("short code conflict")
	// ErrNotFound means the short code is not stored.
	ErrNotFound = errors.New("short code not found")
)

// URLRepository is UrlRepository from the method spec.
// Implementations must use parameterized SQL (SEC-016).
type URLRepository interface {
	Save(ctx context.Context, record domain.URLRecord) error
	FindByShortCode(ctx context.Context, code string) (*domain.URLRecord, error)
	Exists(ctx context.Context, code string) (bool, error)
}

// ClickCounter is ClickCounter from the method spec.
// IncrementAtomic must be a single atomic UPDATE (SEC-016).
type ClickCounter interface {
	IncrementAtomic(ctx context.Context, shortCode string) (int, error)
	GetCount(ctx context.Context, shortCode string) (int, error)
}
