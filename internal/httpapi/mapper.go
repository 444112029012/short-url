package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/444112029012/short-url/internal/domain"
)

// ErrorMapper implements ErrorMapper.toHttpResponse.
// The body is only {error:{type,message}} using a fixed message catalog.
type ErrorMapper struct{}

func NewErrorMapper() ErrorMapper { return ErrorMapper{} }

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// ToHTTPResponse maps an AppError onto the four stable HTTP results.
// Unknown types become internal_error. Caller messages are not copied unless
// they are one of the fixed catalog strings for that type.
func (ErrorMapper) ToHTTPResponse(err *domain.AppError) response {
	if err == nil {
		err = domain.Internal()
	}
	status, typ := statusFor(err.Type)
	payload := errorEnvelope{Error: errorBody{
		Type:    string(typ),
		Message: safeMessage(typ, err.Message),
	}}
	raw, encErr := json.Marshal(payload)
	if encErr != nil {
		return response{
			status: http.StatusInternalServerError,
			header: jsonHeader(),
			body:   []byte(`{"error":{"type":"internal_error","message":"An unexpected error occurred"}}`),
		}
	}
	return response{status: status, header: jsonHeader(), body: raw}
}

func statusFor(typ domain.ErrorType) (int, domain.ErrorType) {
	switch typ {
	case domain.TypeInvalidURL:
		return http.StatusBadRequest, domain.TypeInvalidURL
	case domain.TypeNotFound:
		return http.StatusNotFound, domain.TypeNotFound
	case domain.TypeRateLimited:
		return http.StatusTooManyRequests, domain.TypeRateLimited
	case domain.TypeInternal:
		return http.StatusInternalServerError, domain.TypeInternal
	default:
		return http.StatusInternalServerError, domain.TypeInternal
	}
}

func safeMessage(typ domain.ErrorType, message string) string {
	switch typ {
	case domain.TypeInvalidURL:
		if message == domain.MsgInvalidShortCode || message == domain.MsgURLValidationFailed {
			return message
		}
		return domain.MsgURLValidationFailed
	case domain.TypeNotFound:
		return domain.MsgShortCodeNotFound
	case domain.TypeRateLimited:
		return domain.MsgTooManyRequests
	default:
		return domain.MsgUnexpected
	}
}
