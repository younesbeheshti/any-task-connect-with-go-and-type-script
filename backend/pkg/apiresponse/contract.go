package apiresponse

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	govalidator "github.com/go-playground/validator/v10"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/validator"
)

// ErrorBody matches front/docs/api-contracts.md error envelope.
type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail is the nested error object.
type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// WriteError writes a frontend-compatible error response.
func WriteError(c *gin.Context, err error) {
	status, detail := mapError(err)
	c.JSON(status, ErrorBody{Error: detail})
}

func mapError(err error) (int, ErrorDetail) {
	var valErr validator.ValidationError
	if errors.As(err, &valErr) {
		return http.StatusBadRequest, ErrorDetail{
			Code:    "VALIDATION_FAILED",
			Message: "اطلاعات وارد شده معتبر نیست",
			Fields:  valErr.Fields,
		}
	}

	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return appErr.HTTPStatus, ErrorDetail{
			Code:    appErr.Code,
			Message: appErr.Message,
			Fields:  appErr.Fields,
		}
	}

	// gin's c.ShouldBindJSON validates via go-playground tags and returns raw
	// ValidationErrors; surface them as a 400 with per-field info instead of a 500.
	var bindErrs govalidator.ValidationErrors
	if errors.As(err, &bindErrs) {
		fields := make(map[string]string, len(bindErrs))
		for _, fe := range bindErrs {
			fields[strings.ToLower(fe.Field())] = "فیلد الزامی یا نامعتبر است"
		}
		return http.StatusBadRequest, ErrorDetail{
			Code:    "VALIDATION_FAILED",
			Message: "اطلاعات وارد شده معتبر نیست",
			Fields:  fields,
		}
	}

	// Malformed or empty JSON request body → 400, not 500.
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return http.StatusBadRequest, ErrorDetail{
			Code:    "INVALID_BODY",
			Message: "بدنه درخواست نامعتبر است",
		}
	}

	switch {
	case errors.Is(err, apperrors.ErrUnauthorized):
		return http.StatusUnauthorized, ErrorDetail{Code: "UNAUTHORIZED", Message: "لطفاً وارد حساب کاربری شوید"}
	case errors.Is(err, apperrors.ErrForbidden):
		return http.StatusForbidden, ErrorDetail{Code: "FORBIDDEN", Message: "دسترسی مجاز نیست"}
	case errors.Is(err, apperrors.ErrNotFound):
		return http.StatusNotFound, ErrorDetail{Code: "NOT_FOUND", Message: "مورد درخواستی یافت نشد"}
	case errors.Is(err, apperrors.ErrConflict):
		return http.StatusConflict, ErrorDetail{Code: "CONFLICT", Message: "این اطلاعات قبلاً ثبت شده است"}
	case errors.Is(err, apperrors.ErrValidation):
		return http.StatusBadRequest, ErrorDetail{Code: "VALIDATION_FAILED", Message: "اطلاعات وارد شده معتبر نیست"}
	default:
		return http.StatusInternalServerError, ErrorDetail{Code: "INTERNAL_ERROR", Message: "خطای داخلی سرور"}
	}
}

// JSON writes a raw JSON response (for contract-shaped success bodies).
func JSON(c *gin.Context, status int, body any) {
	c.JSON(status, body)
}

// NoContent writes 204.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
