package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/younesbeheshti/any-task-connect/backend/configs"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/database"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/jwt"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/validator"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// TestLogger returns a no-op zap logger for tests.
func TestLogger(t *testing.T) *zap.Logger {
	t.Helper()
	log, err := zap.NewDevelopment()
	require.NoError(t, err)
	t.Cleanup(func() { _ = log.Sync() })
	return log
}

// TestConfig returns minimal test configuration.
func TestConfig() *configs.Config {
	return &configs.Config{
		App: configs.AppConfig{
			Name:        "taskbridge-test",
			Environment: "development",
			Port:        "8080",
		},
		Database: configs.DatabaseConfig{
			MaxOpenConns:    5,
			MaxIdleConns:    2,
			ConnMaxLifetime: time.Minute,
		},
		JWT: configs.JWTConfig{
			Secret:     "test-secret-key-minimum-32-characters",
			AccessTTL:  900,
			RefreshTTL: 604800,
		},
	}
}

// MockDB opens an in-memory SQLite database for unit tests.
func MockDB(t *testing.T) *database.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	require.NoError(t, err)
	return &database.DB{DB: db}
}

// MockRedis starts an in-process Redis server.
func MockRedis(t *testing.T) (*redis.Client, func()) {
	t.Helper()
	s, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	cleanup := func() {
		_ = client.Close()
		s.Close()
	}
	t.Cleanup(cleanup)
	return client, cleanup
}

// MockRabbitMQ is a no-op publisher for tests without a broker.
type MockRabbitMQ struct{}

func (m *MockRabbitMQ) HealthCheck(_ context.Context) error { return nil }

// NewJWTService creates a JWT service for tests.
func NewJWTService(cfg *configs.Config) *jwt.Service {
	return jwt.NewService(cfg.JWT)
}

// NewValidator creates a validator for tests.
func NewValidator() *ValidatorWrapper {
	return &ValidatorWrapper{Validator: validator.New()}
}

// ValidatorWrapper exposes the test validator.
type ValidatorWrapper struct {
	*validator.Validator
}
