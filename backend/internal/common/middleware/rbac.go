package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/response"
)

// RBAC checks granular permissions for the authenticated role.
func RBAC(permissions ...common.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(RoleKey)
		if !exists {
			response.ErrorResponse(c, apperrors.New("UNAUTHORIZED", "authentication required", 401, apperrors.ErrUnauthorized))
			c.Abort()
			return
		}

		role, ok := roleVal.(common.Role)
		if !ok {
			response.ErrorResponse(c, apperrors.New("FORBIDDEN", "invalid role", 403, apperrors.ErrForbidden))
			c.Abort()
			return
		}

		for _, perm := range permissions {
			if !common.HasPermission(role, perm) {
				response.ErrorResponse(c, apperrors.New("FORBIDDEN", "permission denied", 403, apperrors.ErrForbidden))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
