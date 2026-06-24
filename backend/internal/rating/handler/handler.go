package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	"github.com/younesbeheshti/any-task-connect/backend/internal/rating/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/rating/service"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// POST /tasks/:id/reviews
func (h *Handler) Create(c *gin.Context) {
	reviewerID := mustUserID(c)
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	var req struct {
		ReviewedUserID string `json:"reviewedUserId" binding:"required"`
		Rating         int    `json:"rating"         binding:"required,min=1,max=5"`
		Comment        string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	reviewedID, err := uuid.Parse(req.ReviewedUserID)
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه کاربر ارزیابی‌شونده نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	review, err := h.svc.Create(c.Request.Context(), domain.CreateReviewInput{
		TaskID:         taskID,
		ReviewerID:     reviewerID,
		ReviewedUserID: reviewedID,
		Rating:         req.Rating,
		Comment:        req.Comment,
	})
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusCreated, review)
}

// GET /tasks/:id/reviews
func (h *Handler) ListByTask(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	reviews, err := h.svc.GetByTask(c.Request.Context(), taskID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"reviews": reviews})
}

// GET /users/:id/reviews
func (h *Handler) ListByUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	reviews, err := h.svc.GetByUser(c.Request.Context(), userID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"reviews": reviews})
}

func mustUserID(c *gin.Context) uuid.UUID {
	id, _ := uuid.Parse(c.GetString(middleware.UserIDKey))
	return id
}
