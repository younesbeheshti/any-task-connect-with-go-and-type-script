package database

import (
	"context"
	"fmt"
	"time"

	"github.com/younesbeheshti/any-task-connect/backend/configs"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DB wraps the GORM database connection.
type DB struct {
	*gorm.DB
	log *zap.Logger
}

// Connect establishes a PostgreSQL connection with pooling and registers models.
func Connect(ctx context.Context, cfg configs.DatabaseConfig, log *zap.Logger, models ...any) (*DB, error) {
	gormCfg := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if len(models) > 0 {
		if err := RegisterModels(db, models...); err != nil {
			return nil, err
		}
	}

	log.Info("database connected",
		zap.String("host", cfg.Host),
		zap.String("database", cfg.Name),
	)

	return &DB{DB: db, log: log}, nil
}

// RegisterModels registers GORM models for development auto-migration.
// Production deployments should use golang-migrate SQL migrations instead.
func RegisterModels(db *gorm.DB, models ...any) error {
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("auto migrate model %T: %w", model, err)
		}
	}
	return nil
}

// Close closes the underlying connection pool.
func (d *DB) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// HealthCheck verifies database connectivity.
func (d *DB) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
