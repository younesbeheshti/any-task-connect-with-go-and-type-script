package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/category/service"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
)


// Handler serves category HTTP endpoints.
type Handler struct {
	svc *service.CategoryService
}

// NewHandler creates a category Handler.
func NewHandler(svc *service.CategoryService) *Handler {
	return &Handler{svc: svc}
}

type createCategoryRequest struct {
	Title string `json:"title" binding:"required"`
	Icon  string `json:"icon"`
}

type updateCategoryRequest struct {
	Title string `json:"title" binding:"required"`
	Icon  string `json:"icon"`
}


func (h *Handler) List(c *gin.Context) {
	cats, err := h.svc.List(c.Request.Context())
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"categories": cats})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.JSON(c, http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_ID", "message": "شناسه نامعتبر است"}})
		return
	}
	cat, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, cat)
}

func (h *Handler) Create(c *gin.Context) {
	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	cat, err := h.svc.Create(c.Request.Context(), req.Title, req.Icon)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusCreated, cat)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.JSON(c, http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_ID", "message": "شناسه نامعتبر است"}})
		return
	}
	var req updateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	cat, err := h.svc.Update(c.Request.Context(), id, req.Title, req.Icon)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, cat)
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
