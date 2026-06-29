package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	"github.com/younesbeheshti/any-task-connect/backend/internal/file/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/file/service"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/apiresponse"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// uploadedFile is the metadata returned per file. The shape matches the chat
// domain.Attachment value object so chat/task forms can use it directly.
type uploadedFile struct {
	ID   uuid.UUID `json:"id"`
	URL  string    `json:"url"`
	Name string    `json:"name"`
	Mime string    `json:"mime"`
	Size int64     `json:"size"`
}

func toUploaded(f *domain.File) uploadedFile {
	return uploadedFile{
		ID:   f.ID,
		URL:  "/v1/files/" + f.ID.String(),
		Name: f.OriginalName,
		Mime: f.MimeType,
		Size: f.Size,
	}
}

// Upload accepts one or more files via multipart/form-data under the "files" field.
func (h *Handler) Upload(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("INVALID_FILE", "بارگذاری فایل نامعتبر است", 422, apperrors.ErrValidation))
		return
	}
	headers := form.File["files"]
	if len(headers) == 0 {
		apiresponse.WriteError(c, apperrors.New("INVALID_FILE", "فایلی ارسال نشده است", 422, apperrors.ErrValidation))
		return
	}

	results := make([]uploadedFile, 0, len(headers))
	for _, fh := range headers {
		f, err := h.svc.Upload(c.Request.Context(), userID, fh)
		if err != nil {
			apiresponse.WriteError(c, err)
			return
		}
		results = append(results, toUploaded(f))
	}

	apiresponse.JSON(c, http.StatusOK, gin.H{"files": results})
}

// Download streams a stored file with its original filename.
func (h *Handler) Download(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.WriteError(c, apperrors.New("NOT_FOUND", "فایل یافت نشد", 404, apperrors.ErrNotFound))
		return
	}
	f, err := h.svc.GetByID(c.Request.Context(), id)
	if err == common.ErrNotFound {
		apiresponse.WriteError(c, apperrors.New("NOT_FOUND", "فایل یافت نشد", 404, apperrors.ErrNotFound))
		return
	}
	if err != nil {
		apiresponse.WriteError(c, err)
		return
	}
	c.FileAttachment(f.Path, f.OriginalName)
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
