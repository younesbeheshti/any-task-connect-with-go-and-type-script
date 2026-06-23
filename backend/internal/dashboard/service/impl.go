package service

import (
	"context"

	"gorm.io/gorm"
)

// DashboardService implements Service using direct DB queries.
type DashboardService struct {
	db *gorm.DB
}

// NewDashboardService creates a new DashboardService.
func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{db: db}
}

func (s *DashboardService) GetAdminStats(ctx context.Context) (*Stats, error) {
	var stats Stats

	if err := s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`).Scan(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM tasks WHERE deleted_at IS NULL`).Scan(&stats.TotalTasks).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM tasks WHERE status = 'OPEN' AND deleted_at IS NULL`).Scan(&stats.OpenTasks).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM tasks WHERE status IN ('COMPLETED','WAITING_FOR_VERIFICATION','VERIFIED','PAID') AND deleted_at IS NULL`).Scan(&stats.CompletedTasks).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM applications`).Scan(&stats.TotalApplications).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func (s *DashboardService) GetUserStats(ctx context.Context, userID string) (*UserStats, error) {
	var stats UserStats

	if err := s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM tasks WHERE requester_id = ? AND deleted_at IS NULL`, userID).Scan(&stats.PostedTasks).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM tasks WHERE requester_id = ? AND status IN ('COMPLETED','WAITING_FOR_VERIFICATION','VERIFIED','PAID') AND deleted_at IS NULL`, userID).Scan(&stats.CompletedTasks).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM tasks WHERE requester_id = ? AND status IN ('OPEN','ASSIGNED','IN_PROGRESS') AND deleted_at IS NULL`, userID).Scan(&stats.ActiveTasks).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM applications WHERE agent_id = ?`, userID).Scan(&stats.TotalApplications).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
