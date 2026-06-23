package dto

import "time"

// Verification matches frontend User.verification schema.
type Verification struct {
	Phone      bool   `json:"phone"`
	Email      bool   `json:"email"`
	NationalID bool   `json:"nationalId"`
	Level      string `json:"level"`
}

// User matches front/docs/api-contracts.md User schema.
type User struct {
	ID             string       `json:"id"`
	FullName       string       `json:"fullName"`
	Phone          string       `json:"phone"`
	Email          *string      `json:"email"`
	City           string       `json:"city,omitempty"`
	Role           string       `json:"role"`
	AvatarURL      *string      `json:"avatarUrl"`
	Verification   Verification `json:"verification"`
	Rating         float64      `json:"rating"`
	CompletedCount int          `json:"completedCount"`
	CreatedAt      time.Time    `json:"createdAt"`
}

// AuthResponse is the login/register response body.
type AuthResponse struct {
	Token        string   `json:"token"`
	RefreshToken string   `json:"refreshToken,omitempty"`
	User         User     `json:"user"`
	Permissions  []string `json:"permissions,omitempty"`
	Role         string   `json:"role,omitempty"`
}
