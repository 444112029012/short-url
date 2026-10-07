package domain

import "time"

// URLRecord is the UrlMapping row persisted in DS1.
// Visitor IP and User-Agent are intentionally absent (SEC-004).
type URLRecord struct {
	ShortCode  string
	LongURL    string
	ClickCount int
	CreatedAt  time.Time
}

// ValidatedURL is a long URL that passed UrlValidator.validateLongUrl.
type ValidatedURL struct {
	URL string
}

// URLCreated is the createShortUrl success value.
type URLCreated struct {
	ShortCode string
	ShortURL  string
	LongURL   string
}

// URLStats is the getStats success value.
type URLStats struct {
	ShortCode  string
	ClickCount int
}

// RedirectTarget is the redirect success value. LongURL is the stored target.
type RedirectTarget struct {
	LongURL string
}
