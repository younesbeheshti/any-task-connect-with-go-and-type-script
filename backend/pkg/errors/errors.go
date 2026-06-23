package errors

import (
	stderrors "errors"
	"net/http"
)

// Application sentinel errors.
var (
	ErrUnauthorized = stderrors.New("unauthorized")
	ErrForbidden    = stderrors.New("forbidden")
	ErrValidation   = stderrors.New("validation failed")
	ErrNotFound     = stderrors.New("resource not found")
	ErrConflict     = stderrors.New("conflict")
	ErrInternal     = stderrors.New("internal server error")
)

// AppError carries HTTP status and client-facing details.
type AppError struct {
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	HTTPStatus int               `json:"-"`
	Fields     map[string]string `json:"fields,omitempty"`
	Err        error             `json:"-"`
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "application error"
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates an AppError.
func New(code, message string, status int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
		Err:        err,
	}
}

// Validation creates a validation AppError with field details.
func Validation(fields map[string]string) *AppError {
	return &AppError{
		Code:       "VALIDATION_FAILED",
		Message:    "اطلاعات وارد شده معتبر نیست",
		HTTPStatus: http.StatusBadRequest,
		Fields:     fields,
		Err:        ErrValidation,
	}
}

// HTTPStatus maps errors to HTTP status codes.
func HTTPStatus(err error) int {
	var appErr *AppError
	if stderrors.As(err, &appErr) && appErr.HTTPStatus != 0 {
		return appErr.HTTPStatus
	}

	switch {
	case stderrors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case stderrors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case stderrors.Is(err, ErrValidation):
		return http.StatusBadRequest
	case stderrors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case stderrors.Is(err, ErrConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// Code returns a stable error code for clients.
func Code(err error) string {
	var appErr *AppError
	if stderrors.As(err, &appErr) && appErr.Code != "" {
		return appErr.Code
	}

	switch {
	case stderrors.Is(err, ErrUnauthorized):
		return "UNAUTHORIZED"
	case stderrors.Is(err, ErrForbidden):
		return "FORBIDDEN"
	case stderrors.Is(err, ErrValidation):
		return "VALIDATION_FAILED"
	case stderrors.Is(err, ErrNotFound):
		return "NOT_FOUND"
	case stderrors.Is(err, ErrConflict):
		return "CONFLICT"
	default:
		return "INTERNAL_ERROR"
	}
}

// Message returns a safe client-facing message.
func Message(err error) string {
	var appErr *AppError
	if stderrors.As(err, &appErr) && appErr.Message != "" {
		return appErr.Message
	}
	return "An unexpected error occurred"
}

// Fields returns validation field errors when present.
func Fields(err error) map[string]string {
	var appErr *AppError
	if stderrors.As(err, &appErr) {
		return appErr.Fields
	}
	return nil
}
