package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	walletrepo "github.com/younesbeheshti/any-task-connect/backend/internal/wallet/repository"
	"github.com/younesbeheshti/any-task-connect/backend/internal/withdraw/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/withdraw/repository"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
	"gorm.io/gorm"
)

type WithdrawService struct {
	repo       repository.Repository
	walletRepo walletrepo.Repository
	publisher  common.Publisher
	db         *gorm.DB
}

func NewWithdrawService(repo repository.Repository, walletRepo walletrepo.Repository, db *gorm.DB, publisher common.Publisher) *WithdrawService {
	return &WithdrawService{repo: repo, walletRepo: walletRepo, db: db, publisher: publisher}
}

func (s *WithdrawService) CreateWithdraw(ctx context.Context, input domain.CreateInput) (*domain.WithdrawRequest, error) {
	if input.Amount <= 0 {
		return nil, apperrors.New("INVALID_AMOUNT", "مبلغ درخواستی نامعتبر است", 422, apperrors.ErrValidation)
	}

	var result *domain.WithdrawRequest
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Resolve wallet by userID if not provided.
		wRepo := s.walletRepo.WithTx(tx)
		wallet, err := wRepo.GetByUserIDForUpdate(ctx, input.UserID)
		if err == common.ErrNotFound {
			return apperrors.New("NOT_FOUND", "کیف پول یافت نشد", 404, apperrors.ErrNotFound)
		}
		if err != nil {
			return err
		}
		if wallet.AvailableBalance < input.Amount {
			return apperrors.New("INSUFFICIENT_FUNDS", "موجودی کیف پول کافی نیست", 422, apperrors.ErrValidation)
		}

		// Move funds from available to pending_withdraw.
		wallet.AvailableBalance -= input.Amount
		wallet.PendingWithdrawBalance += input.Amount
		if err := wRepo.Update(ctx, wallet); err != nil {
			return err
		}

		req := &domain.WithdrawRequest{
			ID:          uuid.New(),
			UserID:      input.UserID,
			WalletID:    wallet.ID,
			Amount:      input.Amount,
			Status:      domain.StatusPending,
			BankAccount: input.BankAccount,
			ShebaNumber: input.ShebaNumber,
			Description: input.Description,
		}
		if err := s.repo.Create(ctx, req); err != nil {
			return err
		}
		result = req
		return nil
	})
	if err != nil {
		return nil, err
	}

	if s.publisher != nil {
		_ = s.publisher.Publish(ctx, common.EventWithdrawCreated, map[string]any{
			"withdrawId": result.ID.String(),
			"userId":     result.UserID.String(),
			"amount":     result.Amount,
		})
	}
	return result, nil
}

func (s *WithdrawService) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.WithdrawRequest, error) {
	req, err := s.repo.GetByID(ctx, id)
	if err == common.ErrNotFound {
		return nil, apperrors.New("NOT_FOUND", "درخواست برداشت یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if req.UserID != userID {
		return nil, apperrors.New("FORBIDDEN", "دسترسی مجاز نیست", 403, apperrors.ErrForbidden)
	}
	return req, nil
}

func (s *WithdrawService) ListUserWithdraws(ctx context.Context, userID uuid.UUID, pg common.PaginationParams) (*common.PaginatedResult[domain.WithdrawRequest], error) {
	items, total, err := s.repo.ListByUser(ctx, userID, pg)
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
	return &common.PaginatedResult[domain.WithdrawRequest]{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func (s *WithdrawService) Approve(ctx context.Context, id, adminID uuid.UUID) error {
	req, err := s.repo.GetByID(ctx, id)
	if err == common.ErrNotFound {
		return apperrors.New("NOT_FOUND", "درخواست برداشت یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return err
	}
	if req.Status != domain.StatusPending {
		return apperrors.New("INVALID_STATUS", "وضعیت درخواست برای تایید مجاز نیست", 422, apperrors.ErrConflict)
	}
	if err := s.repo.UpdateStatus(ctx, id, domain.StatusApproved, &adminID); err != nil {
		return err
	}
	if s.publisher != nil {
		_ = s.publisher.Publish(ctx, common.EventWithdrawApproved, map[string]any{
			"withdrawId": id.String(),
			"adminId":    adminID.String(),
		})
	}
	return nil
}

func (s *WithdrawService) Reject(ctx context.Context, id, adminID uuid.UUID) error {
	req, err := s.repo.GetByID(ctx, id)
	if err == common.ErrNotFound {
		return apperrors.New("NOT_FOUND", "درخواست برداشت یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return err
	}
	if req.Status != domain.StatusPending {
		return apperrors.New("INVALID_STATUS", "وضعیت درخواست برای رد مجاز نیست", 422, apperrors.ErrConflict)
	}
	if err := s.repo.UpdateStatus(ctx, id, domain.StatusRejected, &adminID); err != nil {
		return err
	}
	if s.publisher != nil {
		_ = s.publisher.Publish(ctx, common.EventWithdrawRejected, map[string]any{
			"withdrawId": id.String(),
			"adminId":    adminID.String(),
		})
	}
	return nil
}

func (s *WithdrawService) ListAll(ctx context.Context, pg common.PaginationParams) (*common.PaginatedResult[domain.WithdrawRequest], error) {
	items, total, err := s.repo.ListAll(ctx, pg)
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
	return &common.PaginatedResult[domain.WithdrawRequest]{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}
