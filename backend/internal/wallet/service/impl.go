package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	ledgerdomain "github.com/younesbeheshti/any-task-connect/backend/internal/ledger/domain"
	ledgerservice "github.com/younesbeheshti/any-task-connect/backend/internal/ledger/service"
	paymentdomain "github.com/younesbeheshti/any-task-connect/backend/internal/payment/domain"
	paymentrepo "github.com/younesbeheshti/any-task-connect/backend/internal/payment/repository"
	revenueservice "github.com/younesbeheshti/any-task-connect/backend/internal/revenue/service"
	walletdomain "github.com/younesbeheshti/any-task-connect/backend/internal/wallet/domain"
	walletrepo "github.com/younesbeheshti/any-task-connect/backend/internal/wallet/repository"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
	"gorm.io/gorm"
)

type WalletService struct {
	walletRepo walletrepo.Repository
	txRepo     paymentrepo.Repository
	ledger     ledgerservice.Service
	revenue    revenueservice.Service
	publisher  common.Publisher
	cache      common.Cache
	db         *gorm.DB
	commission int64
}

func NewWalletService(
	walletRepo walletrepo.Repository,
	txRepo paymentrepo.Repository,
	ledger ledgerservice.Service,
	revenue revenueservice.Service,
	publisher common.Publisher,
	cache common.Cache,
	db *gorm.DB,
	commission int64,
) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
		txRepo:     txRepo,
		ledger:     ledger,
		revenue:    revenue,
		publisher:  publisher,
		cache:      cache,
		db:         db,
		commission: commission,
	}
}

func (s *WalletService) EnsureWallet(ctx context.Context, userID uuid.UUID) (*walletdomain.Wallet, error) {
	existing, err := s.walletRepo.GetByUserID(ctx, userID)
	if err == nil {
		return existing, nil
	}
	if err != common.ErrNotFound {
		return nil, err
	}
	wallet := &walletdomain.Wallet{
		ID:       uuid.New(),
		UserID:   userID,
		Currency: "IRT",
	}
	if err := s.walletRepo.Create(ctx, wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *WalletService) GetWallet(ctx context.Context, userID uuid.UUID) (*walletdomain.Wallet, error) {
	if s.cache != nil {
		var cached walletdomain.Wallet
		if err := s.cache.Get(ctx, common.KeyWalletCache+userID.String(), &cached); err == nil {
			return &cached, nil
		}
	}
	wallet, err := s.EnsureWallet(ctx, userID)
	if err != nil {
		return nil, err
	}
	if s.cache != nil {
		_ = s.cache.Set(ctx, common.KeyWalletCache+userID.String(), wallet, 0)
	}
	return wallet, nil
}

func (s *WalletService) TopUp(ctx context.Context, userID uuid.UUID, amount int64, cardID *uuid.UUID, returnURL string) (string, error) {
	if amount <= 0 {
		return "", apperrors.New("INVALID_AMOUNT", "مبلغ شارژ نامعتبر است", 422, apperrors.ErrValidation)
	}
	_, err := s.EnsureWallet(ctx, userID)
	if err != nil {
		return "", err
	}
	// Placeholder: real payment gateway integration would go here.
	return "https://payment.example.com/pay/" + uuid.New().String(), nil
}

func (s *WalletService) Withdraw(ctx context.Context, userID uuid.UUID, amount int64, iban string) error {
	if amount <= 0 {
		return apperrors.New("INVALID_AMOUNT", "مبلغ برداشت نامعتبر است", 422, apperrors.ErrValidation)
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.walletRepo.WithTx(tx)
		wallet, err := repo.GetByUserIDForUpdate(ctx, userID)
		if err == common.ErrNotFound {
			return apperrors.New("NOT_FOUND", "کیف پول یافت نشد", 404, apperrors.ErrNotFound)
		}
		if err != nil {
			return err
		}
		if !wallet.CanWithdraw(amount) {
			return apperrors.New("INSUFFICIENT_FUNDS", "موجودی کیف پول کافی نیست", 422, apperrors.ErrValidation)
		}
		wallet.AvailableBalance -= amount
		wallet.PendingWithdrawBalance += amount
		if err := repo.Update(ctx, wallet); err != nil {
			return err
		}
		s.invalidateCache(ctx, userID)
		return nil
	})
}

func (s *WalletService) GetTransactions(ctx context.Context, userID uuid.UUID, filter paymentdomain.TransactionFilter, pg common.PaginationParams) (*common.PaginatedResult[paymentdomain.Transaction], error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err == common.ErrNotFound {
		return &common.PaginatedResult[paymentdomain.Transaction]{Items: []paymentdomain.Transaction{}, Page: 1, PageSize: 20}, nil
	}
	if err != nil {
		return nil, err
	}
	filter.WalletID = &wallet.ID
	txs, total, err := s.txRepo.List(ctx, filter, pg)
	if err != nil {
		return nil, err
	}
	page := pg.Page
	if page < 1 {
		page = 1
	}
	pageSize := pg.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	return &common.PaginatedResult[paymentdomain.Transaction]{
		Items:    txs,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func (s *WalletService) GetStatistics(ctx context.Context, userID uuid.UUID) (*WalletStats, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err == common.ErrNotFound {
		return &WalletStats{}, nil
	}
	if err != nil {
		return nil, err
	}
	filter := paymentdomain.TransactionFilter{WalletID: &wallet.ID}
	_, count, err := s.txRepo.List(ctx, filter, common.PaginationParams{Page: 1, PageSize: 1})
	if err != nil {
		count = 0
	}
	return &WalletStats{
		TotalEarned: wallet.TotalEarned,
		TotalSpent:  wallet.TotalSpent,
		TxCount:     count,
	}, nil
}

func (s *WalletService) LockEscrow(ctx context.Context, input paymentdomain.EscrowLockInput) error {
	total := input.Budget + input.Fee
	return s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.walletRepo.WithTx(tx)
		txRepo := s.txRepo.WithTx(tx)

		wallet, err := repo.GetByUserIDForUpdate(ctx, input.RequesterID)
		if err == common.ErrNotFound {
			return apperrors.New("INSUFFICIENT_FUNDS", "موجودی کیف پول کافی نیست", 422, apperrors.ErrValidation)
		}
		if err != nil {
			return err
		}
		if !wallet.CanLock(total) {
			return apperrors.New("INSUFFICIENT_FUNDS", "موجودی کیف پول کافی نیست", 422, apperrors.ErrValidation)
		}

		wallet.AvailableBalance -= total
		wallet.LockedBalance += total
		if err := repo.Update(ctx, wallet); err != nil {
			return err
		}

		lockTx := &paymentdomain.Transaction{
			ID:              uuid.New(),
			WalletID:        wallet.ID,
			TaskID:          &input.TaskID,
			Amount:          total,
			Currency:        "IRT",
			Type:            paymentdomain.TransactionTypeEscrowLock,
			Status:          paymentdomain.TransactionStatusLocked,
			ReferenceNumber: fmt.Sprintf("TXN-%s", uuid.New().String()[:8]),
			Description:     "بلوکه شدن وجه برای تسک",
		}
		if err := txRepo.Create(ctx, lockTx); err != nil {
			return err
		}

		if s.ledger != nil {
			entries := []ledgerdomain.Entry{
				{
					ID:            uuid.New(),
					TransactionID: lockTx.ID,
					Account:       ledgerdomain.AccountUserWallet,
					Debit:         total,
					Credit:        0,
					BalanceAfter:  wallet.AvailableBalance,
				},
				{
					ID:            uuid.New(),
					TransactionID: lockTx.ID,
					Account:       ledgerdomain.AccountEscrow,
					Debit:         0,
					Credit:        total,
					BalanceAfter:  total,
				},
			}
			_ = s.ledger.Record(ctx, entries)
		}

		s.invalidateCache(ctx, input.RequesterID)

		if s.publisher != nil {
			_ = s.publisher.Publish(ctx, common.EventEscrowLocked, map[string]any{
				"taskId":      input.TaskID.String(),
				"requesterId": input.RequesterID.String(),
				"amount":      total,
			})
		}
		return nil
	})
}

func (s *WalletService) ReleaseEscrow(ctx context.Context, input paymentdomain.EscrowReleaseInput) error {
	fee := input.Amount * s.commission / 100
	agentAmount := input.Amount - fee

	return s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.walletRepo.WithTx(tx)
		txRepo := s.txRepo.WithTx(tx)

		requesterWallet, err := repo.GetByUserIDForUpdate(ctx, input.RequesterID)
		if err != nil {
			return err
		}
		agentWallet, err := repo.GetByUserIDForUpdate(ctx, input.AgentID)
		if err == common.ErrNotFound {
			agentWallet, err = s.createWalletInTx(ctx, repo, input.AgentID)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}

		totalLocked := input.Amount + input.Fee
		requesterWallet.LockedBalance -= totalLocked
		requesterWallet.TotalSpent += input.Amount
		if err := repo.Update(ctx, requesterWallet); err != nil {
			return err
		}

		agentWallet.AvailableBalance += agentAmount
		agentWallet.TotalEarned += agentAmount
		if err := repo.Update(ctx, agentWallet); err != nil {
			return err
		}

		releaseTx := &paymentdomain.Transaction{
			ID:              uuid.New(),
			WalletID:        requesterWallet.ID,
			TaskID:          &input.TaskID,
			Amount:          totalLocked,
			Currency:        "IRT",
			Type:            paymentdomain.TransactionTypeEscrowRelease,
			Status:          paymentdomain.TransactionStatusReleased,
			ReferenceNumber: fmt.Sprintf("TXN-%s", uuid.New().String()[:8]),
			Description:     "آزادسازی وجه امانت",
		}
		if err := txRepo.Create(ctx, releaseTx); err != nil {
			return err
		}

		paymentTx := &paymentdomain.Transaction{
			ID:              uuid.New(),
			WalletID:        agentWallet.ID,
			TaskID:          &input.TaskID,
			Amount:          agentAmount,
			Currency:        "IRT",
			Type:            paymentdomain.TransactionTypePayment,
			Status:          paymentdomain.TransactionStatusCompleted,
			ReferenceNumber: fmt.Sprintf("TXN-%s", uuid.New().String()[:8]),
			Description:     "پرداخت به مجری",
		}
		if err := txRepo.Create(ctx, paymentTx); err != nil {
			return err
		}

		if fee > 0 {
			commTx := &paymentdomain.Transaction{
				ID:              uuid.New(),
				WalletID:        requesterWallet.ID,
				TaskID:          &input.TaskID,
				Amount:          fee,
				Currency:        "IRT",
				Type:            paymentdomain.TransactionTypeCommission,
				Status:          paymentdomain.TransactionStatusCompleted,
				ReferenceNumber: fmt.Sprintf("TXN-%s", uuid.New().String()[:8]),
				Description:     "کارمزد پلتفرم",
			}
			if err := txRepo.Create(ctx, commTx); err != nil {
				return err
			}

			if s.revenue != nil {
				_ = s.revenue.RecordCommission(ctx, &input.TaskID, fee)
			}
		}

		s.invalidateCache(ctx, input.RequesterID)
		s.invalidateCache(ctx, input.AgentID)

		if s.publisher != nil {
			_ = s.publisher.Publish(ctx, common.EventEscrowReleased, map[string]any{
				"taskId":  input.TaskID.String(),
				"agentId": input.AgentID.String(),
				"amount":  agentAmount,
			})
		}
		return nil
	})
}

func (s *WalletService) RefundEscrow(ctx context.Context, taskID, requesterID uuid.UUID, amount, fee int64) error {
	total := amount + fee
	return s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.walletRepo.WithTx(tx)
		txRepo := s.txRepo.WithTx(tx)

		wallet, err := repo.GetByUserIDForUpdate(ctx, requesterID)
		if err != nil {
			return err
		}

		wallet.LockedBalance -= total
		wallet.AvailableBalance += total
		if err := repo.Update(ctx, wallet); err != nil {
			return err
		}

		refundTx := &paymentdomain.Transaction{
			ID:              uuid.New(),
			WalletID:        wallet.ID,
			TaskID:          &taskID,
			Amount:          total,
			Currency:        "IRT",
			Type:            paymentdomain.TransactionTypeRefund,
			Status:          paymentdomain.TransactionStatusCompleted,
			ReferenceNumber: fmt.Sprintf("TXN-%s", uuid.New().String()[:8]),
			Description:     "استرداد وجه پس از لغو تسک",
		}
		if err := txRepo.Create(ctx, refundTx); err != nil {
			return err
		}

		s.invalidateCache(ctx, requesterID)
		return nil
	})
}

func (s *WalletService) createWalletInTx(ctx context.Context, repo walletrepo.Repository, userID uuid.UUID) (*walletdomain.Wallet, error) {
	wallet := &walletdomain.Wallet{
		ID:       uuid.New(),
		UserID:   userID,
		Currency: "IRT",
	}
	if err := repo.Create(ctx, wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *WalletService) invalidateCache(ctx context.Context, userID uuid.UUID) {
	if s.cache != nil {
		_ = s.cache.Delete(ctx, common.KeyWalletCache+userID.String())
	}
}
