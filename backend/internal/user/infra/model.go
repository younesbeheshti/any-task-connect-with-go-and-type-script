package infra

import (
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/database"
	"gorm.io/gorm"
)

// UserModel is the GORM persistence model for users.
type UserModel struct {
	database.BaseModel
	FullName             string         `gorm:"column:full_name;size:200;not null"`
	Phone                string         `gorm:"column:phone;size:20;not null;uniqueIndex"`
	Email                *string        `gorm:"column:email;size:255;uniqueIndex"`
	PasswordHash         string         `gorm:"column:password_hash;size:255;not null"`
	Role                 string         `gorm:"column:role;type:user_role;not null;default:REQUESTER"`
	NationalID           *string        `gorm:"column:national_id;size:20"`
	Avatar               *string        `gorm:"column:avatar"`
	Bio                  *string        `gorm:"column:bio"`
	CityID               *uuid.UUID     `gorm:"column:city_id;type:uuid"`
	Rating               float64        `gorm:"column:rating;type:numeric(3,2);not null;default:0"`
	RatingCount          int            `gorm:"column:rating_count;not null;default:0"`
	CompletedTasks       int            `gorm:"column:completed_tasks;not null;default:0"`
	IsVerified           bool           `gorm:"column:is_verified;not null;default:false"`
	IsActive             bool           `gorm:"column:is_active;not null;default:true"`
	PhoneVerified        bool           `gorm:"column:phone_verified;not null;default:false"`
	EmailVerified        bool           `gorm:"column:email_verified;not null;default:false"`
	NationalIDVerified   bool           `gorm:"column:national_id_verified;not null;default:false"`
	VerificationLevel    string         `gorm:"column:verification_level;size:20;not null;default:none"`
	VerificationStatus   string         `gorm:"column:verification_status;size:30;not null;default:pending"`
	VerificationReason   *string        `gorm:"column:verification_reason"`
	VerifiedAt           *time.Time     `gorm:"column:verified_at"`
	WalletBalance        int64          `gorm:"column:wallet_balance;not null;default:0"`
	LockedBalance        int64          `gorm:"column:locked_balance;not null;default:0"`
	City                 *CityModel     `gorm:"foreignKey:CityID;references:ID"`
}

func (UserModel) TableName() string { return "users" }

// CityModel is a lightweight city reference.
type CityModel struct {
	ID       uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Title    string    `gorm:"column:title"`
	Province string    `gorm:"column:province"`
}

func (CityModel) TableName() string { return "cities" }

// PermissionModel maps to the permissions table.
type PermissionModel struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Name        string    `gorm:"column:name;size:100;uniqueIndex;not null"`
	Description string    `gorm:"column:description"`
}

func (PermissionModel) TableName() string { return "permissions" }

// RolePermissionModel maps role to permission.
type RolePermissionModel struct {
	Role       common.Role `gorm:"column:role;type:user_role;primaryKey"`
	Permission string      `gorm:"column:permission;size:100;primaryKey"`
}

func (RolePermissionModel) TableName() string { return "role_permissions" }

// RegisterModels registers user-related GORM models.
func RegisterModels(db *gorm.DB) error {
	return db.AutoMigrate(&UserModel{}, &CityModel{}, &PermissionModel{}, &RolePermissionModel{})
}
