package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/younesbeheshti/any-task-connect/backend/pkg/jwt"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/response"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

const (
	UserIDKey = "user_id"
	RoleKey   = "role"
	JTIKey    = "jti"
)

// Auth validates JWT access tokens.
func Auth(jwtService *jwtpkg.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c.GetHeader("Authorization"))
		if token == "" {
			response.ErrorResponse(c, apperrors.New("UNAUTHORIZED", "missing authorization token", 401, apperrors.ErrUnauthorized))
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			response.ErrorResponse(c, apperrors.New("UNAUTHORIZED", "invalid or expired token", 401, apperrors.ErrUnauthorized))
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(RoleKey, claims.Role)
		c.Set(JTIKey, claims.JTI)
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

// OptionalAuth attaches user context when a valid token is present.
func OptionalAuth(jwtService *jwtpkg.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c.GetHeader("Authorization"))
		if token == "" {
			c.Next()
			return
		}

		claims, err := jwtService.ValidateAccessToken(token)
		if err == nil {
			c.Set(UserIDKey, claims.UserID)
			c.Set(RoleKey, claims.Role)
			c.Set(JTIKey, claims.JTI)
		}
		c.Next()
	}
}
