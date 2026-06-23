package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/user/domain"
	userinfra "github.com/younesbeheshti/any-task-connect/backend/internal/user/infra"
	"github.com/younesbeheshti/any-task-connect/backend/internal/user/repository"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

// UserService implements user profile use cases.
type UserService struct {
	repo repository.Repository
	city *userinfra.GormRepository
}

// NewUserService creates a UserService.
func NewUserService(repo repository.Repository) *UserService {
	cityRepo, _ := repo.(*userinfra.GormRepository)
	return &UserService{repo: repo, city: cityRepo}
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetPublicProfile(ctx context.Context, id uuid.UUID) (*domain.PublicProfile, error) {
	user, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "کاربر یافت نشد", 404, err)
	}
	if err != nil {
		return nil, err
	}
	profile := user.ToPublicProfile()
	return &profile, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, input domain.UpdateProfileInput) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if input.FullName != nil {
		user.FullName = *input.FullName
	}
	if input.Email != nil {
		user.Email = input.Email
	}
	if input.Avatar != nil {
		user.Avatar = input.Avatar
	}
	if input.Bio != nil {
		user.Bio = input.Bio
	}
	if input.CityID != nil {
		user.CityID = input.CityID
	}
	if input.NationalID != nil {
		user.NationalID = input.NationalID
	}
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, userID)
}

func (s *UserService) UpdateProfileByCityTitle(ctx context.Context, userID uuid.UUID, input domain.UpdateProfileInput, cityTitle *string) (*domain.User, error) {
	if cityTitle != nil && s.city != nil {
		cityID, err := s.city.GetCityByTitle(ctx, *cityTitle)
		if errors.Is(err, common.ErrNotFound) {
			return nil, apperrors.Validation(map[string]string{"city": "شهر انتخاب شده معتبر نیست"})
		}
		if err != nil {
			return nil, err
		}
		input.CityID = cityID
	}
	return s.UpdateProfile(ctx, userID, input)
}

func (s *UserService) DeleteAvatar(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	empty := ""
	return s.UpdateProfile(ctx, userID, domain.UpdateProfileInput{Avatar: &empty})
}

func (s *UserService) GetStats(ctx context.Context, userID uuid.UUID) (*domain.UserStats, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &domain.UserStats{
		Rating:         user.Rating,
		CompletedCount: user.CompletedTasks,
		RatingCount:    user.RatingCount,
		WalletBalance:  user.WalletBalance,
		LockedBalance:  user.LockedBalance,
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context, filter repository.UserFilter, pg common.PaginationParams) (*common.PaginatedResult[domain.User], error) {
	items, total, err := s.repo.List(ctx, filter, pg)
	if err != nil {
		return nil, err
	}
	return &common.PaginatedResult[domain.User]{
		Items: items, Page: pg.Page, PageSize: pg.PageSize, Total: total,
	}, nil
}

func (s *UserService) Deactivate(ctx context.Context, userID uuid.UUID) error {
	return s.repo.SoftDelete(ctx, userID)
}

func (s *UserService) UpdateRatingAggregate(ctx context.Context, userID uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, userID)
	return err
}
