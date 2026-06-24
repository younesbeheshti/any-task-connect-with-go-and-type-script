package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/chat/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/chat/service"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// GET /tasks/:id/messages
func (h *Handler) ListMessages(c *gin.Context) {
	userID := mustUserID(c)
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	var before *uuid.UUID
	if raw := c.Query("before"); raw != "" {
		id, err := uuid.Parse(raw)
		if err == nil {
			before = &id
		}
	}
	limit := 50
	msgs, err := h.svc.ListMessages(c.Request.Context(), taskID, userID, before, limit)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"messages": msgs})
}

// POST /tasks/:id/messages
func (h *Handler) SendMessage(c *gin.Context) {
	userID := mustUserID(c)
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	var req struct {
		ReceiverID string             `json:"receiverId" binding:"required"`
		Message    string             `json:"message"`
		Attachment *domain.Attachment `json:"attachment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	receiverID, err := uuid.Parse(req.ReceiverID)
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه گیرنده نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	msg, err := h.svc.SendMessage(c.Request.Context(), domain.SendMessageInput{
		TaskID:     taskID,
		SenderID:   userID,
		ReceiverID: receiverID,
		Message:    req.Message,
		Attachment: req.Attachment,
	})
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusCreated, msg)
}

// GET /chats
func (h *Handler) ListChats(c *gin.Context) {
	userID := mustUserID(c)
	chats, err := h.svc.ListChats(c.Request.Context(), userID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"chats": chats})
}

// POST /tasks/:id/messages/read
func (h *Handler) MarkRead(c *gin.Context) {
	userID := mustUserID(c)
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	if err := h.svc.MarkRead(c.Request.Context(), taskID, userID); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"ok": true})
}

func mustUserID(c *gin.Context) uuid.UUID {
	id, _ := uuid.Parse(c.GetString(middleware.UserIDKey))
	return id
}
