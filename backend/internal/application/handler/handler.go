package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/application/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/application/service"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

// Handler serves application HTTP endpoints.
type Handler struct {
	svc *service.ApplicationService
}

// NewHandler creates an application Handler.
func NewHandler(svc *service.ApplicationService) *Handler {
	return &Handler{svc: svc}
}

type submitApplicationRequest struct {
	ProposalMessage string `json:"message" binding:"required"`
	ProposedPrice   *int64 `json:"price"`
	ETA             string `json:"eta"`
}

// Submit handles POST /tasks/:id/applications — agent submits a proposal.
func (h *Handler) Submit(c *gin.Context) {
	agentID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_FIELD", "شناسه تسک نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	var req submitApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}

	input := domain.SubmitApplicationInput{
		TaskID:          taskID,
		AgentID:         agentID,
		ProposalMessage: req.ProposalMessage,
		ProposedPrice:   req.ProposedPrice,
		ETA:             req.ETA,
	}
	app, err := h.svc.Submit(c.Request.Context(), input)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusCreated, app)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.JSON(c, http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_ID", "message": "شناسه نامعتبر است"}})
		return
	}
	app, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, app)
}

func (h *Handler) ListByTask(c *gin.Context) {
	requesterID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.JSON(c, http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_ID", "message": "شناسه تسک نامعتبر است"}})
		return
	}
	apps, err := h.svc.ListByTask(c.Request.Context(), taskID, requesterID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"applications": apps})
}

func (h *Handler) ListMyApplications(c *gin.Context) {
	agentID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apps, err := h.svc.ListByAgent(c.Request.Context(), agentID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"applications": apps})
}

func (h *Handler) Accept(c *gin.Context) {
	requesterID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.JSON(c, http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_ID", "message": "شناسه نامعتبر است"}})
		return
	}
	app, err := h.svc.Accept(c.Request.Context(), id, requesterID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, app)
}

func (h *Handler) Reject(c *gin.Context) {
	requesterID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.JSON(c, http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_ID", "message": "شناسه نامعتبر است"}})
		return
	}
	app, err := h.svc.Reject(c.Request.Context(), id, requesterID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, app)
}

func (h *Handler) Withdraw(c *gin.Context) {
	agentID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.JSON(c, http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_ID", "message": "شناسه نامعتبر است"}})
		return
	}
	app, err := h.svc.Withdraw(c.Request.Context(), id, agentID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, app)
}

func currentUserID(c *gin.Context) (uuid.UUID, error) {
	idStr := c.GetString(middleware.UserIDKey)
	if idStr == "" {
		return uuid.Nil, apperrors.New("UNAUTHORIZED", "لطفاً وارد حساب کاربری شوید", 401, apperrors.ErrUnauthorized)
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, apperrors.New("UNAUTHORIZED", "شناسه کاربر نامعتبر است", 401, apperrors.ErrUnauthorized)
	}
	return id, nil
}
