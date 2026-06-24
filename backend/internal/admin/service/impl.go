package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/admin/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	paymentdomain "github.com/younesbeheshti/any-task-connect/backend/internal/payment/domain"
	paymentrepo "github.com/younesbeheshti/any-task-connect/backend/internal/payment/repository"
	userdomain "github.com/younesbeheshti/any-task-connect/backend/internal/user/domain"
	userrepo "github.com/younesbeheshti/any-task-connect/backend/internal/user/repository"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
	"gorm.io/gorm"
)

type AdminService struct {
	userRepo    userrepo.Repository
	txRepo      paymentrepo.Repository
	db          *gorm.DB
}

func NewAdminService(userRepo userrepo.Repository, txRepo paymentrepo.Repository, db *gorm.DB) *AdminService {
	return &AdminService{
		userRepo: userRepo,
		txRepo:   txRepo,
		db:       db,
	}
}

func (s *AdminService) GetMetrics(ctx context.Context) (*domain.DashboardMetrics, error) {
	type counts struct {
		TotalUsers     int64
		ActiveUsers    int64
		TotalTasks     int64
		OpenTasks      int64
		CompletedTasks int64
	}
	var c counts
	s.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&c.TotalUsers)
	s.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM users WHERE is_active = true AND deleted_at IS NULL").Scan(&c.ActiveUsers)
	s.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM tasks WHERE deleted_at IS NULL").Scan(&c.TotalTasks)
	s.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM tasks WHERE status NOT IN ('CANCELLED','PAID','VERIFIED') AND deleted_at IS NULL").Scan(&c.OpenTasks)
	s.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM tasks WHERE status IN ('COMPLETED','VERIFIED','PAID') AND deleted_at IS NULL").Scan(&c.CompletedTasks)

	return &domain.DashboardMetrics{
		Users:   domain.MetricGroup{Total: c.TotalUsers, Active: c.ActiveUsers},
		Tasks:   domain.MetricGroup{Total: c.TotalTasks, Active: c.OpenTasks, Open: c.CompletedTasks},
		Revenue: domain.MetricGroup{},
	}, nil
}

func (s *AdminService) ListUsers(ctx context.Context, filter userrepo.UserFilter, pg common.PaginationParams) (*common.PaginatedResult[userdomain.User], error) {
	items, total, err := s.userRepo.List(ctx, filter, pg)
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
	return &common.PaginatedResult[userdomain.User]{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func (s *AdminService) UpdateUser(ctx context.Context, userID uuid.UUID, update domain.UserAdminUpdate) (*userdomain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err == common.ErrNotFound {
		return nil, apperrors.New("NOT_FOUND", "کاربر یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if update.IsActive != nil {
		user.IsActive = *update.IsActive
	}
	if update.Role != nil {
		r := common.RoleFromAPI(*update.Role)
		if r.IsValid() {
			user.Role = r
		}
	}
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AdminService) ListTransactions(ctx context.Context, filter paymentdomain.TransactionFilter, pg common.PaginationParams) (*common.PaginatedResult[paymentdomain.Transaction], error) {
	items, total, err := s.txRepo.ListAll(ctx, filter, pg)
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
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func (s *AdminService) GetRevenueReport(ctx context.Context, rangeKey string) (*domain.RevenueReport, error) {
	days := 30
	switch rangeKey {
	case "7d":
		days = 7
	case "90d":
		days = 90
	}
	type row struct {
		Date  string
		Value int64
	}
	var rows []row
	s.db.WithContext(ctx).Raw(`
		SELECT DATE(created_at)::text AS date, SUM(amount) AS value
		FROM transactions
		WHERE type = 'FEE' AND created_at >= NOW() - INTERVAL '? days'
		GROUP BY DATE(created_at)
		ORDER BY date
	`, days).Scan(&rows)

	points := make([]domain.DataPoint, len(rows))
	var total int64
	for i, r := range rows {
		points[i] = domain.DataPoint{Date: r.Date, Value: r.Value}
		total += r.Value
	}
	return &domain.RevenueReport{Range: rangeKey, Total: total, Points: points}, nil
}

func (s *AdminService) ResolveDispute(ctx context.Context, disputeID uuid.UUID, resolution domain.DisputeResolution) error {
	// Placeholder — disputes domain not yet implemented
	return apperrors.New("NOT_IMPLEMENTED", "مدیریت اعتراض در نسخه بعدی پیاده‌سازی خواهد شد", 501, apperrors.ErrInternal)
}

var _ Service = (*AdminService)(nil)

