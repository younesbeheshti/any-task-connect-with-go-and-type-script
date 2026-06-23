package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	"github.com/younesbeheshti/any-task-connect/backend/internal/user/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/user/service"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/validator"
)

// Handler serves user profile HTTP endpoints.
type Handler struct {
	users     *service.UserService
	validator *validator.Validator
}

// NewHandler creates a Handler.
func NewHandler(users *service.UserService, v *validator.Validator) *Handler {
	return &Handler{users: users, validator: v}
}

// RegisterRoutes mounts user routes on the router group.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/me", h.GetMe)
	rg.PATCH("/me", h.UpdateMe)
	rg.DELETE("/me/avatar", h.DeleteAvatar)
	rg.GET("/me/stats", h.GetStats)
	rg.GET("/users/:id/public", h.GetPublicProfile)
}

// GetMe godoc
// @Summary      Get current user
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  dto.User
// @Failure      401  {object}  apiresponse.ErrorBody
// @Router       /me [get]
func (h *Handler) GetMe(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	user, err := h.users.GetByID(c.Request.Context(), userID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, ToUserDTO(user))
}

// UpdateMe godoc
// @Summary      Update current user profile
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  UpdateProfileRequest  true  "Profile fields"
// @Success      200   {object}  dto.User
// @Router       /me [patch]
func (h *Handler) UpdateMe(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	input := domain.UpdateProfileInput{
		FullName: req.FullName, Avatar: req.AvatarURL, Bio: req.Bio,
	}
	user, err := h.users.UpdateProfileByCityTitle(c.Request.Context(), userID, input, req.City)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, ToUserDTO(user))
}

// DeleteAvatar godoc
// @Summary      Delete user avatar
// @Tags         users
// @Security     BearerAuth
// @Success      200  {object}  dto.User
// @Router       /me/avatar [delete]
func (h *Handler) DeleteAvatar(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	user, err := h.users.DeleteAvatar(c.Request.Context(), userID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, ToUserDTO(user))
}

// GetStats godoc
// @Summary      Get user statistics
// @Tags         users
// @Security     BearerAuth
// @Success      200  {object}  dto.UserStats
// @Router       /me/stats [get]
func (h *Handler) GetStats(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	stats, err := h.users.GetStats(c.Request.Context(), userID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, stats)
}

// GetPublicProfile godoc
// @Summary      Get public user profile
// @Tags         users
// @Produce      json
// @Param        id  path  string  true  "User ID"
// @Success      200  {object}  dto.User
// @Router       /users/{id}/public [get]
func (h *Handler) GetPublicProfile(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	profile, err := h.users.GetPublicProfile(c.Request.Context(), id)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, ToPublicUserDTO(profile))
}

func currentUserID(c *gin.Context) (uuid.UUID, error) {
	idStr := c.GetString(middleware.UserIDKey)
	if idStr == "" {
		return uuid.Nil, apperrors.New("UNAUTHORIZED", "لطفاً وارد حساب کاربری شوید", 401, apperrors.ErrUnauthorized)
	}
	return uuid.Parse(idStr)
}
