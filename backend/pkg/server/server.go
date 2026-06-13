package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/younesbeheshti/any-task-connect/backend/configs"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	"go.uber.org/zap"
)

// Server wraps the HTTP server and Gin engine.
type Server struct {
	cfg    configs.ServerConfig
	port   string
	engine *gin.Engine
	http   *http.Server
	log    *zap.Logger
}

// Dependencies holds middleware and route registration callbacks.
type Dependencies struct {
	Config  *configs.Config
	Logger  *zap.Logger
	RegisterRoutes func(r *gin.Engine)
}

// New creates a configured Gin server.
func New(deps Dependencies) *Server {
	if deps.Config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(
		middleware.Recovery(deps.Logger),
		middleware.RequestID(),
		middleware.Logging(deps.Logger),
		middleware.CORS(deps.Config.CORS),
		middleware.SecurityHeaders(),
		middleware.Timeout(deps.Config.Server.ReadTimeout),
	)

	if deps.RegisterRoutes != nil {
		deps.RegisterRoutes(engine)
	}

	addr := fmt.Sprintf(":%s", deps.Config.App.Port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  deps.Config.Server.ReadTimeout,
		WriteTimeout: deps.Config.Server.WriteTimeout,
		IdleTimeout:  deps.Config.Server.IdleTimeout,
	}

	return &Server{
		cfg:    deps.Config.Server,
		port:   deps.Config.App.Port,
		engine: engine,
		http:   httpServer,
		log:    deps.Logger,
	}
}

// Engine returns the Gin engine for route registration.
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

// Start begins listening for HTTP requests.
func (s *Server) Start() error {
	s.log.Info("starting http server", zap.String("port", s.port))
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	s.log.Info("shutting down http server")
	return s.http.Shutdown(ctx)
}
