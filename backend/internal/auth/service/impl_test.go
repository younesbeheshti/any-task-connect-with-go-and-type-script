package service_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/younesbeheshti/any-task-connect/backend/configs"
	authdomain "github.com/younesbeheshti/any-task-connect/backend/internal/auth/domain"
	authinfra "github.com/younesbeheshti/any-task-connect/backend/internal/auth/infra"
	authservice "github.com/younesbeheshti/any-task-connect/backend/internal/auth/service"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	userinfra "github.com/younesbeheshti/any-task-connect/backend/internal/user/infra"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/jwt"
	rediscache "github.com/younesbeheshti/any-task-connect/backend/pkg/redis"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthTest(t *testing.T) *authservice.AuthService {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.Exec(`CREATE TABLE users (
		id TEXT PRIMARY KEY, full_name TEXT NOT NULL, phone TEXT NOT NULL UNIQUE,
		email TEXT UNIQUE, password_hash TEXT NOT NULL, role TEXT NOT NULL DEFAULT 'REQUESTER',
		national_id TEXT, avatar TEXT, bio TEXT, city_id TEXT,
		rating REAL NOT NULL DEFAULT 0, rating_count INTEGER NOT NULL DEFAULT 0,
		completed_tasks INTEGER NOT NULL DEFAULT 0, is_verified INTEGER NOT NULL DEFAULT 0,
		is_active INTEGER NOT NULL DEFAULT 1, phone_verified INTEGER NOT NULL DEFAULT 0,
		email_verified INTEGER NOT NULL DEFAULT 0, national_id_verified INTEGER NOT NULL DEFAULT 0,
		verification_level TEXT NOT NULL DEFAULT 'none', verification_status TEXT NOT NULL DEFAULT 'pending',
		verification_reason TEXT, verified_at DATETIME,
		wallet_balance INTEGER NOT NULL DEFAULT 0, locked_balance INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME
	)`).Error)
	require.NoError(t, db.Exec(`CREATE TABLE refresh_tokens (
		id TEXT PRIMARY KEY, user_id TEXT NOT NULL, token_hash TEXT NOT NULL UNIQUE,
		expires_at DATETIME NOT NULL, revoked_at DATETIME, created_at DATETIME
	)`).Error)

	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	cache := rediscache.NewCache(rdb)

	jwtCfg := configs.JWTConfig{
		Secret: "test-secret-key-minimum-32-characters-long",
		AccessTTL: 900, RefreshTTL: 604800,
	}
	jwtSvc := jwt.NewService(jwtCfg)

	return authservice.NewAuthService(
		userinfra.NewGormRepository(db),
		authinfra.NewGormRepository(db),
		jwtSvc, cache,
		authinfra.NewSessionStore(cache, jwtCfg.RefreshTokenDuration()),
		authinfra.NewOTPStore(cache),
		authinfra.NewLockoutStore(cache),
		jwtCfg.AccessTokenDuration(), jwtCfg.RefreshTokenDuration(),
	)
}

func TestRegisterAndLogin(t *testing.T) {
	svc := setupAuthTest(t)
	ctx := context.Background()

	tokens, user, err := svc.Register(ctx, authdomain.RegisterInput{
		FullName: "Test User", Phone: "09120000001",
		Password: "Secure1!", Role: common.RoleRequester,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.Equal(t, "Test User", user.FullName)

	loginTokens, loginUser, err := svc.Login(ctx, authdomain.LoginInput{
		Phone: "+989120000001", Password: "Secure1!",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, loginTokens.AccessToken)
	assert.Equal(t, user.ID, loginUser.ID)
}

func TestOTPVerifyFlow(t *testing.T) {
	svc := setupAuthTest(t)
	ctx := context.Background()

	cityID := uuid.New()
	nationalID := "0012345678"
	_, _, err := svc.Register(ctx, authdomain.RegisterInput{
		FullName: "OTP User", Phone: "09120000002",
		Password: "Secure1!", Role: common.RoleAgent,
		CityID: &cityID, NationalID: &nationalID,
	})
	require.NoError(t, err)

	err = svc.SendPhoneOTP(ctx, "+989120000002")
	require.NoError(t, err)
}

func TestPermissionsForRole(t *testing.T) {
	perms := common.PermissionsForRole(common.RoleRequester)
	assert.Contains(t, perms, "task:create")
	assert.NotContains(t, common.PermissionsForRole(common.RoleAgent), "task:create")
}
