package api

import (
	"github.com/gin-gonic/gin"
	apphandler "github.com/younesbeheshti/any-task-connect/backend/internal/application/handler"
	appservice "github.com/younesbeheshti/any-task-connect/backend/internal/application/service"
	authhandler "github.com/younesbeheshti/any-task-connect/backend/internal/auth/handler"
	authservice "github.com/younesbeheshti/any-task-connect/backend/internal/auth/service"
	cathandler "github.com/younesbeheshti/any-task-connect/backend/internal/category/handler"
	catservice "github.com/younesbeheshti/any-task-connect/backend/internal/category/service"
	cityhandler "github.com/younesbeheshti/any-task-connect/backend/internal/city/handler"
	cityservice "github.com/younesbeheshti/any-task-connect/backend/internal/city/service"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common/middleware"
	dashhandler "github.com/younesbeheshti/any-task-connect/backend/internal/dashboard/handler"
	dashservice "github.com/younesbeheshti/any-task-connect/backend/internal/dashboard/service"
	paymenthandler "github.com/younesbeheshti/any-task-connect/backend/internal/payment/handler"
	revenuehandler "github.com/younesbeheshti/any-task-connect/backend/internal/revenue/handler"
	taskhandler "github.com/younesbeheshti/any-task-connect/backend/internal/task/handler"
	taskservice "github.com/younesbeheshti/any-task-connect/backend/internal/task/service"
	userhandler "github.com/younesbeheshti/any-task-connect/backend/internal/user/handler"
	userservice "github.com/younesbeheshti/any-task-connect/backend/internal/user/service"
	wallethandler "github.com/younesbeheshti/any-task-connect/backend/internal/wallet/handler"
	withdrawhandler "github.com/younesbeheshti/any-task-connect/backend/internal/withdraw/handler"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/validator"
)

// Handlers groups HTTP handlers for API routes.
type Handlers struct {
	Auth        *authhandler.Handler
	User        *userhandler.Handler
	Category    *cathandler.Handler
	City        *cityhandler.Handler
	Task        *taskhandler.Handler
	Application *apphandler.Handler
	Dashboard   *dashhandler.Handler
	Wallet      *wallethandler.Handler
	Payment     *paymenthandler.Handler
	Revenue     *revenuehandler.Handler
	Withdraw    *withdrawhandler.Handler
}

// NewHandlers creates API handlers.
func NewHandlers(
	auth *authservice.AuthService,
	users *userservice.UserService,
	categories *catservice.CategoryService,
	cities *cityservice.CityService,
	tasks *taskservice.TaskService,
	applications *appservice.ApplicationService,
	dashboard dashservice.Service,
	wallet *wallethandler.Handler,
	payment *paymenthandler.Handler,
	revenue *revenuehandler.Handler,
	withdraw *withdrawhandler.Handler,
	v *validator.Validator,
) *Handlers {
	return &Handlers{
		Auth:        authhandler.NewHandler(auth, v),
		User:        userhandler.NewHandler(users, v),
		Category:    cathandler.NewHandler(categories),
		City:        cityhandler.NewHandler(cities),
		Task:        taskhandler.NewHandler(tasks),
		Application: apphandler.NewHandler(applications),
		Dashboard:   dashhandler.NewHandler(dashboard),
		Wallet:      wallet,
		Payment:     payment,
		Revenue:     revenue,
		Withdraw:    withdraw,
	}
}

// RegisterRoutes mounts /v1 API routes per front/docs/api-contracts.md.
func RegisterRoutes(r *gin.Engine, h *Handlers, auth *authservice.AuthService) {
	v1 := r.Group("/v1")

	// Auth routes (public).
	h.Auth.RegisterRoutes(v1)

	// Public read routes.
	v1.GET("/categories", h.Category.List)
	v1.GET("/categories/:id", h.Category.GetByID)
	v1.GET("/cities", h.City.List)
	v1.GET("/cities/:id", h.City.GetByID)
	v1.GET("/tasks", h.Task.List)
	v1.GET("/tasks/:id", h.Task.GetByPublicID)

	// Protected routes (auth required).
	protected := v1.Group("")
	protected.Use(middleware.AuthJWT(auth))

	// User profile.
	h.User.RegisterRoutes(protected)

	// Task protected actions.
	protected.POST("/tasks", h.Task.Create)
	protected.PATCH("/tasks/:id", h.Task.Update)
	protected.DELETE("/tasks/:id", h.Task.Delete)
	protected.GET("/tasks/:id/timeline", h.Task.GetTimeline)
	protected.POST("/tasks/:id/cancel", h.Task.Cancel)
	protected.POST("/tasks/:id/start", h.Task.Start)
	protected.POST("/tasks/:id/complete", h.Task.Complete)
	protected.POST("/tasks/:id/verify", h.Task.Verify)

	// Application routes.
	protected.POST("/tasks/:id/applications", h.Application.Submit)
	protected.GET("/tasks/:id/applications", h.Application.ListByTask)
	protected.GET("/applications/:id", h.Application.GetByID)
	protected.GET("/me/applications", h.Application.ListMyApplications)
	protected.POST("/applications/:id/accept", h.Application.Accept)
	protected.POST("/applications/:id/reject", h.Application.Reject)
	protected.POST("/applications/:id/withdraw", h.Application.Withdraw)

	// Dashboard user stats.
	protected.GET("/dashboard/stats", h.Dashboard.GetUserStats)

	// Wallet routes.
	protected.GET("/wallet", h.Wallet.GetWallet)
	protected.GET("/wallet/history", h.Wallet.GetHistory)
	protected.GET("/wallet/statistics", h.Wallet.GetStatistics)
	protected.POST("/wallet/topup", h.Wallet.TopUp)
	protected.POST("/wallet/withdraw", h.Withdraw.Create)

	// Transaction routes.
	protected.GET("/transactions", h.Payment.ListTransactions)
	protected.GET("/transactions/:id", h.Payment.GetTransaction)

	// Withdraw routes.
	protected.GET("/withdraws", h.Withdraw.ListMine)
	protected.GET("/withdraws/:id", h.Withdraw.GetByID)

	// Admin routes (auth + admin permission required).
	admin := protected.Group("/admin")
	admin.Use(middleware.PermissionRequired(common.PermAdminDashboard))
	admin.POST("/categories", h.Category.Create)
	admin.PATCH("/categories/:id", h.Category.Update)
	admin.DELETE("/categories/:id", h.Category.Delete)
	admin.POST("/cities", h.City.Create)
	admin.PATCH("/cities/:id", h.City.Update)
	admin.DELETE("/cities/:id", h.City.Delete)
	admin.GET("/dashboard/admin-stats", h.Dashboard.GetAdminStats)

	// Admin financial routes.
	admin.GET("/revenue", h.Revenue.GetRevenue)
	admin.GET("/revenue/statistics", h.Revenue.GetStatistics)
	admin.GET("/revenue/daily", h.Revenue.GetDaily)
	admin.GET("/revenue/monthly", h.Revenue.GetMonthly)
	admin.GET("/financial/overview", h.Wallet.GetAdminOverview)
	admin.GET("/financial/revenue", h.Revenue.GetRevenue)
	admin.GET("/financial/transactions", h.Payment.ListAllTransactions)
	admin.GET("/financial/withdraws", h.Withdraw.ListAll)
	admin.POST("/withdraws/:id/approve", h.Withdraw.Approve)
	admin.POST("/withdraws/:id/reject", h.Withdraw.Reject)
}
