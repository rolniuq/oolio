package integration

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"oolio/internal/app/handler"
	"oolio/internal/app/middleware"
	"oolio/internal/app/router"
)

func TestIntegration_Routing_Products(t *testing.T) {
	// Create mock services
	mockProductService := &MockProductService{}
	mockOrderService := &MockOrderService{}
	mockQueueService := &MockOrderQueueService{}

	// Create mock handlers
	mockProductHandler := handler.NewProductHandler(mockProductService)
	mockOrderHandler := handler.NewOrderHandler(mockOrderService, mockQueueService)

	// Create auth middleware that allows all requests
	authMiddleware := middleware.APIKeyAuth([]string{"any-key"})

	// Create mock rate limit middleware
	mockRateLimiter := &MockRateLimiterService{}
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(mockRateLimiter)

	// Setup router
	router := router.SetupRouter(mockProductHandler, mockOrderHandler, authMiddleware, []gin.HandlerFunc{}, rateLimitMiddleware)

	// Test GET /api/v1/product
	req, _ := http.NewRequest("GET", "/api/v1/product", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Routing_Orders(t *testing.T) {
	// Create mock services
	mockOrderService := &MockOrderService{}
	mockQueueService := &MockOrderQueueService{}

	// Create simple mock handler
	mockHandler := handler.NewOrderHandler(mockOrderService, mockQueueService)

	// Create auth middleware that requires specific key
	authMiddleware := middleware.APIKeyAuth([]string{"test-api-key"})

	// Create mock rate limit middleware
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(&MockRateLimiterService{})

	// Setup router
	router := router.SetupRouter(nil, mockHandler, authMiddleware, []gin.HandlerFunc{}, rateLimitMiddleware)

	// Test POST /api/v1/order with valid API key
	jsonBody := []byte(`{"items": [{"productId": "test-1", "quantity": 2}]}`)
	req, _ := http.NewRequest("POST", "/api/v1/order", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_key", "test-api-key")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should get 202 because order is queued successfully
	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestIntegration_Routing_Orders_Unauthorized(t *testing.T) {
	// Create mock services
	mockOrderService := &MockOrderService{}
	mockQueueService := &MockOrderQueueService{}

	// Create simple mock handler
	mockHandler := handler.NewOrderHandler(mockOrderService, mockQueueService)

	// Create auth middleware that requires specific key
	authMiddleware := middleware.APIKeyAuth([]string{"test-api-key"})

	// Create mock rate limit middleware
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(&MockRateLimiterService{})

	// Setup router
	router := router.SetupRouter(nil, mockHandler, authMiddleware, []gin.HandlerFunc{}, rateLimitMiddleware)

	// Test POST /api/v1/order without API key
	jsonBody := []byte(`{"items": [{"productId": "test-1", "quantity": 2}]}`)
	req, _ := http.NewRequest("POST", "/api/v1/order", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestIntegration_Routing_Orders_InvalidAPIKey(t *testing.T) {
	// Create mock services
	mockOrderService := &MockOrderService{}
	mockQueueService := &MockOrderQueueService{}

	// Create simple mock handler
	mockHandler := handler.NewOrderHandler(mockOrderService, mockQueueService)

	// Create auth middleware that requires specific key
	authMiddleware := middleware.APIKeyAuth([]string{"test-api-key"})

	// Create mock rate limit middleware
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(&MockRateLimiterService{})

	// Setup router
	router := router.SetupRouter(nil, mockHandler, authMiddleware, []gin.HandlerFunc{}, rateLimitMiddleware)

	// Test POST /api/v1/order with invalid API key
	jsonBody := []byte(`{"items": [{"productId": "test-1", "quantity": 2}]}`)
	req, _ := http.NewRequest("POST", "/api/v1/order", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_key", "invalid-key")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestIntegration_Routing_HealthCheck(t *testing.T) {
	// Create mock services for order handler
	mockOrderService := &MockOrderService{}
	mockQueueService := &MockOrderQueueService{}
	mockOrderHandler := handler.NewOrderHandler(mockOrderService, mockQueueService)

	// Create auth middleware
	authMiddleware := middleware.APIKeyAuth([]string{"any-key"})

	// Create mock rate limit middleware
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(&MockRateLimiterService{})

	// Setup router
	router := router.SetupRouter(nil, mockOrderHandler, authMiddleware, []gin.HandlerFunc{}, rateLimitMiddleware)

	// Test GET /health
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Routing_Products_WithAuth(t *testing.T) {
	// Create mock services
	mockProductService := &MockProductService{}
	mockOrderService := &MockOrderService{}
	mockQueueService := &MockOrderQueueService{}

	// Create mock handlers
	mockProductHandler := handler.NewProductHandler(mockProductService)
	mockOrderHandler := handler.NewOrderHandler(mockOrderService, mockQueueService)

	// Create auth middleware that requires specific key
	authMiddleware := middleware.APIKeyAuth([]string{"test-api-key"})

	// Create mock rate limit middleware
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(&MockRateLimiterService{})

	// Setup router
	router := router.SetupRouter(mockProductHandler, mockOrderHandler, authMiddleware, []gin.HandlerFunc{}, rateLimitMiddleware)

	// Test GET /api/v1/product (should work even with auth)
	req, _ := http.NewRequest("GET", "/api/v1/product", nil)
	req.Header.Set("api_key", "test-api-key")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
