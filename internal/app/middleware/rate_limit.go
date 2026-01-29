package middleware

import (
	"net/http"
	"strconv"
	"time"

	"oolio/internal/app/services"

	"github.com/gin-gonic/gin"
)

type RateLimitMiddleware struct {
	rateLimiter services.RateLimiterService
}

func NewRateLimitMiddleware(rateLimiter services.RateLimiterService) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		rateLimiter: rateLimiter,
	}
}

// RateLimit creates a middleware that limits requests based on the provided parameters
func (m *RateLimitMiddleware) RateLimit(requestsPerMinute int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If rate limiter is nil (e.g., in tests), skip rate limiting
		if m.rateLimiter == nil {
			c.Next()
			return
		}

		// Use IP address as the key for rate limiting
		key := "rate_limit:" + c.ClientIP()

		// Check if request is allowed
		allowed, err := m.rateLimiter.AllowRequest(c.Request.Context(), key, requestsPerMinute, window)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Rate limiter error",
			})
			c.Abort()
			return
		}

		if !allowed {
			// Get remaining tokens for response headers
			remaining, _ := m.rateLimiter.GetRemainingTokens(c.Request.Context(), key, requestsPerMinute)

			c.Header("X-RateLimit-Limit", strconv.Itoa(requestsPerMinute))
			c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		// Add rate limit headers for successful requests
		remaining, _ := m.rateLimiter.GetRemainingTokens(c.Request.Context(), key, requestsPerMinute)
		c.Header("X-RateLimit-Limit", strconv.Itoa(requestsPerMinute))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))

		c.Next()
	}
}

// RateLimitByUser creates a middleware that limits requests per user (requires user ID in context)
func (m *RateLimitMiddleware) RateLimitByUser(requestsPerMinute int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If rate limiter is nil (e.g., in tests), skip rate limiting
		if m.rateLimiter == nil {
			c.Next()
			return
		}

		// Try to get user ID from context or use IP as fallback
		userID, exists := c.Get("user_id")
		if !exists {
			userID = c.ClientIP()
		}

		key := "rate_limit:user:" + userID.(string)

		allowed, err := m.rateLimiter.AllowRequest(c.Request.Context(), key, requestsPerMinute, window)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Rate limiter error",
			})
			c.Abort()
			return
		}

		if !allowed {
			remaining, _ := m.rateLimiter.GetRemainingTokens(c.Request.Context(), key, requestsPerMinute)

			c.Header("X-RateLimit-Limit", strconv.Itoa(requestsPerMinute))
			c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		remaining, _ := m.rateLimiter.GetRemainingTokens(c.Request.Context(), key, requestsPerMinute)
		c.Header("X-RateLimit-Limit", strconv.Itoa(requestsPerMinute))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))

		c.Next()
	}
}
