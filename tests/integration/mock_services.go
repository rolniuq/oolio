package integration

import (
	"context"
	"time"

	"oolio/internal/app/models"
)

// MockOrderService implements OrderService for testing
type MockOrderService struct{}

func (m *MockOrderService) CreateOrder(ctx context.Context, orderReq *models.OrderReq) (*models.Order, error) {
	return &models.Order{
		ID:        "test-order-id",
		Total:     100.0,
		Discounts: 0.0,
		Items:     orderReq.Items,
	}, nil
}

func (m *MockOrderService) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	return &models.Order{
		ID:        id,
		Total:     100.0,
		Discounts: 0.0,
		Items:     []models.OrderItem{},
	}, nil
}

// MockOrderQueueService implements OrderQueueService for testing
type MockOrderQueueService struct{}

func (m *MockOrderQueueService) AddOrderToQueue(ctx context.Context, orderReq *models.OrderReq) (*models.OrderQueueItem, error) {
	return &models.OrderQueueItem{
		ID:       "test-queue-id",
		OrderReq: *orderReq,
		Status:   "pending",
	}, nil
}

func (m *MockOrderQueueService) GetCompletedOrders(ctx context.Context) ([]*models.OrderQueueItem, error) {
	return []*models.OrderQueueItem{}, nil
}

func (m *MockOrderQueueService) ProcessBatch(ctx context.Context, batchSize int) (*models.BatchProcessResult, error) {
	return &models.BatchProcessResult{
		Processed: 0,
		Failed:    0,
		Errors:    []string{},
		Items:     []models.OrderQueueItem{},
	}, nil
}

func (m *MockOrderQueueService) GetQueueStatus(ctx context.Context) (map[string]int, error) {
	return map[string]int{
		"pending":    0,
		"processing": 0,
		"completed":  0,
		"failed":     0,
	}, nil
}

func (m *MockOrderQueueService) StartWorker(ctx context.Context, interval time.Duration, batchSize int) {
	// Mock implementation does nothing
}

func (m *MockOrderQueueService) GetOrderFromQueue(ctx context.Context, itemID string) (*models.OrderQueueItem, error) {
	return &models.OrderQueueItem{
		ID:       itemID,
		OrderReq: models.OrderReq{},
		Status:   "pending",
	}, nil
}

// MockRateLimiterService implements RateLimiterService for testing
type MockRateLimiterService struct{}

func (m *MockRateLimiterService) AllowRequest(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	// Always allow requests in tests
	return true, nil
}

func (m *MockRateLimiterService) IsAllowed(ctx context.Context, key string) (bool, error) {
	return true, nil
}

func (m *MockRateLimiterService) GetRemainingTokens(ctx context.Context, key string, limit int) (int, error) {
	return limit, nil
}

func (m *MockRateLimiterService) ResetKey(ctx context.Context, key string) error {
	return nil
}

// MockProductService implements ProductService for testing
type MockProductService struct{}

func (m *MockProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	return []models.Product{
		{
			ID:       "test-product-1",
			Name:     "Test Product 1",
			Price:    10.99,
			Category: "Waffle",
			Image: models.Image{
				Thumbnail: "https://example.com/thumb1.jpg",
				Mobile:    "https://example.com/mobile1.jpg",
				Tablet:    "https://example.com/tablet1.jpg",
				Desktop:   "https://example.com/desktop1.jpg",
			},
		},
		{
			ID:       "test-product-2",
			Name:     "Test Product 2",
			Price:    15.99,
			Category: "Waffle",
			Image: models.Image{
				Thumbnail: "https://example.com/thumb2.jpg",
				Mobile:    "https://example.com/mobile2.jpg",
				Tablet:    "https://example.com/tablet2.jpg",
				Desktop:   "https://example.com/desktop2.jpg",
			},
		},
	}, nil
}

func (m *MockProductService) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	if id == "test-product-1" {
		return &models.Product{
			ID:       "test-product-1",
			Name:     "Test Product 1",
			Price:    10.99,
			Category: "Waffle",
			Image: models.Image{
				Thumbnail: "https://example.com/thumb1.jpg",
				Mobile:    "https://example.com/mobile1.jpg",
				Tablet:    "https://example.com/tablet1.jpg",
				Desktop:   "https://example.com/desktop1.jpg",
			},
		}, nil
	}
	return nil, nil
}

func (m *MockProductService) CreateProduct(ctx context.Context, product *models.Product) error {
	// Mock implementation does nothing
	return nil
}

func (m *MockProductService) UpdateProduct(ctx context.Context, product *models.Product) error {
	// Mock implementation does nothing
	return nil
}

func (m *MockProductService) DeleteProduct(ctx context.Context, id string) error {
	// Mock implementation does nothing
	return nil
}

func (m *MockProductService) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	if id == "test-product-1" {
		return &models.Product{
			ID:       "test-product-1",
			Name:     "Test Product 1",
			Price:    10.99,
			Category: "Waffle",
			Image: models.Image{
				Thumbnail: "https://example.com/thumb1.jpg",
				Mobile:    "https://example.com/mobile1.jpg",
				Tablet:    "https://example.com/tablet1.jpg",
				Desktop:   "https://example.com/desktop1.jpg",
			},
		}, nil
	}
	return nil, nil
}
