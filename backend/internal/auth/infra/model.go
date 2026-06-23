package infra

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshTokenModel maps to refresh_tokens table.
type RefreshTokenModel struct {
	ID        uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID      `gorm:"column:user_id;type:uuid;not null;index"`
	TokenHash string         `gorm:"column:token_hash;size:255;not null;uniqueIndex"`
	ExpiresAt time.Time      `gorm:"column:expires_at;not null"`
	RevokedAt *time.Time     `gorm:"column:revoked_at"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
}

func (RefreshTokenModel) TableName() string { return "refresh_tokens" }

// PasswordResetModel maps to password_reset_tokens table.
type PasswordResetModel struct {
	ID        uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID  `gorm:"column:user_id;type:uuid;not null;index"`
	TokenHash string     `gorm:"column:token_hash;size:255;not null;uniqueIndex"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null"`
	UsedAt    *time.Time `gorm:"column:used_at"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (PasswordResetModel) TableName() string { return "password_reset_tokens" }

// RegisterModels registers auth-related GORM models.
func RegisterModels(db *gorm.DB) error {
	return db.AutoMigrate(&RefreshTokenModel{}, &PasswordResetModel{})
}
