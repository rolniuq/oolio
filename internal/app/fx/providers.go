package fx

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"oolio/internal/app/handler"
	"oolio/internal/app/middleware"
	"oolio/internal/app/repository"
	"oolio/internal/app/router"
	"oolio/internal/app/services"
	"oolio/internal/app/worker"
	"oolio/internal/config"
	"oolio/internal/database"
)

// Config Module
var ConfigModule = fx.Module("config",
	fx.Provide(config.Load),
)

// Database Module
var DatabaseModule = fx.Module("database",
	fx.Provide(database.NewDatabase),
	fx.Provide(func(d *database.Database) *sql.DB { return d.DB }),
)

// Repository Module
var RepositoryModule = fx.Module("repository",
	fx.Provide(repository.NewProductRepository),
	fx.Provide(repository.NewOrderRepository),
	fx.Provide(repository.NewOrderQueueRepository),
)

// Service Module
var ServiceModule = fx.Module("service",
	fx.Provide(
		services.NewProductService,
		services.NewOrderService,
		services.NewOrderQueueService,
		NewRateLimiterService,
		NewCouponService,
	),
)

// Handler Module
var HandlerModule = fx.Module("handler",
	fx.Provide(
		handler.NewProductHandler,
		NewOrderHandler,
	),
)

// Middleware Module
var MiddlewareModule = fx.Module("middleware",
	fx.Provide(
		NewAuthMiddleware,
		NewErrorHandlerMiddleware,
		NewRateLimitMiddleware,
	),
)

// Worker Module
var WorkerModule = fx.Module("worker",
	fx.Provide(NewOrderWorker),
)

// Router Module
var RouterModule = fx.Module("router",
	fx.Provide(NewRouter),
)

// Custom provider for Coupon Service
func NewCouponService(cfg *config.Config) services.CouponService {
	return services.NewCouponService(cfg.Coupon.BaseURL)
}

// Custom provider for Auth Middleware
func NewAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return middleware.APIKeyAuth([]string{cfg.API.APIKey})
}

// Custom provider for Error Handler Middleware
func NewErrorHandlerMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.ErrorHandler(),
		middleware.PanicRecovery(),
	}
}

// Custom provider for Rate Limiter Service
func NewRateLimiterService(cfg *config.Config) services.RateLimiterService {
	return services.NewRateLimiterService(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
}

// Custom provider for Rate Limit Middleware
func NewRateLimitMiddleware(rateLimiter services.RateLimiterService) *middleware.RateLimitMiddleware {
	return middleware.NewRateLimitMiddleware(rateLimiter)
}

// Custom provider for OrderHandler
func NewOrderHandler(orderService services.OrderService, queueService services.OrderQueueService) *handler.OrderHandler {
	return handler.NewOrderHandler(orderService, queueService)
}

// Custom provider for Router
func NewRouter(
	productHandler *handler.ProductHandler,
	orderHandler *handler.OrderHandler,
	authMiddleware gin.HandlerFunc,
	errorMiddleware []gin.HandlerFunc,
	rateLimitMiddleware *middleware.RateLimitMiddleware,
) *gin.Engine {
	return router.SetupRouter(
		productHandler,
		orderHandler,
		authMiddleware,
		errorMiddleware,
		rateLimitMiddleware,
	)
}

// Custom provider for Order Worker
func NewOrderWorker(queueService services.OrderQueueService) *worker.OrderWorker {
	return worker.NewOrderWorker(queueService, 5*time.Second, 10) // Process every 5 seconds, batch size 10
}

// Application Modules
var AppModule = fx.Options(
	ConfigModule,
	DatabaseModule,
	RepositoryModule,
	ServiceModule,
	HandlerModule,
	MiddlewareModule,
	WorkerModule,
	RouterModule,
)
