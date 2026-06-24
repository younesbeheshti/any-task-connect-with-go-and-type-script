package infra

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/withdraw/domain"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, w *domain.WithdrawRequest) error {
	m := WithdrawModel{
		ID:          w.ID,
		UserID:      w.UserID,
		WalletID:    w.WalletID,
		Amount:      w.Amount,
		Status:      string(w.Status),
		BankAccount: w.BankAccount,
		ShebaNumber: w.ShebaNumber,
		Description: w.Description,
	}
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if m.Status == "" {
		m.Status = string(domain.StatusPending)
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	w.ID = m.ID
	w.CreatedAt = m.CreatedAt
	w.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.WithdrawRequest, error) {
	var m WithdrawModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m), nil
}

func (r *GormRepository) ListByUser(ctx context.Context, userID uuid.UUID, pg common.PaginationParams) ([]domain.WithdrawRequest, int64, error) {
	var models []WithdrawModel
	var count int64
	q := r.db.WithContext(ctx).Model(&WithdrawModel{}).Where("user_id = ?", userID)
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at desc").Offset(pg.Offset()).Limit(pg.Limit()).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	return toSlice(models), count, nil
}

func (r *GormRepository) ListAll(ctx context.Context, pg common.PaginationParams) ([]domain.WithdrawRequest, int64, error) {
	var models []WithdrawModel
	var count int64
	q := r.db.WithContext(ctx).Model(&WithdrawModel{})
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at desc").Offset(pg.Offset()).Limit(pg.Limit()).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	return toSlice(models), count, nil
}

func (r *GormRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.Status, approvedBy *uuid.UUID) error {
	updates := map[string]any{
		"status":     string(status),
		"updated_at": time.Now(),
	}
	if approvedBy != nil {
		now := time.Now()
		updates["approved_by"] = approvedBy
		updates["approved_at"] = now
	}
	return r.db.WithContext(ctx).Model(&WithdrawModel{}).Where("id = ?", id).Updates(updates).Error
}

func toDomain(m WithdrawModel) *domain.WithdrawRequest {
	return &domain.WithdrawRequest{
		ID:          m.ID,
		UserID:      m.UserID,
		WalletID:    m.WalletID,
		Amount:      m.Amount,
		Status:      domain.Status(m.Status),
		BankAccount: m.BankAccount,
		ShebaNumber: m.ShebaNumber,
		Description: m.Description,
		ApprovedBy:  m.ApprovedBy,
		ApprovedAt:  m.ApprovedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func toSlice(models []WithdrawModel) []domain.WithdrawRequest {
	out := make([]domain.WithdrawRequest, len(models))
	for i, m := range models {
		out[i] = *toDomain(m)
	}
	return out
}
