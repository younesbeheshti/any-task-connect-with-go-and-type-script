package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/chat/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/chat/service"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
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
	taskRef := c.Param("id")
	var before *uuid.UUID
	if raw := c.Query("before"); raw != "" {
		id, err := uuid.Parse(raw)
		if err == nil {
			before = &id
		}
	}
	limit := 50
	msgs, err := h.svc.ListMessages(c.Request.Context(), taskRef, userID, before, limit)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"messages": msgs})
}

// POST /tasks/:id/messages
func (h *Handler) SendMessage(c *gin.Context) {
	userID := mustUserID(c)
	taskRef := c.Param("id")
	var req struct {
		Message    string             `json:"message"`
		Attachment *domain.Attachment `json:"attachment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	// Receiver is derived from the task (the other participant); the client only
	// supplies the message/attachment.
	msg, err := h.svc.SendMessage(c.Request.Context(), domain.SendMessageInput{
		TaskRef:    taskRef,
		SenderID:   userID,
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
	taskRef := c.Param("id")
	if err := h.svc.MarkRead(c.Request.Context(), taskRef, userID); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"ok": true})
}

func mustUserID(c *gin.Context) uuid.UUID {
	id, _ := uuid.Parse(c.GetString(middleware.UserIDKey))
	return id
}
