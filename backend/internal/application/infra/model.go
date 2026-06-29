package infra

import (
	"time"

	"github.com/google/uuid"
)

// AgentModel is a minimal user reference for application associations.
type AgentModel struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	FullName       string    `gorm:"column:full_name"`
	Avatar         *string   `gorm:"column:avatar"`
	Rating         float64   `gorm:"column:rating"`
	CompletedTasks int       `gorm:"column:completed_tasks"`
	IsVerified     bool      `gorm:"column:is_verified"`
}

func (AgentModel) TableName() string { return "users" }

// CityRefModel is a minimal city reference for resolving a task's city title.
type CityRefModel struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title string    `gorm:"column:title"`
}

func (CityRefModel) TableName() string { return "cities" }

// TaskRefModel is a minimal task reference for an agent's application listing.
type TaskRefModel struct {
	ID       uuid.UUID    `gorm:"type:uuid;primaryKey"`
	PublicID string       `gorm:"column:public_id"`
	Title    string       `gorm:"column:title"`
	Budget   int64        `gorm:"column:budget"`
	Deadline time.Time    `gorm:"column:deadline"`
	Status   string       `gorm:"column:status"`
	CityID   *uuid.UUID   `gorm:"column:city_id"`
	City     CityRefModel `gorm:"foreignKey:CityID;references:ID"`
}

func (TaskRefModel) TableName() string { return "tasks" }

// ApplicationModel is the GORM persistence model for applications.
type ApplicationModel struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TaskID                 uuid.UUID  `gorm:"type:uuid;not null;index"`
	AgentID                uuid.UUID  `gorm:"type:uuid;not null"`
	ProposalMessage        string     `gorm:"type:text;not null"`
	ExpectedCompletionTime *time.Time `gorm:"type:date"`
	ProposedPrice          *int64
	ETA                    string    `gorm:"size:100"`
	Status                 string    `gorm:"size:20;not null;default:'PENDING'"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	Agent                  AgentModel   `gorm:"foreignKey:AgentID;references:ID"`
	Task                   TaskRefModel `gorm:"foreignKey:TaskID;references:ID"`
}

func (ApplicationModel) TableName() string { return "applications" }
