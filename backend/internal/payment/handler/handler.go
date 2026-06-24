package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	"github.com/younesbeheshti/any-task-connect/backend/internal/payment/domain"
	paymentrepo "github.com/younesbeheshti/any-task-connect/backend/internal/payment/repository"
	walletrepo "github.com/younesbeheshti/any-task-connect/backend/internal/wallet/repository"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

type Handler struct {
	txRepo     paymentrepo.Repository
	walletRepo walletrepo.Repository
}

func NewHandler(txRepo paymentrepo.Repository, walletRepo walletrepo.Repository) *Handler {
	return &Handler{txRepo: txRepo, walletRepo: walletRepo}
}

func (h *Handler) ListTransactions(c *gin.Context) {
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
	wallet, err := h.walletRepo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		apiresponse.JSON(c, http.StatusOK, gin.H{"items": []any{}, "page": 1, "pageSize": 20, "total": 0})
		return
	}
	filter := domain.TransactionFilter{WalletID: &wallet.ID}
	txs, total, err := h.txRepo.List(c.Request.Context(), filter, pg)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	responses := make([]txResponse, len(txs))
	for i, t := range txs {
		responses[i] = toResponse(t)
	}
	page := pg.Page
	if page < 1 {
		page = 1
	}
	pageSize := pg.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{
		"items":    responses,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

func (h *Handler) GetTransaction(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	tx, err := h.txRepo.GetByID(c.Request.Context(), id)
	if err == common.ErrNotFound {
		apiresponse.WriteError(c, apperrors.New("NOT_FOUND", "تراکنش یافت نشد", 404, apperrors.ErrNotFound))
		return
	}
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, toResponse(*tx))
}

func (h *Handler) ListAllTransactions(c *gin.Context) {
	var pg common.PaginationParams
	if err := c.ShouldBindQuery(&pg); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	txs, total, err := h.txRepo.ListAll(c.Request.Context(), domain.TransactionFilter{}, pg)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	responses := make([]txResponse, len(txs))
	for i, t := range txs {
		responses[i] = toResponse(t)
	}
	page := pg.Page
	if page < 1 {
		page = 1
	}
	pageSize := pg.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{
		"items":    responses,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

type txResponse struct {
	ID              string `json:"id"`
	Amount          int64  `json:"amount"`
	Currency        string `json:"currency"`
	Type            string `json:"type"`
	Status          string `json:"status"`
	ReferenceNumber string `json:"referenceNumber"`
	Description     string `json:"description"`
	Date            string `json:"date"`
}

func toResponse(t domain.Transaction) txResponse {
	return txResponse{
		ID:              t.ID.String(),
		Amount:          t.Amount,
		Currency:        t.Currency,
		Type:            t.Type.APIType(),
		Status:          t.Status.APIStatus(),
		ReferenceNumber: t.ReferenceNumber,
		Description:     t.Description,
		Date:            t.CreatedAt.Format("2006-01-02T15:04:05Z"),
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
