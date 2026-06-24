package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	"github.com/younesbeheshti/any-task-connect/backend/internal/notification/service"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// GET /notifications
func (h *Handler) List(c *gin.Context) {
	userID := mustUserID(c)
	unreadOnly := c.Query("unread") == "true"
	var pg common.PaginationParams
	_ = c.ShouldBindQuery(&pg)
	result, err := h.svc.List(c.Request.Context(), userID, unreadOnly, pg)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, result)
}

// PATCH /notifications/:id/read
func (h *Handler) MarkRead(c *gin.Context) {
	userID := mustUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	if err := h.svc.MarkRead(c.Request.Context(), id, userID); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"ok": true})
}

// POST /notifications/read-all
func (h *Handler) MarkAllRead(c *gin.Context) {
	userID := mustUserID(c)
	if err := h.svc.MarkAllRead(c.Request.Context(), userID); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"ok": true})
}

func mustUserID(c *gin.Context) uuid.UUID {
	id, _ := uuid.Parse(c.GetString(middleware.UserIDKey))
	return id
}
