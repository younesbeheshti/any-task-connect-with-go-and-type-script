package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/younesbeheshti/any-task-connect/backend/configs"
	_ "github.com/younesbeheshti/any-task-connect/backend/docs"
	"github.com/younesbeheshti/any-task-connect/backend/internal/health"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/database"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/jwt"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/logger"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/rabbitmq"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/redis"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/server"
	"go.uber.org/zap"
)

// App holds wired infrastructure dependencies.
type App struct {
	Config   *configs.Config
	Logger   *zap.Logger
	DB       *database.DB
	Redis    *redis.Client
	RabbitMQ *rabbitmq.Client
	JWT      *jwt.Service
	Server   *server.Server
}

// Run initializes infrastructure and starts the HTTP server.
func Run() error {
	cfg, err := configs.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log, err := logger.Init(logger.Config{
		Environment: cfg.App.Environment,
		AppName:     cfg.App.Name,
	})
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	defer logger.Sync()

	logger.Startup("application",
		zap.String("environment", cfg.App.Environment),
		zap.String("port", cfg.App.Port),
	)

	ctx := context.Background()

	db, err := database.Connect(ctx, cfg.Database, log)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	rdb, err := redis.Connect(ctx, cfg.Redis, log)
	if err != nil {
		_ = db.Close()
		return fmt.Errorf("connect redis: %w", err)
	}

	mq, err := rabbitmq.Connect(ctx, cfg.RabbitMQ, log)
	if err != nil {
		_ = rdb.Close()
		_ = db.Close()
		return fmt.Errorf("connect rabbitmq: %w", err)
	}

	if err := mq.SetupEventQueues(ctx); err != nil {
		log.Warn("setup event queues", zap.Error(err))
	}

	jwtService := jwt.NewService(cfg.JWT)
	healthHandler := health.NewHandler(db, rdb, mq)

	srv := server.New(server.Dependencies{
		Config: cfg,
		Logger: log,
		RegisterRoutes: func(r *gin.Engine) {
			healthHandler.RegisterRoutes(r)
			r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		},
	})

	app := &App{
		Config:   cfg,
		Logger:   log,
		DB:       db,
		Redis:    rdb,
		RabbitMQ: mq,
		JWT:      jwtService,
		Server:   srv,
	}

	return app.serve()
}

func (a *App) serve() error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- a.Server.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		a.shutdown()
		return err
	case <-quit:
		a.Logger.Info("shutdown signal received")
		a.shutdown()
		return nil
	}
}

func (a *App) shutdown() {
	ctx := context.Background()
	if err := a.Server.Shutdown(ctx); err != nil {
		a.Logger.Error("server shutdown", zap.Error(err))
	}
	if a.RabbitMQ != nil {
		_ = a.RabbitMQ.Close()
	}
	if a.Redis != nil {
		_ = a.Redis.Close()
	}
	if a.DB != nil {
		_ = a.DB.Close()
	}
	a.Logger.Info("application stopped")
}
