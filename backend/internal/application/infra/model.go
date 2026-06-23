package infra

import (
	"time"

	"github.com/google/uuid"
)

// AgentModel is a minimal user reference for application associations.
type AgentModel struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	FullName string    `gorm:"column:full_name"`
	Avatar   *string   `gorm:"column:avatar"`
	Rating   float64   `gorm:"column:rating"`
}

func (AgentModel) TableName() string { return "users" }

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
	Agent                  AgentModel `gorm:"foreignKey:AgentID;references:ID"`
}

func (ApplicationModel) TableName() string { return "applications" }
