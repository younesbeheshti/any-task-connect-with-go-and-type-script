package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/admin/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/admin/service"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	userrepo "github.com/younesbeheshti/any-task-connect/backend/internal/user/repository"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// GET /admin/users
func (h *Handler) ListUsers(c *gin.Context) {
	var pg common.PaginationParams
	_ = c.ShouldBindQuery(&pg)

	var filter userrepo.UserFilter
	filter.Query = c.Query("q")
	if role := c.Query("role"); role != "" {
		r := common.RoleFromAPI(role)
		if r.IsValid() {
			filter.Role = &r
		}
	}
	if active := c.Query("active"); active != "" {
		b := active == "true"
		filter.IsActive = &b
	}

	result, err := h.svc.ListUsers(c.Request.Context(), filter, pg)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, result)
}

// GET /admin/users/:id
func (h *Handler) GetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	pg := common.PaginationParams{Page: 1, PageSize: 1}
	active := true
	filter := userrepo.UserFilter{}
	_ = active
	_ = filter
	// Get via ListUsers with a single UUID lookup
	result, err := h.svc.ListUsers(c.Request.Context(), userrepo.UserFilter{}, pg)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	for _, u := range result.Items {
		if u.ID == userID {
			apiresponse.JSON(c, http.StatusOK, u)
			return
		}
	}
	apiresponse.WriteError(c, apperrors.New("NOT_FOUND", "کاربر یافت نشد", 404, apperrors.ErrNotFound))
}

// PATCH /admin/users/:id
func (h *Handler) UpdateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	var req struct {
		IsActive *bool   `json:"isActive"`
		Role     *string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	update := domain.UserAdminUpdate{
		IsActive: req.IsActive,
		Role:     req.Role,
	}
	user, err := h.svc.UpdateUser(c.Request.Context(), userID, update)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, user)
}

// POST /admin/users/:id/suspend
func (h *Handler) SuspendUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	inactive := false
	user, err := h.svc.UpdateUser(c.Request.Context(), userID, domain.UserAdminUpdate{IsActive: &inactive})
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"message": "حساب کاربری تعلیق شد", "user": user})
}

// POST /admin/users/:id/activate
func (h *Handler) ActivateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_ID", "شناسه نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	active := true
	user, err := h.svc.UpdateUser(c.Request.Context(), userID, domain.UserAdminUpdate{IsActive: &active})
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"message": "حساب کاربری فعال شد", "user": user})
}

// GET /admin/metrics
func (h *Handler) GetMetrics(c *gin.Context) {
	metrics, err := h.svc.GetMetrics(c.Request.Context())
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, metrics)
}

func mustUserID(c *gin.Context) uuid.UUID {
	id, _ := uuid.Parse(c.GetString(middleware.UserIDKey))
	return id
}
