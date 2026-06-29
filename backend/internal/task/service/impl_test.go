package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/task/domain"
	taskinfra "github.com/younesbeheshti/any-task-connect/backend/internal/task/infra"
	taskservice "github.com/younesbeheshti/any-task-connect/backend/internal/task/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// noopPublisher implements common.Publisher for tests.
type noopPublisher struct{}

func (p *noopPublisher) Publish(_ context.Context, _ string, _ any) error { return nil }

func setupTaskTest(t *testing.T) (*taskservice.TaskService, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	require.NoError(t, db.Exec(`CREATE TABLE categories (
		id TEXT PRIMARY KEY, title TEXT NOT NULL UNIQUE, icon TEXT NOT NULL DEFAULT '',
		created_at DATETIME
	)`).Error)
	require.NoError(t, db.Exec(`CREATE TABLE cities (
		id TEXT PRIMARY KEY, title TEXT NOT NULL, province TEXT NOT NULL DEFAULT '',
		created_at DATETIME
	)`).Error)
	require.NoError(t, db.Exec(`CREATE TABLE users (
		id TEXT PRIMARY KEY, full_name TEXT, phone TEXT, email TEXT,
		password_hash TEXT DEFAULT '', role TEXT DEFAULT 'REQUESTER',
		rating REAL DEFAULT 0, rating_count INTEGER DEFAULT 0,
		completed_tasks INTEGER DEFAULT 0, is_verified INTEGER DEFAULT 0,
		is_active INTEGER DEFAULT 1, phone_verified INTEGER DEFAULT 0,
		email_verified INTEGER DEFAULT 0, national_id_verified INTEGER DEFAULT 0,
		verification_level TEXT DEFAULT 'none', wallet_balance INTEGER DEFAULT 0,
		locked_balance INTEGER DEFAULT 0,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME
	)`).Error)
	require.NoError(t, db.Exec(`CREATE TABLE tasks (
		id TEXT PRIMARY KEY, public_id TEXT NOT NULL UNIQUE,
		title TEXT NOT NULL, description TEXT NOT NULL,
		category_id TEXT NOT NULL, city_id TEXT NOT NULL,
		budget INTEGER NOT NULL, escrow_fee INTEGER NOT NULL DEFAULT 0,
		status TEXT NOT NULL DEFAULT 'CREATED',
		deadline DATETIME NOT NULL, requester_id TEXT NOT NULL,
		assigned_agent_id TEXT, attachment_urls TEXT NOT NULL DEFAULT '[]',
		view_count INTEGER NOT NULL DEFAULT 0,
		application_count INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME
	)`).Error)
	require.NoError(t, db.Exec(`CREATE TABLE task_timeline (
		id TEXT PRIMARY KEY, task_id TEXT NOT NULL,
		from_status TEXT, to_status TEXT NOT NULL,
		actor_id TEXT, note TEXT, created_at DATETIME
	)`).Error)

	// Seed a category and city.
	catID := uuid.New()
	cityID := uuid.New()
	require.NoError(t, db.Exec(`INSERT INTO categories (id, title, icon, created_at) VALUES (?, 'اداری', 'FileText', ?)`,
		catID.String(), time.Now()).Error)
	require.NoError(t, db.Exec(`INSERT INTO cities (id, title, province, created_at) VALUES (?, 'تهران', 'تهران', ?)`,
		cityID.String(), time.Now()).Error)

	repo := taskinfra.NewGormRepository(db)
	svc := taskservice.NewTaskService(repo, &noopPublisher{})
	return svc, db
}

func seedIDs(t *testing.T, db *gorm.DB) (categoryID, cityID, requesterID uuid.UUID) {
	t.Helper()
	var cat struct{ ID string }
	require.NoError(t, db.Raw("SELECT id FROM categories LIMIT 1").Scan(&cat).Error)
	categoryID, _ = uuid.Parse(cat.ID)

	var city struct{ ID string }
	require.NoError(t, db.Raw("SELECT id FROM cities LIMIT 1").Scan(&city).Error)
	cityID, _ = uuid.Parse(city.ID)

	requesterID = uuid.New()
	require.NoError(t, db.Exec(`INSERT INTO users (id, full_name, phone, password_hash, role, created_at, updated_at) VALUES (?, 'Test User', '09120000000', 'hash', 'REQUESTER', ?, ?)`,
		requesterID.String(), time.Now(), time.Now()).Error)
	return
}

func TestTaskCreate(t *testing.T) {
	svc, db := setupTaskTest(t)
	catID, cityID, reqID := seedIDs(t, db)

	task, escrow, err := svc.Create(context.Background(), domain.CreateTaskInput{
		Title:       "دریافت ریزنمرات",
		Description: "از دانشگاه تهران ریزنمرات بگیرید",
		CategoryID:  catID,
		CityID:      cityID,
		Budget:      750000,
		Deadline:    time.Now().Add(7 * 24 * time.Hour),
		RequesterID: reqID,
	})

	require.NoError(t, err)
	assert.Equal(t, "TB-1", task.PublicID)
	assert.Equal(t, domain.TaskStatusOpen, task.Status)
	assert.Equal(t, int64(60000), escrow.Fee) // 8% of 750000
	assert.Equal(t, "IRT", escrow.Currency)
}

func TestTaskStateMachine_ValidTransitions(t *testing.T) {
	task := &domain.Task{Status: domain.TaskStatusOpen}
	assert.True(t, task.CanTransitionTo(domain.TaskStatusAssigned))
	assert.True(t, task.CanTransitionTo(domain.TaskStatusCancelled))
	assert.False(t, task.CanTransitionTo(domain.TaskStatusPaid))
	assert.False(t, task.CanTransitionTo(domain.TaskStatusInProgress))
}

func TestTaskStateMachine_InvalidTransition(t *testing.T) {
	task := &domain.Task{Status: domain.TaskStatusPaid}
	assert.False(t, task.CanTransitionTo(domain.TaskStatusOpen))
	assert.False(t, task.CanTransitionTo(domain.TaskStatusCancelled))
}

func TestTaskCancel(t *testing.T) {
	svc, db := setupTaskTest(t)
	catID, cityID, reqID := seedIDs(t, db)

	created, _, err := svc.Create(context.Background(), domain.CreateTaskInput{
		Title:       "کار لغوشدنی",
		Description: "این تسک لغو خواهد شد برای تست",
		CategoryID:  catID,
		CityID:      cityID,
		Budget:      500000,
		Deadline:    time.Now().Add(3 * 24 * time.Hour),
		RequesterID: reqID,
	})
	require.NoError(t, err)
	assert.Equal(t, domain.TaskStatusOpen, created.Status)

	cancelled, err := svc.Cancel(context.Background(), created.ID, reqID)
	require.NoError(t, err)
	assert.Equal(t, domain.TaskStatusCancelled, cancelled.Status)
}

func TestTaskList(t *testing.T) {
	svc, db := setupTaskTest(t)
	catID, cityID, reqID := seedIDs(t, db)

	for i := 0; i < 3; i++ {
		_, _, err := svc.Create(context.Background(), domain.CreateTaskInput{
			Title:       "تسک آزمایشی",
			Description: "این یک تسک آزمایشی است برای تست لیست",
			CategoryID:  catID,
			CityID:      cityID,
			Budget:      int64(500000 + i*100000),
			Deadline:    time.Now().Add(7 * 24 * time.Hour),
			RequesterID: reqID,
		})
		require.NoError(t, err)
	}

	result, err := svc.List(context.Background(), domain.TaskFilter{}, common.PaginationParams{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(3), result.Total)
	assert.Len(t, result.Items, 3)
}

func TestTaskAPIStatus(t *testing.T) {
	cases := []struct {
		status domain.TaskStatus
		api    string
	}{
		{domain.TaskStatusCreated, "posted"},
		{domain.TaskStatusOpen, "awaiting_applicants"},
		{domain.TaskStatusAssigned, "accepted"},
		{domain.TaskStatusInProgress, "in_progress"},
		{domain.TaskStatusCompleted, "completed"},
		{domain.TaskStatusWaitingForVerification, "awaiting_verification"},
		{domain.TaskStatusVerified, "verified"},
		{domain.TaskStatusPaid, "paid"},
		{domain.TaskStatusCancelled, "cancelled"},
	}
	for _, c := range cases {
		assert.Equal(t, c.api, c.status.APIStatus(), "status %s", c.status)
	}
}
