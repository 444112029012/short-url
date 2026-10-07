package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/444112029012/short-url/internal/domain"
	"github.com/444112029012/short-url/internal/ports"
)

const defaultGenerateAttempts = 8

// CodeGenerator is ShortCodeGenerator.generate.
type CodeGenerator interface {
	Generate() (string, *domain.AppError)
}

// ShortURLApplicationService implements createShortUrl and getStats.
type ShortURLApplicationService struct {
	validator   *domain.URLValidator
	generator   CodeGenerator
	repo        ports.URLRepository
	counter     ports.ClickCounter
	baseURL     string
	maxAttempts int
	now         func() time.Time
}

func NewShortURLApplicationService(
	validator *domain.URLValidator,
	generator CodeGenerator,
	repo ports.URLRepository,
	counter ports.ClickCounter,
	baseURL string,
) *ShortURLApplicationService {
	return &ShortURLApplicationService{
		validator:   validator,
		generator:   generator,
		repo:        repo,
		counter:     counter,
		baseURL:     strings.TrimRight(baseURL, "/"),
		maxAttempts: defaultGenerateAttempts,
		now:         time.Now,
	}
}

// CreateShortURL validates the long URL, generates a unique short code, and stores click_count 0.
func (s *ShortURLApplicationService) CreateShortURL(ctx context.Context, longURL string) (domain.URLCreated, *domain.AppError) {
	validated, appErr := s.validator.ValidateLongURL(longURL)
	if appErr != nil {
		return domain.URLCreated{}, appErr
	}
	attempts := s.maxAttempts
	if attempts <= 0 {
		attempts = defaultGenerateAttempts
	}
	for i := 0; i < attempts; i++ {
		code, appErr := s.generator.Generate()
		if appErr != nil {
			return domain.URLCreated{}, appErr
		}
		if _, appErr = s.validator.ValidateShortCode(code); appErr != nil {
			return domain.URLCreated{}, domain.Internal()
		}
		exists, err := s.repo.Exists(ctx, code)
		if err != nil {
			return domain.URLCreated{}, domain.InternalWrap(err)
		}
		if exists {
			continue
		}
		now := s.now
		if now == nil {
			now = time.Now
		}
		record := domain.URLRecord{
			ShortCode:  code,
			LongURL:    validated.URL,
			ClickCount: 0,
			CreatedAt:  now().UTC(),
		}
		if err := s.repo.Save(ctx, record); err != nil {
			if errors.Is(err, ports.ErrConflict) {
				continue
			}
			return domain.URLCreated{}, domain.InternalWrap(err)
		}
		return domain.URLCreated{
			ShortCode: code,
			ShortURL:  s.baseURL + "/" + code,
			LongURL:   validated.URL,
		}, nil
	}
	return domain.URLCreated{}, domain.Internal()
}

// GetStats validates the short code and reads click_count. It does not write.
func (s *ShortURLApplicationService) GetStats(ctx context.Context, shortCode string) (domain.URLStats, *domain.AppError) {
	code, appErr := s.validator.ValidateShortCode(shortCode)
	if appErr != nil {
		return domain.URLStats{}, appErr
	}
	n, err := s.counter.GetCount(ctx, code)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return domain.URLStats{}, domain.NotFound(domain.MsgShortCodeNotFound)
		}
		return domain.URLStats{}, domain.InternalWrap(err)
	}
	if n < 0 {
		return domain.URLStats{}, domain.Internal()
	}
	return domain.URLStats{ShortCode: code, ClickCount: n}, nil
}
