package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/younesbeheshti/any-task-connect/backend/internal/auth/domain"
	authservice "github.com/younesbeheshti/any-task-connect/backend/internal/auth/service"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	userdomain "github.com/younesbeheshti/any-task-connect/backend/internal/user/domain"
	userhandler "github.com/younesbeheshti/any-task-connect/backend/internal/user/handler"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/validator"
)

const refreshCookieName = "refresh_token"

// Handler serves authentication HTTP endpoints.
type Handler struct {
	auth      *authservice.AuthService
	validator *validator.Validator
}

// NewHandler creates a Handler.
func NewHandler(auth *authservice.AuthService, v *validator.Validator) *Handler {
	return &Handler{auth: auth, validator: v}
}

// RegisterRoutes mounts auth routes.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/auth/register", h.Register)
	rg.POST("/auth/login", h.Login)
	rg.POST("/auth/refresh", h.Refresh)
	rg.POST("/auth/logout", middleware.AuthJWT(h.auth), h.Logout)
	rg.POST("/auth/logout-all", middleware.AuthJWT(h.auth), h.LogoutAll)
	rg.POST("/auth/change-password", middleware.AuthJWT(h.auth), h.ChangePassword)
	rg.POST("/auth/forgot-password", h.ForgotPassword)
	rg.POST("/auth/reset-password", h.ResetPassword)
	rg.POST("/auth/otp/send", h.SendOTP)
	rg.POST("/auth/otp/verify", h.VerifyOTP)
}

// RegisterRequest matches POST /auth/register contract.
type RegisterRequest struct {
	FullName string `json:"fullName" validate:"required,min=2,max=200"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required,password_strength"`
	Role     string `json:"role" validate:"required,oneof=requester agent"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

// LoginRequest matches POST /auth/login contract (+ optional email).
type LoginRequest struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
}

// OTPRequest matches OTP send contract.
type OTPRequest struct {
	Phone string `json:"phone" validate:"required"`
}

// OTPVerifyRequest matches OTP verify contract.
type OTPVerifyRequest struct {
	Phone string `json:"phone" validate:"required"`
	Code  string `json:"code" validate:"required,len=6"`
}

// ChangePasswordRequest holds password change input.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,password_strength"`
}

// ForgotPasswordRequest holds forgot password input.
type ForgotPasswordRequest struct {
	Phone string `json:"phone" validate:"required"`
}

// ResetPasswordRequest holds reset password input.
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,password_strength"`
}

// RefreshRequest holds refresh token body fallback.
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// AuthResponse matches register/login contract with Phase 3 extensions.
type AuthResponse struct {
	Token        string              `json:"token"`
	RefreshToken string              `json:"refreshToken,omitempty"`
	User         userhandler.UserDTO `json:"user"`
	Permissions  []string            `json:"permissions,omitempty"`
	Role         string              `json:"role,omitempty"`
}

// Register godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  RegisterRequest  true  "Registration data"
// @Success      201   {object}  dto.AuthResponse
// @Router       /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}

	if err := h.validator.ValidateStruct(req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}

	var email *string
	if req.Email != "" {
		email = &req.Email
	}
	tokens, user, err := h.auth.Register(c.Request.Context(), domain.RegisterInput{
		FullName: req.FullName, Phone: req.Phone, Password: req.Password,
		Role: common.RoleFromAPI(req.Role), Email: email,
	})

	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	h.writeAuthResponse(c, http.StatusCreated, tokens, user)
}

// Login godoc
// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  LoginRequest  true  "Credentials"
// @Success      200   {object}  dto.AuthResponse
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	if req.Phone == "" && req.Email == "" {
		apiresponse.WriteError(c, apperrors.Validation(map[string]string{
			"phone": "شماره موبایل یا ایمیل الزامی است",
		}))
		return
	}
	tokens, user, err := h.auth.Login(c.Request.Context(), domain.LoginInput{
		Phone: req.Phone, Email: req.Email, Password: req.Password,
	})
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	h.writeAuthResponse(c, http.StatusOK, tokens, user)
}

// Refresh godoc
// @Summary      Refresh access token
// @Tags         auth
// @Success      200  {object}  dto.AuthResponse
// @Router       /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, _ := c.Cookie(refreshCookieName)
	if refreshToken == "" {
		var req RefreshRequest
		_ = c.ShouldBindJSON(&req)
		refreshToken = req.RefreshToken
	}
	if refreshToken == "" {
		apiresponse.WriteError(c, apperrors.New("UNAUTHORIZED", "توکن تازه‌سازی یافت نشد", 401, common.ErrUnauthorized))
		return
	}
	tokens, err := h.auth.Refresh(c.Request.Context(), domain.RefreshInput{RefreshToken: refreshToken})
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	h.setRefreshCookie(c, tokens.RefreshToken)
	apiresponse.JSON(c, http.StatusOK, gin.H{
		"token":        tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
	})
}

// Logout godoc
// @Summary      Logout current device
// @Tags         auth
// @Security     BearerAuth
// @Success      204
// @Router       /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	userID, jti := middleware.GetAuthContext(c)
	_ = h.auth.Logout(c.Request.Context(), userID, jti)
	h.clearRefreshCookie(c)
	apiresponse.NoContent(c)
}

// LogoutAll godoc
// @Summary      Logout all devices
// @Tags         auth
// @Security     BearerAuth
// @Success      204
// @Router       /auth/logout-all [post]
func (h *Handler) LogoutAll(c *gin.Context) {
	userID, jti := middleware.GetAuthContext(c)
	_ = h.auth.LogoutAll(c.Request.Context(), userID, jti)
	h.clearRefreshCookie(c)
	apiresponse.NoContent(c)
}

// ChangePassword godoc
// @Summary      Change password
// @Tags         auth
// @Security     BearerAuth
// @Success      204
// @Router       /auth/change-password [post]
func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	if err := h.validator.ValidateStruct(req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	userID, _ := middleware.GetAuthContext(c)
	if err := h.auth.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.NoContent(c)
}

// ForgotPassword godoc
// @Summary      Request password reset
// @Tags         auth
// @Success      204
// @Router       /auth/forgot-password [post]
func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	_ = h.auth.ForgotPassword(c.Request.Context(), domain.ForgotPasswordInput{Phone: req.Phone})
	apiresponse.NoContent(c)
}

// ResetPassword godoc
// @Summary      Reset password with token
// @Tags         auth
// @Success      204
// @Router       /auth/reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	if err := h.validator.ValidateStruct(req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	if err := h.auth.ResetPassword(c.Request.Context(), domain.ResetPasswordInput{
		Token: req.Token, NewPassword: req.NewPassword,
	}); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.NoContent(c)
}

// SendOTP godoc
// @Summary      Send phone OTP
// @Tags         auth
// @Success      204
// @Router       /auth/otp/send [post]
func (h *Handler) SendOTP(c *gin.Context) {
	var req OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	if err := h.auth.SendPhoneOTP(c.Request.Context(), req.Phone); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.NoContent(c)
}

// VerifyOTP godoc
// @Summary      Verify phone OTP
// @Tags         auth
// @Success      204
// @Router       /auth/otp/verify [post]
func (h *Handler) VerifyOTP(c *gin.Context) {
	var req OTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	if err := h.auth.VerifyPhoneOTP(c.Request.Context(), domain.OTPInput{
		Phone: req.Phone, Code: req.Code,
	}); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.NoContent(c)
}

func (h *Handler) writeAuthResponse(c *gin.Context, status int, tokens *domain.TokenPair, user *userdomain.User) {
	h.setRefreshCookie(c, tokens.RefreshToken)
	apiresponse.JSON(c, status, AuthResponse{
		Token:        tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         userhandler.ToUserDTO(user),
		Permissions:  common.PermissionsForRole(user.Role),
		Role:         user.Role.APIRole(),
	})
}

func (h *Handler) setRefreshCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(refreshCookieName, token, int((7 * 24 * time.Hour).Seconds()), "/", "", false, true)
}

func (h *Handler) clearRefreshCookie(c *gin.Context) {
	c.SetCookie(refreshCookieName, "", -1, "/", "", false, true)
}
