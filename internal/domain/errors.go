package domain

// Stable public error types. HTTP bodies may contain only these type strings
// and the fixed messages in this file (REQ-009, SEC-003).

type ErrorType string

const (
	TypeInvalidURL  ErrorType = "invalid_url"
	TypeNotFound    ErrorType = "not_found"
	TypeRateLimited ErrorType = "rate_limited"
	TypeInternal    ErrorType = "internal_error"
)

const (
	MsgURLValidationFailed = "URL validation failed"
	MsgInvalidShortCode    = "Invalid short code format"
	MsgShortCodeNotFound   = "Short code not found"
	MsgTooManyRequests     = "Too many requests"
	MsgUnexpected          = "An unexpected error occurred"
)

// AppError is the application error carrier from the method spec.
// Error() returns only the stable type. cause is never serialized.
type AppError struct {
	Type    ErrorType
	Message string
	cause   error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	return string(e.Type)
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

func InvalidURL(message string) *AppError {
	return &AppError{Type: TypeInvalidURL, Message: message}
}

func NotFound(message string) *AppError {
	return &AppError{Type: TypeNotFound, Message: message}
}

func RateLimited(message string) *AppError {
	return &AppError{Type: TypeRateLimited, Message: message}
}

func Internal() *AppError {
	return &AppError{Type: TypeInternal, Message: MsgUnexpected}
}

func InternalWrap(cause error) *AppError {
	return &AppError{Type: TypeInternal, Message: MsgUnexpected, cause: cause}
}
