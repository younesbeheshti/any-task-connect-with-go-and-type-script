package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	"github.com/younesbeheshti/any-task-connect/backend/internal/dashboard/service"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
)

// Handler serves dashboard HTTP endpoints.
type Handler struct {
	svc service.Service
}

// NewHandler creates a dashboard Handler.
func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}


func (h *Handler) GetAdminStats(c *gin.Context) {
	stats, err := h.svc.GetAdminStats(c.Request.Context())
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, stats)
}

func (h *Handler) GetUserStats(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey)
	stats, err := h.svc.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, stats)
}
