package service

import "context"

// Stats holds platform-level statistics for the admin dashboard.
type Stats struct {
	TotalUsers        int64 `json:"totalUsers"`
	TotalTasks        int64 `json:"totalTasks"`
	OpenTasks         int64 `json:"openTasks"`
	CompletedTasks    int64 `json:"completedTasks"`
	TotalApplications int64 `json:"totalApplications"`
}

// UserStats holds per-user task and application statistics.
type UserStats struct {
	PostedTasks      int64 `json:"postedTasks"`
	CompletedTasks   int64 `json:"completedTasks"`
	ActiveTasks      int64 `json:"activeTasks"`
	TotalApplications int64 `json:"totalApplications"`
}

// Service defines dashboard statistics use cases.
type Service interface {
	GetAdminStats(ctx context.Context) (*Stats, error)
	GetUserStats(ctx context.Context, userID string) (*UserStats, error)
}
