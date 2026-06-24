package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	paymentdomain "github.com/younesbeheshti/any-task-connect/backend/internal/payment/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/wallet/service"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetWallet(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	wallet, err := h.svc.GetWallet(c.Request.Context(), userID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, wallet)
}

func (h *Handler) GetHistory(c *gin.Context) {
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
	result, err := h.svc.GetTransactions(c.Request.Context(), userID, paymentdomain.TransactionFilter{}, pg)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, result)
}

func (h *Handler) GetStatistics(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	stats, err := h.svc.GetStatistics(c.Request.Context(), userID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, stats)
}

type topUpRequest struct {
	Amount    int64   `json:"amount" binding:"required,min=1"`
	ReturnURL string  `json:"returnUrl"`
	CardID    *string `json:"cardId"`
}

func (h *Handler) TopUp(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	var req topUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	var cardID *uuid.UUID
	if req.CardID != nil {
		id, err := uuid.Parse(*req.CardID)
		if err == nil {
			cardID = &id
		}
	}
	url, err := h.svc.TopUp(c.Request.Context(), userID, req.Amount, cardID, req.ReturnURL)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"paymentUrl": url})
}

func (h *Handler) GetAdminOverview(c *gin.Context) {
	apiresponse.JSON(c, http.StatusOK, gin.H{"message": "financial overview"})
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
