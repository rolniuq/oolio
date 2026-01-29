package middleware

import (
	"net/http"
	"slices"

	"oolio/internal/app/models"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth(validKeys []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			// Fallback to lowercase for compatibility
			apiKey = c.GetHeader("api_key")
		}
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Code:    http.StatusUnauthorized,
				Type:    "error",
				Message: "API key is required",
			})
			c.Abort()
			return
		}

		// Validate API key
		isValid := slices.Contains(validKeys, apiKey)

		if !isValid {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Code:    http.StatusUnauthorized,
				Type:    "error",
				Message: "Invalid API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// For now, this is a placeholder for future permission-based access control
		// We can extend this to check user roles and permissions from a database or JWT claims
		c.Next()
	}
}

// Optional: Add CORS middleware for development
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, api_key, X-API-Key")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
