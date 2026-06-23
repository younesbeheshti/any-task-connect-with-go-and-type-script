package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/younesbeheshti/any-task-connect/backend/internal/auth/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/auth/infra"
	authrepo "github.com/younesbeheshti/any-task-connect/backend/internal/auth/repository"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	userdomain "github.com/younesbeheshti/any-task-connect/backend/internal/user/domain"
	userrepo "github.com/younesbeheshti/any-task-connect/backend/internal/user/repository"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/jwt"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/phone"
	rediscache "github.com/younesbeheshti/any-task-connect/backend/pkg/redis"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/security"
)

const (
	maxLoginAttempts  int64 = 5
	lockoutWindow           = 15 * time.Minute
	otpTTL                  = 5 * time.Minute
	otpMaxAttempts          = 5
	passwordResetTTL        = 1 * time.Hour
)

// AuthService implements authentication use cases.
type AuthService struct {
	users      userrepo.Repository
	authRepo   authrepo.Repository
	jwt        *jwt.Service
	cache      *rediscache.Cache
	sessions   *infra.SessionStore
	otp        *infra.OTPStore
	lockout    *infra.LockoutStore
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewAuthService creates an AuthService.
func NewAuthService(
	users userrepo.Repository,
	authRepo authrepo.Repository,
	jwtSvc *jwt.Service,
	cache *rediscache.Cache,
	sessions *infra.SessionStore,
	otp *infra.OTPStore,
	lockout *infra.LockoutStore,
	accessTTL, refreshTTL time.Duration,
) *AuthService {
	return &AuthService{
		users: users, authRepo: authRepo, jwt: jwtSvc, cache: cache,
		sessions: sessions, otp: otp, lockout: lockout,
		accessTTL: accessTTL, refreshTTL: refreshTTL,
	}
}

// Register creates a new user account.
func (s *AuthService) Register(ctx context.Context, input authdomain.RegisterInput) (*authdomain.TokenPair, *userdomain.User, error) {

	log.Println("input", input)

	normalized, err := phone.Normalize(input.Phone)
	if err != nil {
		return nil, nil, apperrors.Validation(map[string]string{"phone": "شماره موبایل معتبر نیست"})
	}
	if !input.Role.IsValid() || input.Role == common.RoleAdmin {
		return nil, nil, apperrors.Validation(map[string]string{"role": "نقش انتخاب شده معتبر نیست"})
	}

	log.Println("normalized", normalized)

	if _, err := s.users.GetByPhone(ctx, normalized); err == nil {
		return nil, nil, apperrors.New("CONFLICT", "این شماره موبایل قبلاً ثبت شده است", 409, common.ErrDuplicate)
	} else if !errors.Is(err, common.ErrNotFound) {
		return nil, nil, err
	}

	hash, err := security.HashPassword(input.Password)
	if err != nil {
		return nil, nil, err
	}

	log.Println("hash", hash)

	user := &userdomain.User{
		ID: uuid.New(), FullName: input.FullName, Phone: normalized,
		Email: input.Email, PasswordHash: hash, Role: input.Role,
		IsActive: true, VerificationLevel: "none", VerificationStatus: "pending",
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}

	log.Println("user", user)
	log.Println("what is happenning here?")

	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, common.ErrDuplicate) {
			return nil, nil, apperrors.New("CONFLICT", "این شماره موبایل قبلاً ثبت شده است", 409, err)
		}
		return nil, nil, err
	}

	return s.issueTokens(ctx, user)
}

// Login authenticates by phone or email.
func (s *AuthService) Login(ctx context.Context, input authdomain.LoginInput) (*authdomain.TokenPair, *userdomain.User, error) {
	identifier := input.Phone
	if input.Email != "" {
		identifier = input.Email
	}

	locked, err := s.lockout.IsLocked(ctx, identifier, maxLoginAttempts)
	if err == nil && locked {
		return nil, nil, apperrors.New("RATE_LIMITED", "حساب به دلیل تلاش‌های ناموفق موقتاً قفل شده است", 429, common.ErrForbidden)
	}

	var user *userdomain.User
	if input.Email != "" {
		user, err = s.users.GetByEmail(ctx, input.Email)
	} else {
		normalized, normErr := phone.Normalize(input.Phone)
		if normErr != nil {
			return nil, nil, apperrors.Validation(map[string]string{"phone": "شماره موبایل معتبر نیست"})
		}
		user, err = s.users.GetByPhone(ctx, normalized)
	}
	if errors.Is(err, common.ErrNotFound) {
		_, _ = s.lockout.RecordFailure(ctx, identifier, maxLoginAttempts, lockoutWindow)
		return nil, nil, apperrors.New("UNAUTHORIZED", "شماره موبایل یا رمز عبور اشتباه است", 401, common.ErrUnauthorized)
	}
	if err != nil {
		return nil, nil, err
	}
	if !user.IsActive {
		return nil, nil, apperrors.New("FORBIDDEN", "حساب کاربری غیرفعال است", 403, common.ErrForbidden)
	}
	if err := security.VerifyPassword(user.PasswordHash, input.Password); err != nil {
		_, _ = s.lockout.RecordFailure(ctx, identifier, maxLoginAttempts, lockoutWindow)
		return nil, nil, apperrors.New("UNAUTHORIZED", "شماره موبایل یا رمز عبور اشتباه است", 401, common.ErrUnauthorized)
	}
	_ = s.lockout.Reset(ctx, identifier)
	return s.issueTokens(ctx, user)
}

func (s *AuthService) issueTokens(ctx context.Context, user *userdomain.User) (*authdomain.TokenPair, *userdomain.User, error) {
	access, err := s.jwt.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, nil, err
	}
	refresh, err := s.jwt.GenerateRefreshToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, nil, err
	}

	sessionID, err := s.sessions.Create(ctx, user.ID, user.Role, refresh)
	if err != nil {
		return nil, nil, err
	}
	_ = s.sessions.TrackSession(ctx, user.ID.String(), sessionID)

	tokenHash := security.HashToken(refresh)
	session := &authdomain.Session{
		ID: uuid.New(), UserID: user.ID, TokenHash: tokenHash,
		ExpiresAt: time.Now().UTC().Add(s.refreshTTL), CreatedAt: time.Now().UTC(),
	}
	_ = s.authRepo.CreateSession(ctx, session)

	return &authdomain.TokenPair{
		AccessToken: access, RefreshToken: refresh,
		ExpiresIn: int64(s.accessTTL.Seconds()),
		SessionID: sessionID,
	}, user, nil
}

// Refresh rotates tokens using refresh token.
func (s *AuthService) Refresh(ctx context.Context, input authdomain.RefreshInput) (*authdomain.TokenPair, error) {
	claims, err := s.jwt.ParseClaims(input.RefreshToken)
	if err != nil {
		return nil, apperrors.New("UNAUTHORIZED", "توکن نامعتبر است", 401, common.ErrUnauthorized)
	}
	if claims.TokenType != jwt.TokenTypeRefresh {
		return nil, apperrors.New("UNAUTHORIZED", "توکن نامعتبر است", 401, common.ErrUnauthorized)
	}

	tokenHash := security.HashToken(input.RefreshToken)
	session, err := s.authRepo.GetSessionByTokenHash(ctx, tokenHash)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("UNAUTHORIZED", "نشست منقضی شده است", 401, common.ErrUnauthorized)
	}
	if err != nil {
		return nil, err
	}
	if time.Now().UTC().After(session.ExpiresAt) {
		return nil, apperrors.New("UNAUTHORIZED", "نشست منقضی شده است", 401, common.ErrUnauthorized)
	}

	userID, _ := uuid.Parse(claims.UserID)
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	_ = s.authRepo.RevokeSession(ctx, session.ID)

	access, err := s.jwt.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}
	refresh, err := s.jwt.GenerateRefreshToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}

	newSession := &authdomain.Session{
		ID: uuid.New(), UserID: user.ID, TokenHash: security.HashToken(refresh),
		ExpiresAt: time.Now().UTC().Add(s.refreshTTL), CreatedAt: time.Now().UTC(),
	}
	_ = s.authRepo.CreateSession(ctx, newSession)

	return &authdomain.TokenPair{
		AccessToken: access, RefreshToken: refresh,
		ExpiresIn: int64(s.accessTTL.Seconds()),
	}, nil
}

// Logout revokes current session and blacklists access token JTI.
func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID, jti string) error {
	if jti != "" {
		_ = s.cache.BlacklistJTI(ctx, jti, s.accessTTL)
	}
	return s.authRepo.RevokeAllUserSessions(ctx, userID)
}

// LogoutDevice revokes a single session.
func (s *AuthService) LogoutDevice(ctx context.Context, userID uuid.UUID, sessionID, jti string) error {
	if jti != "" {
		_ = s.cache.BlacklistJTI(ctx, jti, s.accessTTL)
	}
	_ = s.sessions.Revoke(ctx, userID.String(), sessionID)
	return nil
}

// LogoutAll revokes all sessions for user.
func (s *AuthService) LogoutAll(ctx context.Context, userID uuid.UUID, jti string) error {
	if jti != "" {
		_ = s.cache.BlacklistJTI(ctx, jti, s.accessTTL)
	}
	_ = s.sessions.RevokeAll(ctx, userID.String())
	return s.authRepo.RevokeAllUserSessions(ctx, userID)
}

// ChangePassword updates password for authenticated user.
func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := security.VerifyPassword(user.PasswordHash, currentPassword); err != nil {
		return apperrors.New("UNAUTHORIZED", "رمز عبور فعلی اشتباه است", 401, common.ErrUnauthorized)
	}
	hash, err := security.HashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := s.users.UpdatePassword(ctx, userID, hash); err != nil {
		return err
	}
	return s.authRepo.RevokeAllUserSessions(ctx, userID)
}

// ForgotPassword initiates password reset (structure ready).
func (s *AuthService) ForgotPassword(ctx context.Context, input authdomain.ForgotPasswordInput) error {
	normalized, err := phone.Normalize(input.Phone)
	if err != nil {
		return apperrors.Validation(map[string]string{"phone": "شماره موبایل معتبر نیست"})
	}
	user, err := s.users.GetByPhone(ctx, normalized)
	if errors.Is(err, common.ErrNotFound) {
		return nil // do not reveal user existence
	}
	if err != nil {
		return err
	}
	token, err := security.RandomToken(32)
	if err != nil {
		return err
	}
	expires := time.Now().UTC().Add(passwordResetTTL).Unix()
	return s.authRepo.CreatePasswordResetToken(ctx, user.ID, security.HashToken(token), expires)
}

// ResetPassword completes password reset.
func (s *AuthService) ResetPassword(ctx context.Context, input authdomain.ResetPasswordInput) error {
	userID, err := s.authRepo.GetPasswordResetToken(ctx, security.HashToken(input.Token))
	if errors.Is(err, common.ErrNotFound) {
		return apperrors.New("VALIDATION_FAILED", "توکن بازیابی نامعتبر یا منقضی شده", 400, common.ErrValidation)
	}
	if err != nil {
		return err
	}
	hash, err := security.HashPassword(input.NewPassword)
	if err != nil {
		return err
	}
	if err := s.users.UpdatePassword(ctx, userID, hash); err != nil {
		return err
	}
	_ = s.authRepo.MarkPasswordResetTokenUsed(ctx, security.HashToken(input.Token))
	return s.authRepo.RevokeAllUserSessions(ctx, userID)
}

// SendPhoneOTP sends OTP to phone.
func (s *AuthService) SendPhoneOTP(ctx context.Context, phoneRaw string) error {
	normalized, err := phone.Normalize(phoneRaw)
	if err != nil {
		return apperrors.Validation(map[string]string{"phone": "شماره موبایل معتبر نیست"})
	}
	code, err := security.RandomToken(3)
	if err != nil {
		return err
	}
	// Use first 6 hex chars as numeric-like code for dev
	code = code[:6]
	return s.otp.Save(ctx, "phone", normalized, code, otpTTL, otpMaxAttempts)
}

// VerifyPhoneOTP verifies phone OTP.
func (s *AuthService) VerifyPhoneOTP(ctx context.Context, input authdomain.OTPInput) error {
	normalized, err := phone.Normalize(input.Phone)
	if err != nil {
		return apperrors.Validation(map[string]string{"phone": "شماره موبایل معتبر نیست"})
	}
	ok, err := s.otp.Verify(ctx, "phone", normalized, input.Code)
	if err != nil {
		return apperrors.New("VALIDATION_FAILED", err.Error(), 400, common.ErrValidation)
	}
	if !ok {
		return apperrors.New("VALIDATION_FAILED", "کد تأیید اشتباه است", 400, common.ErrValidation)
	}
	user, err := s.users.GetByPhone(ctx, normalized)
	if err == nil {
		user.PhoneVerified = true
		if user.VerificationLevel == "none" {
			user.VerificationLevel = "basic"
		}
		_ = s.users.Update(ctx, user)
	}
	return nil
}

// SendEmailVerification prepares email verification (structure).
func (s *AuthService) SendEmailVerification(ctx context.Context, userID uuid.UUID) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.Email == nil {
		return apperrors.Validation(map[string]string{"email": "ایمیل ثبت نشده است"})
	}
	code, err := security.RandomToken(3)
	if err != nil {
		return err
	}
	return s.otp.Save(ctx, "email", *user.Email, code[:6], otpTTL, otpMaxAttempts)
}

// VerifyEmail verifies email OTP/token.
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	_ = token
	return fmt.Errorf("not implemented")
}

// ValidateAccessToken validates JWT and checks blacklist.
func (s *AuthService) ValidateAccessToken(ctx context.Context, token string) (*authdomain.JWTClaims, error) {
	claims, err := s.jwt.ValidateAccessToken(token)
	if err != nil {
		return nil, err
	}
	blacklisted, err := s.cache.IsJTIBlacklisted(ctx, claims.JTI)
	if err == nil && blacklisted {
		return nil, common.ErrUnauthorized
	}
	userID, _ := uuid.Parse(claims.UserID)
	return &authdomain.JWTClaims{UserID: userID, Role: claims.Role, JTI: claims.JTI}, nil
}
