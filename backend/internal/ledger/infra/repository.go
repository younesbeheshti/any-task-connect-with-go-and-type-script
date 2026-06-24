package infra

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/ledger/domain"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) WithTx(tx *gorm.DB) *GormRepository {
	return &GormRepository{db: tx}
}

func (r *GormRepository) Create(ctx context.Context, entries []domain.Entry) error {
	if len(entries) == 0 {
		return nil
	}
	models := make([]LedgerEntryModel, len(entries))
	for i, e := range entries {
		models[i] = LedgerEntryModel{
			ID:            e.ID,
			TransactionID: e.TransactionID,
			Account:       string(e.Account),
			Debit:         e.Debit,
			Credit:        e.Credit,
			BalanceAfter:  e.BalanceAfter,
		}
		if models[i].ID == uuid.Nil {
			models[i].ID = uuid.New()
		}
	}
	return r.db.WithContext(ctx).Create(&models).Error
}

func (r *GormRepository) GetByTransactionID(ctx context.Context, txID uuid.UUID) ([]domain.Entry, error) {
	var models []LedgerEntryModel
	if err := r.db.WithContext(ctx).Where("transaction_id = ?", txID).Order("created_at asc").Find(&models).Error; err != nil {
		return nil, err
	}
	entries := make([]domain.Entry, len(models))
	for i, m := range models {
		entries[i] = domain.Entry{
			ID:            m.ID,
			TransactionID: m.TransactionID,
			Account:       domain.Account(m.Account),
			Debit:         m.Debit,
			Credit:        m.Credit,
			BalanceAfter:  m.BalanceAfter,
			CreatedAt:     m.CreatedAt,
		}
	}
	return entries, nil
}
