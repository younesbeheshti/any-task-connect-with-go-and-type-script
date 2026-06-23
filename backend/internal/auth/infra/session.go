package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	rediscache "github.com/younesbeheshti/any-task-connect/backend/pkg/redis"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/security"
)

// SessionData stored in Redis for refresh token sessions.
type SessionData struct {
	UserID       string    `json:"userId"`
	SessionID    string    `json:"sessionId"`
	RefreshHash  string    `json:"refreshHash"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

// SessionStore manages refresh sessions in Redis.
type SessionStore struct {
	cache *rediscache.Cache
	ttl   time.Duration
}

// NewSessionStore creates a SessionStore.
func NewSessionStore(cache *rediscache.Cache, ttl time.Duration) *SessionStore {
	return &SessionStore{cache: cache, ttl: ttl}
}

// Create stores a new refresh session.
func (s *SessionStore) Create(ctx context.Context, userID uuid.UUID, role common.Role, refreshToken string) (sessionID string, err error) {
	sessionID = uuid.NewString()
	data := SessionData{
		UserID:      userID.String(),
		SessionID:   sessionID,
		RefreshHash: security.HashToken(refreshToken),
		Role:        string(role),
		CreatedAt:   time.Now().UTC(),
		ExpiresAt:   time.Now().UTC().Add(s.ttl),
	}
	key := rediscache.SessionKey(userID.String(), sessionID)
	if err := s.cache.Set(ctx, key, data, s.ttl); err != nil {
		return "", err
	}
	return sessionID, nil
}

// Validate checks refresh token against stored session.
func (s *SessionStore) Validate(ctx context.Context, userID, sessionID, refreshToken string) error {
	key := rediscache.SessionKey(userID, sessionID)
	var data SessionData
	if err := s.cache.Get(ctx, key, &data); err != nil {
		return fmt.Errorf("session not found")
	}
	if data.RefreshHash != security.HashToken(refreshToken) {
		return fmt.Errorf("invalid refresh token")
	}
	if time.Now().UTC().After(data.ExpiresAt) {
		return fmt.Errorf("session expired")
	}
	return nil
}

// Rotate updates session with a new refresh token hash.
func (s *SessionStore) Rotate(ctx context.Context, userID, sessionID, newRefreshToken string) error {
	key := rediscache.SessionKey(userID, sessionID)
	var data SessionData
	if err := s.cache.Get(ctx, key, &data); err != nil {
		return err
	}
	data.RefreshHash = security.HashToken(newRefreshToken)
	data.ExpiresAt = time.Now().UTC().Add(s.ttl)
	return s.cache.Set(ctx, key, data, s.ttl)
}

// Revoke removes a single session.
func (s *SessionStore) Revoke(ctx context.Context, userID, sessionID string) error {
	return s.cache.Delete(ctx, rediscache.SessionKey(userID, sessionID))
}

// RevokeAll removes all sessions for a user (pattern scan simplified via index key).
func (s *SessionStore) RevokeAll(ctx context.Context, userID string) error {
	// Store session index per user for efficient revocation.
	indexKey := fmt.Sprintf("%s%s:index", common.KeySession, userID)
	raw, err := s.cache.GetString(ctx, indexKey)
	if err != nil {
		return nil
	}
	var ids []string
	_ = json.Unmarshal([]byte(raw), &ids)
	keys := make([]string, 0, len(ids)+1)
	for _, id := range ids {
		keys = append(keys, rediscache.SessionKey(userID, id))
	}
	keys = append(keys, indexKey)
	return s.cache.Delete(ctx, keys...)
}

// TrackSession adds session ID to user index.
func (s *SessionStore) TrackSession(ctx context.Context, userID, sessionID string) error {
	indexKey := fmt.Sprintf("%s%s:index", common.KeySession, userID)
	var ids []string
	if raw, err := s.cache.GetString(ctx, indexKey); err == nil {
		_ = json.Unmarshal([]byte(raw), &ids)
	}
	ids = append(ids, sessionID)
	data, _ := json.Marshal(ids)
	return s.cache.SetString(ctx, indexKey, string(data), s.ttl)
}

// OTPData holds OTP verification state in Redis.
type OTPData struct {
	CodeHash     string    `json:"codeHash"`
	ExpiresAt    time.Time `json:"expiresAt"`
	AttemptCount int       `json:"attemptCount"`
	MaxAttempts  int       `json:"maxAttempts"`
}

// OTPStore manages OTP codes in Redis.
type OTPStore struct {
	cache *rediscache.Cache
}

// NewOTPStore creates an OTPStore.
func NewOTPStore(cache *rediscache.Cache) *OTPStore {
	return &OTPStore{cache: cache}
}

func otpKey(channel, target string) string {
	return fmt.Sprintf("otp:%s:%s", channel, target)
}

// Save stores a hashed OTP code.
func (o *OTPStore) Save(ctx context.Context, channel, target, code string, ttl time.Duration, maxAttempts int) error {
	data := OTPData{
		CodeHash:     security.HashToken(code),
		ExpiresAt:    time.Now().UTC().Add(ttl),
		AttemptCount: 0,
		MaxAttempts:  maxAttempts,
	}
	return o.cache.Set(ctx, otpKey(channel, target), data, ttl)
}

// Verify checks an OTP code.
func (o *OTPStore) Verify(ctx context.Context, channel, target, code string) (bool, error) {
	key := otpKey(channel, target)
	var data OTPData
	if err := o.cache.Get(ctx, key, &data); err != nil {
		return false, fmt.Errorf("otp not found")
	}
	if time.Now().UTC().After(data.ExpiresAt) {
		return false, fmt.Errorf("otp expired")
	}
	if data.AttemptCount >= data.MaxAttempts {
		return false, fmt.Errorf("too many attempts")
	}
	data.AttemptCount++
	_ = o.cache.Set(ctx, key, data, time.Until(data.ExpiresAt))

	if security.HashToken(code) != data.CodeHash {
		return false, nil
	}
	_ = o.cache.Delete(ctx, key)
	return true, nil
}

// LockoutStore tracks failed login attempts.
type LockoutStore struct {
	cache *rediscache.Cache
}

// NewLockoutStore creates a LockoutStore.
func NewLockoutStore(cache *rediscache.Cache) *LockoutStore {
	return &LockoutStore{cache: cache}
}

func lockoutKey(identifier string) string {
	return "lockout:" + identifier
}

// RecordFailure increments failed attempts; returns locked status.
func (l *LockoutStore) RecordFailure(ctx context.Context, identifier string, maxAttempts int64, window time.Duration) (bool, error) {
	key := lockoutKey(identifier)
	count, err := l.cache.Incr(ctx, key)
	if err != nil {
		return false, err
	}
	if count == 1 {
		_ = l.cache.Expire(ctx, key, window)
	}
	return count >= maxAttempts, nil
}

// Reset clears lockout counter on success.
func (l *LockoutStore) Reset(ctx context.Context, identifier string) error {
	return l.cache.Delete(ctx, lockoutKey(identifier))
}

// IsLocked reports whether identifier is locked out.
func (l *LockoutStore) IsLocked(ctx context.Context, identifier string, maxAttempts int64) (bool, error) {
	key := lockoutKey(identifier)
	val, err := l.cache.GetString(ctx, key)
	if err != nil {
		return false, nil
	}
	var attempts int64
	if _, scanErr := fmt.Sscanf(val, "%d", &attempts); scanErr != nil {
		return false, nil
	}
	return attempts >= maxAttempts, nil
}
