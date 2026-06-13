package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

// Envelope is the standard API response shape.
type Envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    any         `json:"data,omitempty"`
	Errors  any         `json:"errors,omitempty"`
	Meta    *Pagination `json:"meta,omitempty"`
}

// Pagination holds list metadata.
type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

// SuccessResponse writes a success envelope.
func SuccessResponse(c *gin.Context, status int, message string, data any) {
	c.JSON(status, Envelope{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse writes an error envelope.
func ErrorResponse(c *gin.Context, err error) {
	status := apperrors.HTTPStatus(err)
	c.JSON(status, Envelope{
		Success: false,
		Message: apperrors.Message(err),
		Errors:  errorDetails(err),
	})
}

func errorDetails(err error) any {
	fields := apperrors.Fields(err)
	if len(fields) > 0 {
		return fields
	}
	return gin.H{"code": apperrors.Code(err)}
}

// PaginatedResponse writes a paginated list response.
func PaginatedResponse(c *gin.Context, message string, items any, page, pageSize int, total int64) {
	c.JSON(http.StatusOK, Envelope{
		Success: true,
		Message: message,
		Data:    items,
		Meta: &Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	})
}

// OK is a shorthand for 200 success responses.
func OK(c *gin.Context, message string, data any) {
	SuccessResponse(c, http.StatusOK, message, data)
}

// Created is a shorthand for 201 success responses.
func Created(c *gin.Context, message string, data any) {
	SuccessResponse(c, http.StatusCreated, message, data)
}

// NoContent writes a 204 response.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
