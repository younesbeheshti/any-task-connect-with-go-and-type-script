package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appinfra "github.com/younesbeheshti/any-task-connect/backend/internal/application/infra"
	appservice "github.com/younesbeheshti/any-task-connect/backend/internal/application/service"
	"github.com/younesbeheshti/any-task-connect/backend/internal/application/domain"
	taskinfra "github.com/younesbeheshti/any-task-connect/backend/internal/task/infra"
	taskservice "github.com/younesbeheshti/any-task-connect/backend/internal/task/service"
	taskdomain "github.com/younesbeheshti/any-task-connect/backend/internal/task/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type noopPublisher struct{}

func (p *noopPublisher) Publish(_ context.Context, _ string, _ any) error { return nil }

func setupDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	sqls := []string{
		`CREATE TABLE categories (id TEXT PRIMARY KEY, title TEXT NOT NULL UNIQUE, icon TEXT DEFAULT '', created_at DATETIME)`,
		`CREATE TABLE cities (id TEXT PRIMARY KEY, title TEXT NOT NULL, province TEXT DEFAULT '', created_at DATETIME)`,
		`CREATE TABLE users (id TEXT PRIMARY KEY, full_name TEXT, phone TEXT, email TEXT, password_hash TEXT DEFAULT '', role TEXT DEFAULT 'REQUESTER', rating REAL DEFAULT 0, rating_count INTEGER DEFAULT 0, completed_tasks INTEGER DEFAULT 0, is_verified INTEGER DEFAULT 0, is_active INTEGER DEFAULT 1, phone_verified INTEGER DEFAULT 0, email_verified INTEGER DEFAULT 0, national_id_verified INTEGER DEFAULT 0, verification_level TEXT DEFAULT 'none', wallet_balance INTEGER DEFAULT 0, locked_balance INTEGER DEFAULT 0, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE tasks (id TEXT PRIMARY KEY, public_id TEXT NOT NULL UNIQUE, title TEXT NOT NULL, description TEXT NOT NULL, category_id TEXT NOT NULL, city_id TEXT NOT NULL, budget INTEGER NOT NULL, escrow_fee INTEGER NOT NULL DEFAULT 0, status TEXT NOT NULL DEFAULT 'CREATED', deadline DATETIME NOT NULL, requester_id TEXT NOT NULL, assigned_agent_id TEXT, attachment_urls TEXT NOT NULL DEFAULT '[]', view_count INTEGER NOT NULL DEFAULT 0, application_count INTEGER NOT NULL DEFAULT 0, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE task_timeline (id TEXT PRIMARY KEY, task_id TEXT NOT NULL, from_status TEXT, to_status TEXT NOT NULL, actor_id TEXT, note TEXT, created_at DATETIME)`,
		`CREATE TABLE applications (id TEXT PRIMARY KEY, task_id TEXT NOT NULL, agent_id TEXT NOT NULL, proposal_message TEXT NOT NULL DEFAULT '', expected_completion_time DATETIME, proposed_price INTEGER, eta TEXT NOT NULL DEFAULT '', status TEXT NOT NULL DEFAULT 'PENDING', created_at DATETIME, updated_at DATETIME)`,
	}
	for _, s := range sqls {
		require.NoError(t, db.Exec(s).Error)
	}

	catID, cityID := uuid.New(), uuid.New()
	require.NoError(t, db.Exec(`INSERT INTO categories (id, title, icon, created_at) VALUES (?, 'اداری', 'FileText', ?)`, catID, time.Now()).Error)
	require.NoError(t, db.Exec(`INSERT INTO cities (id, title, province, created_at) VALUES (?, 'تهران', 'تهران', ?)`, cityID, time.Now()).Error)
	db.Set("catID", catID)
	db.Set("cityID", cityID)
	return db
}

func setupServices(t *testing.T, db *gorm.DB) (*appservice.ApplicationService, *taskservice.TaskService, uuid.UUID, uuid.UUID, uuid.UUID) {
	t.Helper()
	var cat struct{ ID string }
	require.NoError(t, db.Raw("SELECT id FROM categories LIMIT 1").Scan(&cat).Error)
	catID, _ := uuid.Parse(cat.ID)
	var city struct{ ID string }
	require.NoError(t, db.Raw("SELECT id FROM cities LIMIT 1").Scan(&city).Error)
	cityID, _ := uuid.Parse(city.ID)

	requesterID := uuid.New()
	agentID := uuid.New()
	require.NoError(t, db.Exec(`INSERT INTO users (id, full_name, phone, password_hash, role, created_at, updated_at) VALUES (?, 'Requester', '09120000001', 'h', 'REQUESTER', ?, ?)`, requesterID, time.Now(), time.Now()).Error)
	require.NoError(t, db.Exec(`INSERT INTO users (id, full_name, phone, password_hash, role, created_at, updated_at) VALUES (?, 'Agent', '09120000002', 'h', 'AGENT', ?, ?)`, agentID, time.Now(), time.Now()).Error)

	taskRepo := taskinfra.NewGormRepository(db)
	appRepo := appinfra.NewGormRepository(db)
	taskSvc := taskservice.NewTaskService(taskRepo, &noopPublisher{})
	appSvc := appservice.NewApplicationService(appRepo, taskRepo)
	return appSvc, taskSvc, catID, cityID, requesterID
}

func createOpenTask(t *testing.T, taskSvc *taskservice.TaskService, catID, cityID, requesterID uuid.UUID) *taskdomain.Task {
	t.Helper()
	task, _, err := taskSvc.Create(context.Background(), taskdomain.CreateTaskInput{
		Title:       "تسک تست",
		Description: "این یک تسک برای تست پیشنهادها است",
		CategoryID:  catID,
		CityID:      cityID,
		Budget:      800000,
		Deadline:    time.Now().Add(7 * 24 * time.Hour),
		RequesterID: requesterID,
	})
	require.NoError(t, err)
	return task
}

func TestApplicationSubmit(t *testing.T) {
	db := setupDB(t)
	appSvc, taskSvc, catID, cityID, requesterID := setupServices(t, db)
	task := createOpenTask(t, taskSvc, catID, cityID, requesterID)

	agentID := uuid.New()
	require.NoError(t, db.Exec(`INSERT INTO users (id, full_name, phone, password_hash, role, created_at, updated_at) VALUES (?, 'Agent2', '09120000003', 'h', 'AGENT', ?, ?)`, agentID, time.Now(), time.Now()).Error)

	price := int64(700000)
	app, err := appSvc.Submit(context.Background(), domain.SubmitApplicationInput{
		TaskPublicID:    task.PublicID,
		AgentID:         agentID,
		ProposalMessage: "من می‌توانم این کار را انجام دهم",
		ProposedPrice:   &price,
		ETA:             "tomorrow",
	})

	require.NoError(t, err)
	assert.Equal(t, domain.ApplicationStatusPending, app.Status)
	assert.Equal(t, task.ID, app.TaskID)
}

func TestApplicationDuplicateReject(t *testing.T) {
	db := setupDB(t)
	appSvc, taskSvc, catID, cityID, requesterID := setupServices(t, db)
	task := createOpenTask(t, taskSvc, catID, cityID, requesterID)

	agentID := uuid.New()
	require.NoError(t, db.Exec(`INSERT INTO users (id, full_name, phone, password_hash, role, created_at, updated_at) VALUES (?, 'Agent3', '09120000004', 'h', 'AGENT', ?, ?)`, agentID, time.Now(), time.Now()).Error)

	input := domain.SubmitApplicationInput{
		TaskPublicID:    task.PublicID,
		AgentID:         agentID,
		ProposalMessage: "پیشنهاد اول",
	}
	_, err := appSvc.Submit(context.Background(), input)
	require.NoError(t, err)

	_, err = appSvc.Submit(context.Background(), input)
	assert.Error(t, err, "should fail on duplicate application")
}

func TestApplicationAccept(t *testing.T) {
	db := setupDB(t)
	appSvc, taskSvc, catID, cityID, requesterID := setupServices(t, db)
	task := createOpenTask(t, taskSvc, catID, cityID, requesterID)

	agentID := uuid.New()
	require.NoError(t, db.Exec(`INSERT INTO users (id, full_name, phone, password_hash, role, created_at, updated_at) VALUES (?, 'AcceptAgent', '09120000005', 'h', 'AGENT', ?, ?)`, agentID, time.Now(), time.Now()).Error)

	price := int64(750000)
	app, err := appSvc.Submit(context.Background(), domain.SubmitApplicationInput{
		TaskPublicID:    task.PublicID,
		AgentID:         agentID,
		ProposalMessage: "قبول می‌کنم",
		ProposedPrice:   &price,
	})
	require.NoError(t, err)

	accepted, err := appSvc.Accept(context.Background(), app.ID, requesterID)
	require.NoError(t, err)
	assert.Equal(t, domain.ApplicationStatusAccepted, accepted.Status)
}

func TestApplicationWithdraw(t *testing.T) {
	db := setupDB(t)
	appSvc, taskSvc, catID, cityID, requesterID := setupServices(t, db)
	task := createOpenTask(t, taskSvc, catID, cityID, requesterID)

	agentID := uuid.New()
	require.NoError(t, db.Exec(`INSERT INTO users (id, full_name, phone, password_hash, role, created_at, updated_at) VALUES (?, 'WithdrawAgent', '09120000006', 'h', 'AGENT', ?, ?)`, agentID, time.Now(), time.Now()).Error)

	app, err := appSvc.Submit(context.Background(), domain.SubmitApplicationInput{
		TaskPublicID:    task.PublicID,
		AgentID:         agentID,
		ProposalMessage: "پیشنهاد می‌دهم",
	})
	require.NoError(t, err)

	withdrawn, err := appSvc.Withdraw(context.Background(), app.ID, agentID)
	require.NoError(t, err)
	assert.Equal(t, domain.ApplicationStatusWithdrawn, withdrawn.Status)
}
