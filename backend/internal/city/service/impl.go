package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/city/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/city/repository"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

// CityService implements service.Service.
type CityService struct {
	repo repository.Repository
}

// NewCityService creates a new CityService.
func NewCityService(repo repository.Repository) *CityService {
	return &CityService{repo: repo}
}

func (s *CityService) List(ctx context.Context) ([]domain.City, error) {
	return s.repo.List(ctx)
}

func (s *CityService) GetByID(ctx context.Context, id uuid.UUID) (*domain.City, error) {
	city, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "شهر یافت نشد", 404, apperrors.ErrNotFound)
	}
	return city, err
}

func (s *CityService) Create(ctx context.Context, title, province string) (*domain.City, error) {
	existing, err := s.repo.GetByTitle(ctx, title)
	if err == nil && existing != nil {
		return nil, apperrors.New("CONFLICT", "شهر با این نام وجود دارد", 409, apperrors.ErrConflict)
	}
	city := &domain.City{
		ID:        uuid.New(),
		Title:     title,
		Province:  province,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, city); err != nil {
		return nil, err
	}
	return city, nil
}

func (s *CityService) Update(ctx context.Context, id uuid.UUID, title, province string) (*domain.City, error) {
	city, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "شهر یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	city.Title = title
	city.Province = province
	if err := s.repo.Update(ctx, city); err != nil {
		return nil, err
	}
	return city, nil
}

func (s *CityService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, common.ErrNotFound) {
		return apperrors.New("NOT_FOUND", "شهر یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
