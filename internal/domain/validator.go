package domain

import (
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// MaxLongURLLength is L from REQ-011 / OpenAPI maxLength.
const MaxLongURLLength = 2048

var shortCodePattern = regexp.MustCompile(`^[A-Za-z0-9]{8}$`)

// URLValidator implements UrlValidator from the method spec.
type URLValidator struct{}

func NewURLValidator() *URLValidator {
	return &URLValidator{}
}

// ValidateLongURL checks that raw is a non-empty absolute http(s) URL of at most 2048 characters.
func (v *URLValidator) ValidateLongURL(raw string) (ValidatedURL, *AppError) {
	if raw == "" || !utf8.ValidString(raw) || utf8.RuneCountInString(raw) > MaxLongURLLength {
		return ValidatedURL{}, InvalidURL(MsgURLValidationFailed)
	}
	if hasCTLOrSpace(raw) {
		return ValidatedURL{}, InvalidURL(MsgURLValidationFailed)
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed == nil || !parsed.IsAbs() || parsed.Opaque != "" {
		return ValidatedURL{}, InvalidURL(MsgURLValidationFailed)
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return ValidatedURL{}, InvalidURL(MsgURLValidationFailed)
	}
	if parsed.Hostname() == "" {
		return ValidatedURL{}, InvalidURL(MsgURLValidationFailed)
	}
	return ValidatedURL{URL: raw}, nil
}

// ValidateShortCode checks ^[A-Za-z0-9]{8}$ and returns the code unchanged.
func (v *URLValidator) ValidateShortCode(code string) (string, *AppError) {
	if !shortCodePattern.MatchString(code) {
		return "", InvalidURL(MsgInvalidShortCode)
	}
	return code, nil
}

func hasCTLOrSpace(s string) bool {
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b <= 0x20 || b == 0x7f {
			return true
		}
	}
	return false
}
