package middleware

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/response"
)

// Authorization ensures the authenticated user has one of the required roles.
func Authorization(roles ...common.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(RoleKey)
		if !exists {
			response.ErrorResponse(c, apperrors.New("UNAUTHORIZED", "authentication required", 401, apperrors.ErrUnauthorized))
			c.Abort()
			return
		}

		role, ok := roleVal.(common.Role)
		if !ok || !role.IsValid() {
			response.ErrorResponse(c, apperrors.New("FORBIDDEN", "invalid role", 403, apperrors.ErrForbidden))
			c.Abort()
			return
		}

		if len(roles) == 0 || slices.Contains(roles, role) {
			c.Next()
			return
		}

		response.ErrorResponse(c, apperrors.New("FORBIDDEN", "insufficient permissions", 403, apperrors.ErrForbidden))
		c.Abort()
	}
}
