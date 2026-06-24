package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	"github.com/younesbeheshti/any-task-connect/backend/internal/withdraw/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/withdraw/service"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// createWithdrawRequest matches POST /wallet/withdraw per frontend contract.
type createWithdrawRequest struct {
	Amount      int64  `json:"amount" binding:"required,min=1"`
	IBAN        string `json:"iban"`        // frontend contract field
	BankAccount string `json:"bankAccount"` // alternate field
	Description string `json:"description"`
}

func (h *Handler) Create(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	var req createWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	sheba := req.IBAN
	if sheba == "" {
		sheba = req.BankAccount
	}
	input := domain.CreateInput{
		UserID:      userID,
		Amount:      req.Amount,
		ShebaNumber: sheba,
		Description: req.Description,
	}
	w, err := h.svc.CreateWithdraw(c.Request.Context(), input)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusCreated, toResponse(w))
}

func (h *Handler) ListMine(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	var pg common.PaginationParams
	if err := c.ShouldBindQuery(&pg); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	result, err := h.svc.ListUserWithdraws(c.Request.Context(), userID, pg)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, result)
}

func (h *Handler) GetByID(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	w, err := h.svc.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, toResponse(w))
}

func (h *Handler) Approve(c *gin.Context) {
	adminID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	if err := h.svc.Approve(c.Request.Context(), id, adminID); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"message": "درخواست تایید شد"})
}

func (h *Handler) Reject(c *gin.Context) {
	adminID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	if err := h.svc.Reject(c.Request.Context(), id, adminID); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"message": "درخواست رد شد"})
}

func (h *Handler) ListAll(c *gin.Context) {
	var pg common.PaginationParams
	if err := c.ShouldBindQuery(&pg); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	result, err := h.svc.ListAll(c.Request.Context(), pg)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, result)
}

type withdrawResponse struct {
	ID          string  `json:"id"`
	Amount      int64   `json:"amount"`
	Status      string  `json:"status"`
	BankAccount string  `json:"bankAccount"`
	ShebaNumber string  `json:"shebaNumber"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"createdAt"`
}

func toResponse(w *domain.WithdrawRequest) withdrawResponse {
	return withdrawResponse{
		ID:          w.ID.String(),
		Amount:      w.Amount,
		Status:      w.Status.APIStatus(),
		BankAccount: w.BankAccount,
		ShebaNumber: w.ShebaNumber,
		Description: w.Description,
		CreatedAt:   w.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
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
