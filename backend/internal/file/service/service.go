package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/file/domain"
	filerepo "github.com/younesbeheshti/any-task-connect/backend/internal/file/repository"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

// Service defines file upload use cases.
type Service interface {
	Upload(ctx context.Context, userID uuid.UUID, fh *multipart.FileHeader) (*domain.File, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.File, error)
}

type FileService struct {
	repo     filerepo.Repository
	localDir string
	maxSize  int64
}

func NewFileService(repo filerepo.Repository, localDir string, maxSize int64) *FileService {
	return &FileService{repo: repo, localDir: localDir, maxSize: maxSize}
}

func (s *FileService) Upload(ctx context.Context, userID uuid.UUID, fh *multipart.FileHeader) (*domain.File, error) {
	if fh == nil {
		return nil, apperrors.New("INVALID_FILE", "فایلی ارسال نشده است", 422, apperrors.ErrValidation)
	}
	if s.maxSize > 0 && fh.Size > s.maxSize {
		return nil, apperrors.New("FILE_TOO_LARGE", "حجم فایل بیش از حد مجاز است", 422, apperrors.ErrValidation)
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fh.Filename), "."))
	if !domain.ExtAllowed(ext) {
		return nil, apperrors.New("UNSUPPORTED_TYPE", "نوع فایل مجاز نیست", 422, apperrors.ErrValidation)
	}

	src, err := fh.Open()
	if err != nil {
		return nil, apperrors.New("INVALID_FILE", "خواندن فایل ناموفق بود", 422, apperrors.ErrValidation)
	}
	defer src.Close()

	if err := os.MkdirAll(s.localDir, 0o755); err != nil {
		return nil, fmt.Errorf("create upload dir: %w", err)
	}

	storedName := uuid.New().String()
	if ext != "" {
		storedName += "." + ext
	}
	fullPath := filepath.Join(s.localDir, storedName)

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	written, copyErr := io.Copy(dst, src)
	closeErr := dst.Close()
	if copyErr != nil {
		_ = os.Remove(fullPath)
		return nil, fmt.Errorf("write file: %w", copyErr)
	}
	if closeErr != nil {
		_ = os.Remove(fullPath)
		return nil, fmt.Errorf("close file: %w", closeErr)
	}

	mime := fh.Header.Get("Content-Type")
	if mime == "" {
		mime = "application/octet-stream"
	}

	f := &domain.File{
		ID:           uuid.New(),
		OriginalName: fh.Filename,
		StoredName:   storedName,
		Path:         fullPath,
		Size:         written,
		MimeType:     mime,
		UploadedBy:   userID,
	}
	if err := s.repo.Create(ctx, f); err != nil {
		_ = os.Remove(fullPath)
		return nil, err
	}
	return f, nil
}

func (s *FileService) GetByID(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	return s.repo.GetByID(ctx, id)
}
