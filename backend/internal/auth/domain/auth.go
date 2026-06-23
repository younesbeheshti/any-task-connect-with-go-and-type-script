package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
)

// TokenPair holds access and refresh tokens returned on auth.
type TokenPair struct {
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ExpiresIn    int64  `json:"expiresIn"`
	SessionID    string `json:"-"`
}

// RegisterInput holds registration request data.
type RegisterInput struct {
	FullName string
	Phone    string
	Email    *string
	Password string
	Role     common.Role
}

// LoginInput holds login credentials.
type LoginInput struct {
	Phone    string
	Email    string
	Password string
}

// RefreshInput holds token refresh data.
type RefreshInput struct {
	RefreshToken string
}

// ForgotPasswordInput initiates password reset.
type ForgotPasswordInput struct {
	Phone string
}

// ResetPasswordInput completes password reset.
type ResetPasswordInput struct {
	Token       string
	NewPassword string
}

// OTPInput holds phone OTP verification data.
type OTPInput struct {
	Phone string
	Code  string
}

// Session represents a stored refresh session.
type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// JWTClaims are embedded in access tokens.
type JWTClaims struct {
	UserID uuid.UUID
	Role   common.Role
	JTI    string
}
