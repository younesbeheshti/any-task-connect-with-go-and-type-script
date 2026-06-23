package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/younesbeheshti/any-task-connect/backend/configs"
	_ "github.com/younesbeheshti/any-task-connect/backend/docs"
	"github.com/younesbeheshti/any-task-connect/backend/internal/api"
	appinfra "github.com/younesbeheshti/any-task-connect/backend/internal/application/infra"
	appservice "github.com/younesbeheshti/any-task-connect/backend/internal/application/service"
	authinfra "github.com/younesbeheshti/any-task-connect/backend/internal/auth/infra"
	authservice "github.com/younesbeheshti/any-task-connect/backend/internal/auth/service"
	catinfra "github.com/younesbeheshti/any-task-connect/backend/internal/category/infra"
	catservice "github.com/younesbeheshti/any-task-connect/backend/internal/category/service"
	cityinfra "github.com/younesbeheshti/any-task-connect/backend/internal/city/infra"
	cityservice "github.com/younesbeheshti/any-task-connect/backend/internal/city/service"
	dashservice "github.com/younesbeheshti/any-task-connect/backend/internal/dashboard/service"
	"github.com/younesbeheshti/any-task-connect/backend/internal/health"
	taskinfra "github.com/younesbeheshti/any-task-connect/backend/internal/task/infra"
	taskservice "github.com/younesbeheshti/any-task-connect/backend/internal/task/service"
	userinfra "github.com/younesbeheshti/any-task-connect/backend/internal/user/infra"
	userservice "github.com/younesbeheshti/any-task-connect/backend/internal/user/service"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/database"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/jwt"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/logger"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/rabbitmq"
	redispkg "github.com/younesbeheshti/any-task-connect/backend/pkg/redis"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/server"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/validator"
	"go.uber.org/zap"
)

// App holds wired dependencies.
type App struct {
	Config *configs.Config
	Logger *zap.Logger
	DB     *database.DB
	Redis  *redispkg.Client
	Server *server.Server
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

	rdb, err := redispkg.Connect(ctx, cfg.Redis, log)
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
	_ = mq.SetupEventQueues(ctx)

	jwtService := jwt.NewService(cfg.JWT)
	cache := redispkg.NewCache(rdb.Client)
	refreshTTL := cfg.JWT.RefreshTokenDuration()
	accessTTL := cfg.JWT.AccessTokenDuration()

	// Repositories.
	userRepo := userinfra.NewGormRepository(db.DB)
	authRepo := authinfra.NewGormRepository(db.DB)
	sessions := authinfra.NewSessionStore(cache, refreshTTL)
	otpStore := authinfra.NewOTPStore(cache)
	lockout := authinfra.NewLockoutStore(cache)

	catRepo := catinfra.NewGormRepository(db.DB)
	cityRepo := cityinfra.NewGormRepository(db.DB)
	taskRepo := taskinfra.NewGormRepository(db.DB)
	appRepo := appinfra.NewGormRepository(db.DB)

	// Services.
	authSvc := authservice.NewAuthService(
		userRepo, authRepo, jwtService, cache,
		sessions, otpStore, lockout, accessTTL, refreshTTL,
	)
	userSvc := userservice.NewUserService(userRepo)
	catSvc := catservice.NewCategoryService(catRepo)
	citySvc := cityservice.NewCityService(cityRepo)
	taskSvc := taskservice.NewTaskService(taskRepo, rabbitmq.NewPublisher(mq))
	appSvc := appservice.NewApplicationService(appRepo, taskRepo)
	dashSvc := dashservice.NewDashboardService(db.DB)

	v := validator.New()
	handlers := api.NewHandlers(authSvc, userSvc, catSvc, citySvc, taskSvc, appSvc, dashSvc, v)
	healthHandler := health.NewHandler(db, rdb, mq)

	srv := server.New(server.Dependencies{
		Config: cfg,
		Logger: log,
		RegisterRoutes: func(r *gin.Engine) {
			healthHandler.RegisterRoutes(r)
			api.RegisterRoutes(r, handlers, authSvc)
			r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		},
	})

	app := &App{Config: cfg, Logger: log, DB: db, Redis: rdb, Server: srv}
	return app.serve(mq)
}

func (a *App) serve(mq *rabbitmq.Client) error {
	errCh := make(chan error, 1)
	go func() { errCh <- a.Server.Start() }()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		a.shutdown(mq)
		return err
	case <-quit:
		a.Logger.Info("shutdown signal received")
		a.shutdown(mq)
		return nil
	}
}

func (a *App) shutdown(mq *rabbitmq.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = a.Server.Shutdown(ctx)
	if mq != nil {
		_ = mq.Close()
	}
	if a.Redis != nil {
		_ = a.Redis.Close()
	}
	if a.DB != nil {
		_ = a.DB.Close()
	}
	a.Logger.Info("application stopped")
}
