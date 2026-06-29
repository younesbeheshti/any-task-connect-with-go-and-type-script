package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	authservice "github.com/younesbeheshti/any-task-connect/backend/internal/auth/service"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

// AuthJWT validates JWT access tokens with blacklist check.
func AuthJWT(auth *authservice.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c.GetHeader("Authorization"))
		if token == "" {
			apiresponse.WriteError(c, apperrors.New("UNAUTHORIZED", "لطفاً وارد حساب کاربری شوید", 401, apperrors.ErrUnauthorized))
			c.Abort()
			return
		}
		claims, err := auth.ValidateAccessToken(c.Request.Context(), token)
		if err != nil {
			apiresponse.WriteError(c, apperrors.New("UNAUTHORIZED", "توکن نامعتبر یا منقضی شده", 401, apperrors.ErrUnauthorized))
			c.Abort()
			return
		}
		c.Set(UserIDKey, claims.UserID.String())
		c.Set(RoleKey, claims.Role)
		c.Set(JTIKey, claims.JTI)
		c.Next()
	}
}

// OptionalAuthJWT sets the user context when a valid token is present, but never
// aborts when it is missing or invalid. Use on public routes that adjust their
// response for authenticated users (e.g. "mine" filters on the task list).
func OptionalAuthJWT(auth *authservice.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c.GetHeader("Authorization"))
		if token == "" {
			c.Next()
			return
		}
		claims, err := auth.ValidateAccessToken(c.Request.Context(), token)
		if err != nil {
			c.Next()
			return
		}
		c.Set(UserIDKey, claims.UserID.String())
		c.Set(RoleKey, claims.Role)
		c.Set(JTIKey, claims.JTI)
		c.Next()
	}
}

// GetAuthContext returns authenticated user ID and JTI.
func GetAuthContext(c *gin.Context) (uuid.UUID, string) {
	userID, _ := uuid.Parse(c.GetString(UserIDKey))
	return userID, c.GetString(JTIKey)
}

// CurrentUser loads the authenticated user into context (optional enrichment).
func CurrentUser(auth *authservice.AuthService) gin.HandlerFunc {
	return AuthJWT(auth)
}

// PermissionRequired checks granular permission for the current role.
func PermissionRequired(perm common.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(RoleKey)
		if !exists {
			apiresponse.WriteError(c, apperrors.New("UNAUTHORIZED", "لطفاً وارد حساب کاربری شوید", 401, apperrors.ErrUnauthorized))
			c.Abort()
			return
		}
		role, ok := roleVal.(common.Role)
		if !ok || !common.HasPermission(role, perm) {
			apiresponse.WriteError(c, apperrors.New("FORBIDDEN", "دسترسی مجاز نیست", 403, apperrors.ErrForbidden))
			c.Abort()
			return
		}
		c.Next()
	}
}

func extractBearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
