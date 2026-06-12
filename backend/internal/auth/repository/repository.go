package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/auth/domain"
)

// Repository defines auth persistence operations.
type Repository interface {
	CreateSession(ctx context.Context, session *domain.Session) error
	GetSessionByTokenHash(ctx context.Context, tokenHash string) (*domain.Session, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error

	CreatePasswordResetToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt int64) error
	GetPasswordResetToken(ctx context.Context, tokenHash string) (userID uuid.UUID, err error)
	MarkPasswordResetTokenUsed(ctx context.Context, tokenHash string) error

	SavePhoneOTP(ctx context.Context, phone, codeHash string, expiresAt int64) error
	VerifyPhoneOTP(ctx context.Context, phone, codeHash string) (bool, error)
}
