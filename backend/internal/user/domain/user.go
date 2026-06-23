package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
)

// User is the core user aggregate root.
type User struct {
	ID                   uuid.UUID  `json:"id"`
	FullName             string     `json:"fullName"`
	Phone                string     `json:"phone"`
	Email                *string    `json:"email"`
	PasswordHash         string     `json:"-"`
	Role                 common.Role `json:"role"`
	NationalID           *string    `json:"nationalId,omitempty"`
	Avatar               *string    `json:"avatarUrl"`
	Bio                  *string    `json:"bio,omitempty"`
	CityID               *uuid.UUID `json:"-"`
	Rating               float64    `json:"rating"`
	RatingCount          int        `json:"ratingCount"`
	CompletedTasks       int        `json:"completedCount"`
	IsVerified           bool       `json:"isVerified"`
	IsActive             bool       `json:"isActive"`
	PhoneVerified        bool       `json:"-"`
	EmailVerified        bool       `json:"-"`
	NationalIDVerified   bool       `json:"-"`
	VerificationLevel    string     `json:"-"`
	VerificationStatus   string     `json:"-"`
	VerificationReason   *string    `json:"-"`
	VerifiedAt           *time.Time `json:"-"`
	WalletBalance        int64      `json:"-"`
	LockedBalance        int64      `json:"-"`
	CityTitle            string     `json:"city,omitempty"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`
	DeletedAt            *time.Time `json:"-"`
}

// Verification holds public verification status for API responses.
type Verification struct {
	Phone      bool   `json:"phone"`
	Email      bool   `json:"email"`
	NationalID bool   `json:"nationalId"`
	Level      string `json:"level"`
}

func (u *User) Verification() Verification {
	return Verification{
		Phone:      u.PhoneVerified,
		Email:      u.EmailVerified,
		NationalID: u.NationalIDVerified,
		Level:      u.VerificationLevel,
	}
}

// PublicProfile is a safe subset for public user pages.
type PublicProfile struct {
	ID             uuid.UUID    `json:"id"`
	FullName       string       `json:"fullName"`
	Avatar         *string      `json:"avatarUrl"`
	Rating         float64      `json:"rating"`
	CompletedCount int          `json:"completedCount"`
	Verification   Verification `json:"verification"`
	CreatedAt      time.Time    `json:"createdAt"`
}

func (u *User) ToPublicProfile() PublicProfile {
	return PublicProfile{
		ID:             u.ID,
		FullName:       u.FullName,
		Avatar:         u.Avatar,
		Rating:         u.Rating,
		CompletedCount: u.CompletedTasks,
		Verification:   u.Verification(),
		CreatedAt:      u.CreatedAt,
	}
}

// UserStats holds user statistics for profile view.
type UserStats struct {
	Rating         float64 `json:"rating"`
	CompletedCount int     `json:"completedCount"`
	RatingCount    int     `json:"ratingCount"`
	WalletBalance  int64   `json:"walletBalance"`
	LockedBalance  int64   `json:"lockedBalance"`
}

type UpdateProfileInput struct {
	FullName  *string
	Email     *string
	Avatar    *string
	Bio       *string
	CityID    *uuid.UUID
	NationalID *string
}
