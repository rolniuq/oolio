package router

import (
	"time"

	"oolio/internal/app/handler"
	"oolio/internal/app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	productHandler *handler.ProductHandler,
	orderHandler *handler.OrderHandler,
	authMiddleware gin.HandlerFunc,
	errorMiddleware []gin.HandlerFunc,
	rateLimitMiddleware *middleware.RateLimitMiddleware,
) *gin.Engine {
	r := gin.Default()

	// Apply global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Apply CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Apply error handling middleware
	for _, mw := range errorMiddleware {
		r.Use(mw)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Service is running",
		})
	})

	// Product routes (no authentication required)
	v1 := r.Group("/api/v1")
	{
		// Product endpoints (authentication + rate limiting)
		products := v1.Group("/product").Use(authMiddleware, rateLimitMiddleware.RateLimit(100, time.Minute))
		{
			products.GET("/", productHandler.ListProducts)
			products.GET("/:productId", productHandler.GetProduct)
		}

		// Also support direct access without trailing slash to avoid redirect
		v1.GET("/product", authMiddleware, rateLimitMiddleware.RateLimit(100, time.Minute), productHandler.ListProducts)

		// Order endpoints (authentication + rate limiting)
		orders := v1.Group("/order").Use(authMiddleware, rateLimitMiddleware.RateLimit(50, time.Minute))
		{
			orders.POST("", orderHandler.PlaceOrder)
			orders.GET("", orderHandler.ListOrders)
			orders.GET("/:orderId", orderHandler.GetOrder)
		}

		// Queue status endpoint (authentication + rate limiting)
		v1.GET("/queue/status", authMiddleware, rateLimitMiddleware.RateLimit(30, time.Minute), orderHandler.GetQueueStatus)
	}

	return r
}
