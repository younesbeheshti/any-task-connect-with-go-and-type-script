package infra

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/user/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/user/repository"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/database"
	"gorm.io/gorm"
)

// GormRepository implements user repository with GORM.
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a GormRepository.
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, user *domain.User) error {
	m := toModel(user)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return common.ErrDuplicate
		}
		return err
	}
	created := toDomain(&m)
	*user = created
	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var m UserModel
	err := r.db.WithContext(ctx).Preload("City").First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	u := toDomain(&m)
	return &u, nil
}

func (r *GormRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var m UserModel
	err := r.db.WithContext(ctx).Preload("City").Where("phone = ?", phone).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	u := toDomain(&m)
	return &u, nil
}

func (r *GormRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m UserModel
	err := r.db.WithContext(ctx).Preload("City").Where("email = ?", email).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	u := toDomain(&m)
	return &u, nil
}

func (r *GormRepository) Update(ctx context.Context, user *domain.User) error {
	m := toModel(user)
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *GormRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	return r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", userID).
		Update("password_hash", passwordHash).Error
}

func (r *GormRepository) UpdateRating(ctx context.Context, userID uuid.UUID, rating float64, count int) error {
	return r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", userID).
		Updates(map[string]any{"rating": rating, "rating_count": count}).Error
}

func (r *GormRepository) IncrementCompletedTasks(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", userID).
		UpdateColumn("completed_tasks", gorm.Expr("completed_tasks + 1")).Error
}

func (r *GormRepository) List(ctx context.Context, filter repository.UserFilter, pg common.PaginationParams) ([]domain.User, int64, error) {
	q := r.db.WithContext(ctx).Model(&UserModel{})
	if filter.Query != "" {
		like := "%" + filter.Query + "%"
		q = q.Where("full_name ILIKE ? OR phone ILIKE ?", like, like)
	}
	if filter.Role != nil {
		q = q.Where("role = ?", string(*filter.Role))
	}
	if filter.IsActive != nil {
		q = q.Where("is_active = ?", *filter.IsActive)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var models []UserModel
	if err := q.Preload("City").Offset(pg.Offset()).Limit(pg.Limit()).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	users := make([]domain.User, len(models))
	for i := range models {
		users[i] = toDomain(&models[i])
	}
	return users, total, nil
}

func (r *GormRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&UserModel{}, "id = ?", id).Error
}

// GetCityByTitle resolves a city title to its ID.
func (r *GormRepository) GetCityByTitle(ctx context.Context, title string) (*uuid.UUID, error) {
	var city CityModel
	err := r.db.WithContext(ctx).Where("title = ?", title).First(&city).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &city.ID, nil
}

func toModel(u *domain.User) UserModel {
	m := UserModel{
		BaseModel: database.BaseModel{ID: u.ID, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt},
		FullName: u.FullName, Phone: u.Phone, Email: u.Email,
		PasswordHash: u.PasswordHash, Role: string(u.Role),
		NationalID: u.NationalID, Avatar: u.Avatar, Bio: u.Bio, CityID: u.CityID,
		Rating: u.Rating, RatingCount: u.RatingCount, CompletedTasks: u.CompletedTasks,
		IsVerified: u.IsVerified, IsActive: u.IsActive,
		PhoneVerified: u.PhoneVerified, EmailVerified: u.EmailVerified,
		NationalIDVerified: u.NationalIDVerified, VerificationLevel: u.VerificationLevel,
		VerificationStatus: u.VerificationStatus, VerificationReason: u.VerificationReason,
		VerifiedAt: u.VerifiedAt, WalletBalance: u.WalletBalance, LockedBalance: u.LockedBalance,
	}
	if u.DeletedAt != nil {
		m.DeletedAt = gorm.DeletedAt{Time: *u.DeletedAt, Valid: true}
	}
	return m
}

func toDomain(m *UserModel) domain.User {
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		deletedAt = &t
	}
	u := domain.User{
		ID: m.ID, FullName: m.FullName, Phone: m.Phone, Email: m.Email,
		PasswordHash: m.PasswordHash, Role: common.Role(m.Role),
		NationalID: m.NationalID, Avatar: m.Avatar, Bio: m.Bio, CityID: m.CityID,
		Rating: m.Rating, RatingCount: m.RatingCount, CompletedTasks: m.CompletedTasks,
		IsVerified: m.IsVerified, IsActive: m.IsActive,
		PhoneVerified: m.PhoneVerified, EmailVerified: m.EmailVerified,
		NationalIDVerified: m.NationalIDVerified, VerificationLevel: m.VerificationLevel,
		VerificationStatus: m.VerificationStatus, VerificationReason: m.VerificationReason,
		VerifiedAt: m.VerifiedAt, WalletBalance: m.WalletBalance, LockedBalance: m.LockedBalance,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt, DeletedAt: deletedAt,
	}
	if m.City != nil {
		u.CityTitle = m.City.Title
	}
	return u
}
