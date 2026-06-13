package health

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Checker verifies dependency health.
type Checker interface {
	HealthCheck(ctx context.Context) error
}

// Handler serves health endpoints.
type Handler struct {
	db       Checker
	redis    Checker
	rabbitmq Checker
}

// NewHandler creates a health handler.
func NewHandler(db, redis, rabbitmq Checker) *Handler {
	return &Handler{db: db, redis: redis, rabbitmq: rabbitmq}
}

// RegisterRoutes mounts health endpoints on the router.
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.Health)
	r.GET("/ready", h.Ready)
}

type healthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

// Health godoc
// @Summary      Health check
// @Description  Returns the health status of the application and its dependencies
// @Tags         health
// @Produce      json
// @Success      200  {object}  healthResponse
// @Failure      503  {object}  healthResponse
// @Router       /health [get]
func (h *Handler) Health(c *gin.Context) {
	h.respond(c)
}

// Ready godoc
// @Summary      Readiness check
// @Description  Returns readiness status for orchestrators
// @Tags         health
// @Produce      json
// @Success      200  {object}  healthResponse
// @Failure      503  {object}  healthResponse
// @Router       /ready [get]
func (h *Handler) Ready(c *gin.Context) {
	h.respond(c)
}

func (h *Handler) respond(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	services := make(map[string]string)

	var wg sync.WaitGroup
	var mu sync.Mutex

	checks := []struct {
		name    string
		checker Checker
	}{
		{"database", h.db},
		{"redis", h.redis},
		{"rabbitmq", h.rabbitmq},
	}

	for _, check := range checks {
		wg.Add(1)
		go func(name string, checker Checker) {
			defer wg.Done()
			status := "UP"
			if checker == nil {
				status = "DOWN"
			} else if err := checker.HealthCheck(ctx); err != nil {
				status = "DOWN"
			}
			mu.Lock()
			services[name] = status
			mu.Unlock()
		}(check.name, check.checker)
	}
	wg.Wait()

	overall := "UP"
	for _, status := range services {
		if status != "UP" {
			overall = "DOWN"
			break
		}
	}

	code := http.StatusOK
	if overall == "DOWN" {
		code = http.StatusServiceUnavailable
	}

	c.JSON(code, healthResponse{
		Status:   overall,
		Services: services,
	})
}
