package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/category/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/category/repository"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

// CategoryService implements service.Service.
type CategoryService struct {
	repo repository.Repository
}

// NewCategoryService creates a new CategoryService.
func NewCategoryService(repo repository.Repository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) List(ctx context.Context) ([]domain.Category, error) {
	return s.repo.List(ctx)
}

func (s *CategoryService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "دسته‌بندی یافت نشد", 404, apperrors.ErrNotFound)
	}
	return cat, err
}

func (s *CategoryService) Create(ctx context.Context, title, icon string) (*domain.Category, error) {
	existing, err := s.repo.GetByTitle(ctx, title)
	if err == nil && existing != nil {
		return nil, apperrors.New("CONFLICT", "دسته‌بندی با این نام وجود دارد", 409, apperrors.ErrConflict)
	}
	cat := &domain.Category{
		ID:        uuid.New(),
		Title:     title,
		Icon:      icon,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) Update(ctx context.Context, id uuid.UUID, title, icon string) (*domain.Category, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "دسته‌بندی یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	cat.Title = title
	cat.Icon = icon
	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, common.ErrNotFound) {
		return apperrors.New("NOT_FOUND", "دسته‌بندی یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
