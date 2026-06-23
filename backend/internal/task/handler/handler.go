package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	"github.com/younesbeheshti/any-task-connect/backend/internal/task/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/task/service"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

// Handler serves task HTTP endpoints.
type Handler struct {
	svc *service.TaskService
}

// NewHandler creates a task Handler.
func NewHandler(svc *service.TaskService) *Handler {
	return &Handler{svc: svc}
}


type createTaskRequest struct {
	Title          string   `json:"title" binding:"required"`
	Description    string   `json:"description" binding:"required"`
	CategoryID     string   `json:"categoryId" binding:"required"`
	CityID         string   `json:"cityId" binding:"required"`
	Budget         int64    `json:"budget" binding:"required,min=1"`
	Deadline       string   `json:"deadline" binding:"required"`
	AttachmentURLs []string `json:"attachmentUrls"`
}

type updateTaskRequest struct {
	Title          *string  `json:"title"`
	Description    *string  `json:"description"`
	CategoryID     *string  `json:"categoryId"`
	CityID         *string  `json:"cityId"`
	Budget         *int64   `json:"budget"`
	Deadline       *string  `json:"deadline"`
	AttachmentURLs []string `json:"attachmentUrls"`
}

func (h *Handler) Create(c *gin.Context) {
	requesterID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_FIELD", "شناسه دسته‌بندی نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	cityID, err := uuid.Parse(req.CityID)
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_FIELD", "شناسه شهر نامعتبر است", 400, apperrors.ErrValidation))
		return
	}
	deadline, err := time.Parse("2006-01-02", req.Deadline)
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_FIELD", "فرمت تاریخ نامعتبر است (YYYY-MM-DD)", 400, apperrors.ErrValidation))
		return
	}

	input := domain.CreateTaskInput{
		Title:          req.Title,
		Description:    req.Description,
		CategoryID:     categoryID,
		CityID:         cityID,
		Budget:         req.Budget,
		Deadline:       deadline,
		AttachmentURLs: req.AttachmentURLs,
		RequesterID:    requesterID,
	}

	task, escrow, err := h.svc.Create(c.Request.Context(), input)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}

	apiresponse.JSON(c, http.StatusCreated, CreateTaskResponse{
		Task: toTaskResponse(task),
		Escrow: EscrowResponse{
			Fee:      escrow.Fee,
			Held:     escrow.Held,
			Currency: escrow.Currency,
		},
	})
}

func (h *Handler) GetByPublicID(c *gin.Context) {
	publicID := c.Param("id")
	var viewerID *uuid.UUID
	if idStr := c.GetString(middleware.UserIDKey); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			viewerID = &id
		}
	}
	task, err := h.svc.GetByPublicID(c.Request.Context(), publicID, viewerID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, toTaskResponse(task))
}

func (h *Handler) List(c *gin.Context) {
	var pg common.PaginationParams
	if err := c.ShouldBindQuery(&pg); err != nil {
		apiresponse.WriteError(c, err)
		return
	}

	filter := domain.TaskFilter{
		Query:    c.Query("q"),
		Sort:     c.Query("sort"),
		CityTitle: c.Query("city"),
		CatTitle: c.Query("category"),
	}

	if cityStr := c.Query("cityId"); cityStr != "" {
		if id, err := uuid.Parse(cityStr); err == nil {
			filter.CityID = &id
		}
	}
	if catStr := c.Query("categoryId"); catStr != "" {
		if id, err := uuid.Parse(catStr); err == nil {
			filter.CategoryID = &id
		}
	}
	if statusStr := c.Query("status"); statusStr != "" {
		s := domain.TaskStatus(statusStr)
		filter.Status = &s
	}
	if minStr := c.Query("minBudget"); minStr != "" {
		if v, err := strconv.ParseInt(minStr, 10, 64); err == nil {
			filter.MinBudget = &v
		}
	}
	if maxStr := c.Query("maxBudget"); maxStr != "" {
		if v, err := strconv.ParseInt(maxStr, 10, 64); err == nil {
			filter.MaxBudget = &v
		}
	}

	result, err := h.svc.List(c.Request.Context(), filter, pg)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}

	responses := make([]TaskResponse, len(result.Items))
	for i, t := range result.Items {
		responses[i] = toTaskResponse(&t)
	}

	apiresponse.JSON(c, http.StatusOK, gin.H{
		"items":    responses,
		"page":     result.Page,
		"pageSize": result.PageSize,
		"total":    result.Total,
	})
}

func (h *Handler) Update(c *gin.Context) {
	requesterID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	publicID := c.Param("id")
	task, err := h.svc.GetByPublicID(c.Request.Context(), publicID, nil)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}

	var req updateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.WriteError(c, err)
		return
	}

	input := domain.UpdateTaskInput{
		Title:          req.Title,
		Description:    req.Description,
		Budget:         req.Budget,
		AttachmentURLs: req.AttachmentURLs,
	}
	if req.CategoryID != nil {
		id, err := uuid.Parse(*req.CategoryID)
		if err == nil {
			input.CategoryID = &id
		}
	}
	if req.CityID != nil {
		id, err := uuid.Parse(*req.CityID)
		if err == nil {
			input.CityID = &id
		}
	}
	if req.Deadline != nil {
		d, err := time.Parse("2006-01-02", *req.Deadline)
		if err == nil {
			input.Deadline = &d
		}
	}

	updated, err := h.svc.Update(c.Request.Context(), task.ID, requesterID, input)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, toTaskResponse(updated))
}

func (h *Handler) Delete(c *gin.Context) {
	requesterID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	publicID := c.Param("id")
	task, err := h.svc.GetByPublicID(c.Request.Context(), publicID, nil)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), task.ID, requesterID); err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.NoContent(c)
}

func (h *Handler) Cancel(c *gin.Context) {
	actorID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	publicID := c.Param("id")
	task, err := h.svc.GetByPublicID(c.Request.Context(), publicID, nil)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	updated, err := h.svc.Cancel(c.Request.Context(), task.ID, actorID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, toTaskResponse(updated))
}

func (h *Handler) Start(c *gin.Context) {
	agentID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	publicID := c.Param("id")
	task, err := h.svc.GetByPublicID(c.Request.Context(), publicID, nil)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	updated, err := h.svc.Start(c.Request.Context(), task.ID, agentID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, toTaskResponse(updated))
}

func (h *Handler) Complete(c *gin.Context) {
	agentID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	publicID := c.Param("id")
	task, err := h.svc.GetByPublicID(c.Request.Context(), publicID, nil)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	updated, err := h.svc.Complete(c.Request.Context(), task.ID, agentID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, toTaskResponse(updated))
}

func (h *Handler) Verify(c *gin.Context) {
	requesterID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	publicID := c.Param("id")
	task, err := h.svc.GetByPublicID(c.Request.Context(), publicID, nil)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	updated, err := h.svc.Verify(c.Request.Context(), task.ID, requesterID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, toTaskResponse(updated))
}

func (h *Handler) GetTimeline(c *gin.Context) {
	publicID := c.Param("id")
	task, err := h.svc.GetByPublicID(c.Request.Context(), publicID, nil)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	timeline, err := h.svc.GetTimeline(c.Request.Context(), task.ID)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"timeline": timeline})
}

func currentUserID(c *gin.Context) (uuid.UUID, error) {
	idStr := c.GetString(middleware.UserIDKey)
	if idStr == "" {
		return uuid.Nil, apperrors.New("UNAUTHORIZED", "لطفاً وارد حساب کاربری شوید", 401, apperrors.ErrUnauthorized)
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, apperrors.New("UNAUTHORIZED", "شناسه کاربر نامعتبر است", 401, apperrors.ErrUnauthorized)
	}
	return id, nil
}
