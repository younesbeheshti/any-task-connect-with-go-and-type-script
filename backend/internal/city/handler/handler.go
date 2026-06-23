package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/city/service"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
)

// Handler serves city HTTP endpoints.
type Handler struct {
	svc *service.CityService
}

// NewHandler creates a city Handler.
func NewHandler(svc *service.CityService) *Handler {
	return &Handler{svc: svc}
}

type createCityRequest struct {
	Title    string `json:"title" binding:"required"`
	Province string `json:"province"`
}

type updateCityRequest struct {
	Title    string `json:"title" binding:"required"`
	Province string `json:"province"`
}


func (h *Handler) List(c *gin.Context) {
	cities, err := h.svc.List(c.Request.Context())
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"cities": cities})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.JSON(c, http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_ID", "message": "شناسه نامعتبر است"}})
		return
	}
	city, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, city)
}

func (h *Handler) Create(c *gin.Context) {
	var req createCityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	city, err := h.svc.Create(c.Request.Context(), req.Title, req.Province)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusCreated, city)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.JSON(c, http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_ID", "message": "شناسه نامعتبر است"}})
		return
	}
	var req updateCityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	city, err := h.svc.Update(c.Request.Context(), id, req.Title, req.Province)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, city)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.JSON(c, http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_ID", "message": "شناسه نامعتبر است"}})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.NoContent(c)
}
