package infra

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/application/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"gorm.io/gorm"
)

// GormRepository implements application/repository.Repository using GORM.
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new application GormRepository.
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, app *domain.Application) error {
	m := ApplicationModel{
		ID:                     app.ID,
		TaskID:                 app.TaskID,
		AgentID:                app.AgentID,
		ProposalMessage:        app.ProposalMessage,
		ExpectedCompletionTime: app.ExpectedCompletionTime,
		ProposedPrice:          app.ProposedPrice,
		ETA:                    app.ETA,
		Status:                 string(app.Status),
	}
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	app.ID = m.ID
	app.CreatedAt = m.CreatedAt
	app.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Application, error) {
	var m ApplicationModel
	err := r.db.WithContext(ctx).Preload("Agent").First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m), nil
}

func (r *GormRepository) GetByTaskAndAgent(ctx context.Context, taskID, agentID uuid.UUID) (*domain.Application, error) {
	var m ApplicationModel
	err := r.db.WithContext(ctx).
		Where("task_id = ? AND agent_id = ?", taskID, agentID).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m), nil
}

func (r *GormRepository) ListByTask(ctx context.Context, taskID uuid.UUID) ([]domain.Application, error) {
	var models []ApplicationModel
	if err := r.db.WithContext(ctx).Preload("Agent").Where("task_id = ?", taskID).Order("created_at asc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toSlice(models), nil
}

func (r *GormRepository) ListByAgent(ctx context.Context, agentID uuid.UUID) ([]domain.Application, error) {
	var models []ApplicationModel
	if err := r.db.WithContext(ctx).Where("agent_id = ?", agentID).Order("created_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toSlice(models), nil
}

func (r *GormRepository) Update(ctx context.Context, app *domain.Application) error {
	return r.db.WithContext(ctx).Model(&ApplicationModel{}).Where("id = ?", app.ID).
		Updates(map[string]any{
			"status":     string(app.Status),
			"updated_at": time.Now(),
		}).Error
}

func (r *GormRepository) UpdateStatusByTask(ctx context.Context, taskID uuid.UUID, exceptID uuid.UUID, status domain.ApplicationStatus) error {
	return r.db.WithContext(ctx).Model(&ApplicationModel{}).
		Where("task_id = ? AND id != ? AND status = 'PENDING'", taskID, exceptID).
		Updates(map[string]any{
			"status":     string(status),
			"updated_at": time.Now(),
		}).Error
}

func (r *GormRepository) CountByTask(ctx context.Context, taskID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ApplicationModel{}).Where("task_id = ?", taskID).Count(&count).Error
	return count, err
}

func toDomain(m ApplicationModel) *domain.Application {
	return &domain.Application{
		ID:                     m.ID,
		TaskID:                 m.TaskID,
		AgentID:                m.AgentID,
		ProposalMessage:        m.ProposalMessage,
		ExpectedCompletionTime: m.ExpectedCompletionTime,
		ProposedPrice:          m.ProposedPrice,
		ETA:                    m.ETA,
		Status:                 domain.ApplicationStatus(m.Status),
		CreatedAt:              m.CreatedAt,
		UpdatedAt:              m.UpdatedAt,
	}
}

func toSlice(models []ApplicationModel) []domain.Application {
	apps := make([]domain.Application, len(models))
	for i, m := range models {
		apps[i] = *toDomain(m)
	}
	return apps
}
