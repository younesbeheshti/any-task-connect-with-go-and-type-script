package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/task/domain"
	"gorm.io/gorm"
)

// GormRepository implements task/repository.Repository using GORM.
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new task GormRepository.
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, task *domain.Task) error {
	urlsJSON, err := marshalURLs(task.AttachmentURLs)
	if err != nil {
		return err
	}
	m := TaskModel{
		ID:               task.ID,
		PublicID:         task.PublicID,
		Title:            task.Title,
		Description:      task.Description,
		CategoryID:       task.CategoryID,
		CityID:           task.CityID,
		Budget:           task.Budget,
		EscrowFee:        task.EscrowFee,
		Status:           string(task.Status),
		Deadline:         task.Deadline,
		RequesterID:      task.RequesterID,
		AssignedAgentID:  task.AssignedAgentID,
		AttachmentURLs:   urlsJSON,
		ApplicationCount: task.ApplicantsCount,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	task.ID = m.ID
	task.CreatedAt = m.CreatedAt
	task.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	var m TaskModel
	err := r.db.WithContext(ctx).Preload("Category").Preload("City").
		First(&m, "tasks.id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m)
}

func (r *GormRepository) GetByPublicID(ctx context.Context, publicID string) (*domain.Task, error) {
	var m TaskModel
	err := r.db.WithContext(ctx).Preload("Category").Preload("City").
		Where("tasks.public_id = ?", publicID).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m)
}

func (r *GormRepository) Update(ctx context.Context, task *domain.Task) error {
	urlsJSON, err := marshalURLs(task.AttachmentURLs)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&TaskModel{}).Where("id = ?", task.ID).
		Updates(map[string]any{
			"title":           task.Title,
			"description":     task.Description,
			"category_id":     task.CategoryID,
			"city_id":         task.CityID,
			"budget":          task.Budget,
			"deadline":        task.Deadline,
			"attachment_urls": urlsJSON,
			"updated_at":      time.Now(),
		}).Error
}

func (r *GormRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus) error {
	return r.db.WithContext(ctx).Model(&TaskModel{}).Where("id = ?", id).
		Updates(map[string]any{
			"status":     string(status),
			"updated_at": time.Now(),
		}).Error
}

func (r *GormRepository) List(ctx context.Context, filter domain.TaskFilter, pg common.PaginationParams) ([]domain.Task, int64, error) {
	q := r.db.WithContext(ctx).Model(&TaskModel{}).Preload("Category").Preload("City")

	if filter.CityID != nil {
		q = q.Where("tasks.city_id = ?", filter.CityID)
	} else if filter.CityTitle != "" {
		q = q.Joins("JOIN cities ON cities.id = tasks.city_id").Where("cities.title = ?", filter.CityTitle)
	}
	if filter.CategoryID != nil {
		q = q.Where("tasks.category_id = ?", filter.CategoryID)
	} else if filter.CatTitle != "" {
		q = q.Joins("JOIN categories ON categories.id = tasks.category_id").Where("categories.title = ?", filter.CatTitle)
	}
	if filter.MinBudget != nil {
		q = q.Where("tasks.budget >= ?", *filter.MinBudget)
	}
	if filter.MaxBudget != nil {
		q = q.Where("tasks.budget <= ?", *filter.MaxBudget)
	}
	if filter.Status != nil {
		q = q.Where("status = ?", string(*filter.Status))
	}
	if filter.RequesterID != nil {
		q = q.Where("requester_id = ?", filter.RequesterID)
	}
	if filter.AgentID != nil {
		q = q.Where("assigned_agent_id = ?", filter.AgentID)
	}
	if filter.Query != "" {
		like := "%" + filter.Query + "%"
		q = q.Where("title ILIKE ? OR description ILIKE ?", like, like)
	}

	// Only non-deleted
	q = q.Where("deleted_at IS NULL")

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch filter.Sort {
	case "budget":
		q = q.Order("budget desc")
	case "deadline":
		q = q.Order("deadline asc")
	default:
		q = q.Order("created_at desc")
	}

	var models []TaskModel
	if err := q.Offset(pg.Offset()).Limit(pg.Limit()).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	tasks := make([]domain.Task, 0, len(models))
	for _, m := range models {
		t, err := toDomain(m)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, *t)
	}
	return tasks, total, nil
}

func (r *GormRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&TaskModel{}).Error
}

func (r *GormRepository) AppendTimeline(ctx context.Context, entry *domain.TaskTimeline) error {
	var fromStatus *string
	if entry.FromStatus != nil {
		s := string(*entry.FromStatus)
		fromStatus = &s
	}
	var actorID *uuid.UUID
	if entry.ActorID != nil {
		actorID = entry.ActorID
	}
	m := TaskTimelineModel{
		ID:         entry.ID,
		TaskID:     entry.TaskID,
		FromStatus: fromStatus,
		ToStatus:   string(entry.ToStatus),
		ActorID:    actorID,
		Note:       entry.Note,
	}
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	entry.ID = m.ID
	entry.CreatedAt = m.CreatedAt
	return nil
}

func (r *GormRepository) GetTimeline(ctx context.Context, taskID uuid.UUID) ([]domain.TaskTimeline, error) {
	var models []TaskTimelineModel
	if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Order("created_at asc").Find(&models).Error; err != nil {
		return nil, err
	}
	entries := make([]domain.TaskTimeline, len(models))
	for i, m := range models {
		var fromStatus *domain.TaskStatus
		if m.FromStatus != nil {
			s := domain.TaskStatus(*m.FromStatus)
			fromStatus = &s
		}
		toStatus := domain.TaskStatus(m.ToStatus)
		entries[i] = domain.TaskTimeline{
			ID:         m.ID,
			TaskID:     m.TaskID,
			FromStatus: fromStatus,
			ToStatus:   toStatus,
			ActorID:    m.ActorID,
			Note:       m.Note,
			CreatedAt:  m.CreatedAt,
		}
	}
	return entries, nil
}

func (r *GormRepository) NextPublicID(ctx context.Context) (string, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&TaskModel{}).Count(&count).Error; err != nil {
		return "", fmt.Errorf("next public id: %w", err)
	}
	return fmt.Sprintf("TB-%d", count+1), nil
}

// marshalURLs converts []string to JSON text for storage.
func marshalURLs(urls []string) (string, error) {
	if urls == nil {
		urls = []string{}
	}
	b, err := json.Marshal(urls)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// unmarshalURLs converts stored JSON text back to []string.
func unmarshalURLs(s string) ([]string, error) {
	if s == "" || s == "null" {
		return []string{}, nil
	}
	var urls []string
	if err := json.Unmarshal([]byte(s), &urls); err != nil {
		return nil, err
	}
	return urls, nil
}

func toDomain(m TaskModel) (*domain.Task, error) {
	urls, err := unmarshalURLs(m.AttachmentURLs)
	if err != nil {
		return nil, err
	}
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}
	return &domain.Task{
		ID:              m.ID,
		PublicID:        m.PublicID,
		Title:           m.Title,
		Description:     m.Description,
		CategoryID:      m.CategoryID,
		CityID:          m.CityID,
		Category:        m.Category.Title,
		City:            m.City.Title,
		Budget:          m.Budget,
		EscrowFee:       m.EscrowFee,
		Currency:        "IRR",
		Status:          domain.TaskStatus(m.Status),
		Deadline:        m.Deadline,
		RequesterID:     m.RequesterID,
		AssignedAgentID: m.AssignedAgentID,
		AttachmentURLs:  urls,
		ApplicantsCount: m.ApplicationCount,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		DeletedAt:       deletedAt,
	}, nil
}
