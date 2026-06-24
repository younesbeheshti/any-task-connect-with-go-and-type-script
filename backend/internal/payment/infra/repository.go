package infra

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/payment/domain"
	paymentrepo "github.com/younesbeheshti/any-task-connect/backend/internal/payment/repository"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) WithTx(tx *gorm.DB) paymentrepo.Repository {
	return &GormRepository{db: tx}
}

func (r *GormRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	m := TransactionModel{
		ID:              tx.ID,
		WalletID:        tx.WalletID,
		TaskID:          tx.TaskID,
		Amount:          tx.Amount,
		Currency:        tx.Currency,
		Type:            string(tx.Type),
		Status:          string(tx.Status),
		ReferenceNumber: tx.ReferenceNumber,
		Description:     tx.Description,
	}
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if m.Currency == "" {
		m.Currency = "IRT"
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	tx.ID = m.ID
	tx.CreatedAt = m.CreatedAt
	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	var m TransactionModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m), nil
}

func (r *GormRepository) GetByReference(ctx context.Context, ref string) (*domain.Transaction, error) {
	var m TransactionModel
	err := r.db.WithContext(ctx).First(&m, "reference_number = ?", ref).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m), nil
}

func (r *GormRepository) List(ctx context.Context, filter domain.TransactionFilter, pg common.PaginationParams) ([]domain.Transaction, int64, error) {
	var models []TransactionModel
	var count int64
	q := r.db.WithContext(ctx).Model(&TransactionModel{})
	q = applyFilter(q, filter)
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at desc").Offset(pg.Offset()).Limit(pg.Limit()).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	return toSlice(models), count, nil
}

func (r *GormRepository) ListAll(ctx context.Context, filter domain.TransactionFilter, pg common.PaginationParams) ([]domain.Transaction, int64, error) {
	return r.List(ctx, filter, pg)
}

func applyFilter(q *gorm.DB, filter domain.TransactionFilter) *gorm.DB {
	if filter.WalletID != nil {
		q = q.Where("wallet_id = ?", *filter.WalletID)
	}
	if filter.TaskID != nil {
		q = q.Where("task_id = ?", *filter.TaskID)
	}
	if filter.Type != nil {
		q = q.Where("type = ?", string(*filter.Type))
	}
	if filter.DateFrom != nil {
		q = q.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		q = q.Where("created_at <= ?", *filter.DateTo)
	}
	return q
}

func toDomain(m TransactionModel) *domain.Transaction {
	return &domain.Transaction{
		ID:              m.ID,
		WalletID:        m.WalletID,
		TaskID:          m.TaskID,
		Amount:          m.Amount,
		Currency:        m.Currency,
		Type:            domain.TransactionType(m.Type),
		Status:          domain.TransactionStatus(m.Status),
		ReferenceNumber: m.ReferenceNumber,
		Description:     m.Description,
		CreatedAt:       m.CreatedAt,
	}
}

func toSlice(models []TransactionModel) []domain.Transaction {
	out := make([]domain.Transaction, len(models))
	for i, m := range models {
		out[i] = *toDomain(m)
	}
	return out
}
