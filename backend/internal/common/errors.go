package common

import "errors"

var (
	ErrNotFound            = errors.New("resource not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrConflict            = errors.New("conflict")
	ErrValidation          = errors.New("validation failed")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidTransition   = errors.New("invalid state transition")
	ErrDuplicate           = errors.New("duplicate resource")
)

// AppError carries a user-facing message and optional field-level errors.
type AppError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
	Err     error             `json:"-"`
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

func NewAppError(code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func ValidationError(fields map[string]string) *AppError {
	return &AppError{
		Code:    "VALIDATION_FAILED",
		Message: "validation failed",
		Fields:  fields,
		Err:     ErrValidation,
	}
}
