package infra

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/auth/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"gorm.io/gorm"
)

// GormRepository implements auth repository.
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a GormRepository.
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreateSession(ctx context.Context, session *domain.Session) error {
	m := RefreshTokenModel{
		ID: session.ID, UserID: session.UserID,
		TokenHash: session.TokenHash, ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
	}
	if session.RevokedAt != nil {
		m.RevokedAt = session.RevokedAt
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *GormRepository) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*domain.Session, error) {
	var m RefreshTokenModel
	err := r.db.WithContext(ctx).Where("token_hash = ? AND revoked_at IS NULL", tokenHash).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &domain.Session{
		ID: m.ID, UserID: m.UserID, TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt, RevokedAt: m.RevokedAt, CreatedAt: m.CreatedAt,
	}, nil
}

func (r *GormRepository) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&RefreshTokenModel{}).Where("id = ?", sessionID).
		Update("revoked_at", now).Error
}

func (r *GormRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&RefreshTokenModel{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}

func (r *GormRepository) CreatePasswordResetToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt int64) error {
	m := PasswordResetModel{
		UserID: userID, TokenHash: tokenHash,
		ExpiresAt: time.Unix(expiresAt, 0).UTC(),
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *GormRepository) GetPasswordResetToken(ctx context.Context, tokenHash string) (uuid.UUID, error) {
	var m PasswordResetModel
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", tokenHash, time.Now().UTC()).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, common.ErrNotFound
	}
	if err != nil {
		return uuid.Nil, err
	}
	return m.UserID, nil
}

func (r *GormRepository) MarkPasswordResetTokenUsed(ctx context.Context, tokenHash string) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&PasswordResetModel{}).
		Where("token_hash = ?", tokenHash).Update("used_at", now).Error
}

func (r *GormRepository) SavePhoneOTP(ctx context.Context, phone, codeHash string, expiresAt int64) error {
	_ = ctx
	_ = phone
	_ = codeHash
	_ = expiresAt
	return nil // OTP stored in Redis via OTPStore
}

func (r *GormRepository) VerifyPhoneOTP(ctx context.Context, phone, codeHash string) (bool, error) {
	_ = ctx
	_ = phone
	_ = codeHash
	return false, nil
}
