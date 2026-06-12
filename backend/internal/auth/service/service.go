package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/auth/domain"
	userdomain "github.com/younesbeheshti/any-task-connect/backend/internal/user/domain"
)

// Service defines authentication use cases.
type Service interface {
	Register(ctx context.Context, input domain.RegisterInput) (*domain.TokenPair, *userdomain.User, error)
	Login(ctx context.Context, input domain.LoginInput) (*domain.TokenPair, *userdomain.User, error)
	Refresh(ctx context.Context, input domain.RefreshInput) (*domain.TokenPair, error)
	Logout(ctx context.Context, userID uuid.UUID, jti string) error

	ForgotPassword(ctx context.Context, input domain.ForgotPasswordInput) error
	ResetPassword(ctx context.Context, input domain.ResetPasswordInput) error

	SendPhoneOTP(ctx context.Context, phone string) error
	VerifyPhoneOTP(ctx context.Context, input domain.OTPInput) error

	SendEmailVerification(ctx context.Context, userID uuid.UUID) error
	VerifyEmail(ctx context.Context, token string) error

	ValidateAccessToken(token string) (*domain.JWTClaims, error)
}
