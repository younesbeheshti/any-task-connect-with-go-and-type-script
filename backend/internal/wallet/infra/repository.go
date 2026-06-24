package infra

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/wallet/domain"
	walletrepo "github.com/younesbeheshti/any-task-connect/backend/internal/wallet/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) WithTx(tx *gorm.DB) walletrepo.Repository {
	return &GormRepository{db: tx}
}

func (r *GormRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	m := WalletModel{
		ID:                     wallet.ID,
		UserID:                 wallet.UserID,
		AvailableBalance:       wallet.AvailableBalance,
		LockedBalance:          wallet.LockedBalance,
		PendingWithdrawBalance: wallet.PendingWithdrawBalance,
		TotalEarned:            wallet.TotalEarned,
		TotalSpent:             wallet.TotalSpent,
		Currency:               wallet.Currency,
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
	wallet.ID = m.ID
	wallet.CreatedAt = m.CreatedAt
	wallet.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *GormRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	var m WalletModel
	err := r.db.WithContext(ctx).First(&m, "user_id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m), nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	var m WalletModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m), nil
}

func (r *GormRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	var m WalletModel
	err := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m), nil
}

func (r *GormRepository) GetByUserIDForUpdate(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	var m WalletModel
	err := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&m, "user_id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m), nil
}

func (r *GormRepository) Update(ctx context.Context, wallet *domain.Wallet) error {
	wallet.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Model(&WalletModel{}).Where("id = ?", wallet.ID).
		Updates(map[string]any{
			"available_balance":        wallet.AvailableBalance,
			"locked_balance":           wallet.LockedBalance,
			"pending_withdraw_balance": wallet.PendingWithdrawBalance,
			"total_earned":             wallet.TotalEarned,
			"total_spent":              wallet.TotalSpent,
			"updated_at":               wallet.UpdatedAt,
		}).Error
}

func (r *GormRepository) LockFunds(ctx context.Context, walletID uuid.UUID, amount int64) error {
	return r.db.WithContext(ctx).Model(&WalletModel{}).Where("id = ?", walletID).
		Updates(map[string]any{
			"available_balance": gorm.Expr("available_balance - ?", amount),
			"locked_balance":    gorm.Expr("locked_balance + ?", amount),
			"updated_at":        time.Now(),
		}).Error
}

func (r *GormRepository) UnlockFunds(ctx context.Context, walletID uuid.UUID, amount int64) error {
	return r.db.WithContext(ctx).Model(&WalletModel{}).Where("id = ?", walletID).
		Updates(map[string]any{
			"available_balance": gorm.Expr("available_balance + ?", amount),
			"locked_balance":    gorm.Expr("locked_balance - ?", amount),
			"updated_at":        time.Now(),
		}).Error
}

func (r *GormRepository) ReleaseToAgent(ctx context.Context, requesterWalletID, agentWalletID uuid.UUID, amount int64) error {
	if err := r.db.WithContext(ctx).Model(&WalletModel{}).Where("id = ?", requesterWalletID).
		Updates(map[string]any{
			"locked_balance": gorm.Expr("locked_balance - ?", amount),
			"total_spent":    gorm.Expr("total_spent + ?", amount),
			"updated_at":     time.Now(),
		}).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&WalletModel{}).Where("id = ?", agentWalletID).
		Updates(map[string]any{
			"available_balance": gorm.Expr("available_balance + ?", amount),
			"total_earned":      gorm.Expr("total_earned + ?", amount),
			"updated_at":        time.Now(),
		}).Error
}

func toDomain(m WalletModel) *domain.Wallet {
	return &domain.Wallet{
		ID:                     m.ID,
		UserID:                 m.UserID,
		AvailableBalance:       m.AvailableBalance,
		LockedBalance:          m.LockedBalance,
		PendingWithdrawBalance: m.PendingWithdrawBalance,
		TotalEarned:            m.TotalEarned,
		TotalSpent:             m.TotalSpent,
		Currency:               m.Currency,
		CreatedAt:              m.CreatedAt,
		UpdatedAt:              m.UpdatedAt,
	}
}
