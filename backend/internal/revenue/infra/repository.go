package infra

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/revenue/domain"
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

func (r *GormRepository) Create(ctx context.Context, rev *domain.Revenue) error {
	m := RevenueModel{
		ID:     rev.ID,
		TaskID: rev.TaskID,
		Source: rev.Source,
		Amount: rev.Amount,
	}
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if m.Source == "" {
		m.Source = "commission"
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	rev.ID = m.ID
	rev.CreatedAt = m.CreatedAt
	return nil
}

func (r *GormRepository) GetTotal(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&RevenueModel{}).Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}

func (r *GormRepository) GetMonthToDate(ctx context.Context) (int64, error) {
	var total int64
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	err := r.db.WithContext(ctx).Model(&RevenueModel{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("created_at >= ?", start).
		Scan(&total).Error
	return total, err
}

func (r *GormRepository) GetDaily(ctx context.Context, days int) ([]domain.DailyRevenue, error) {
	type row struct {
		Date   string
		Amount int64
	}
	var rows []row
	since := time.Now().AddDate(0, 0, -days)
	err := r.db.WithContext(ctx).Model(&RevenueModel{}).
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') AS date, COALESCE(SUM(amount), 0) AS amount").
		Where("created_at >= ?", since).
		Group("TO_CHAR(created_at, 'YYYY-MM-DD')").
		Order("date asc").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make([]domain.DailyRevenue, len(rows))
	for i, row := range rows {
		result[i] = domain.DailyRevenue{Date: row.Date, Amount: row.Amount}
	}
	return result, nil
}

func (r *GormRepository) GetMonthly(ctx context.Context, months int) ([]domain.MonthlyRevenue, error) {
	type row struct {
		Month  string
		Amount int64
	}
	var rows []row
	since := time.Now().AddDate(0, -months, 0)
	err := r.db.WithContext(ctx).Model(&RevenueModel{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') AS month, COALESCE(SUM(amount), 0) AS amount").
		Where("created_at >= ?", since).
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Order("month asc").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make([]domain.MonthlyRevenue, len(rows))
	for i, row := range rows {
		result[i] = domain.MonthlyRevenue{Month: row.Month, Amount: row.Amount}
	}
	return result, nil
}

func (r *GormRepository) GetStats(ctx context.Context) (*domain.RevenueStats, error) {
	total, err := r.GetTotal(ctx)
	if err != nil {
		return nil, err
	}
	mtd, err := r.GetMonthToDate(ctx)
	if err != nil {
		return nil, err
	}
	daily, err := r.GetDaily(ctx, 30)
	if err != nil {
		return nil, err
	}
	monthly, err := r.GetMonthly(ctx, 12)
	if err != nil {
		return nil, err
	}
	return &domain.RevenueStats{
		TotalRevenue:   total,
		MonthToDate:    mtd,
		DailyRevenue:   daily,
		MonthlyRevenue: monthly,
	}, nil
}
