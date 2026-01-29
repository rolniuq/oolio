package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"oolio/internal/app/models"
	"oolio/internal/app/repository"
)

// Mock order repository for testing
type mockOrderRepository struct {
	orders []models.Order
}

func NewMockOrderRepository() repository.OrderRepository {
	return &mockOrderRepository{
		orders: []models.Order{
			{
				ID:        "test-order-1",
				Total:     25.99,
				Discounts: 0.0,
				Items: []models.OrderItem{
					{
						ProductID: "test-product-1",
						Quantity:  2,
						Price:     10.99,
					},
				},
				Products: []models.Product{
					{
						ID:       "test-product-1",
						Name:     "Test Product 1",
						Price:    10.99,
						Category: "Waffle",
						Image: models.Image{
							Thumbnail: "http://example.com/thumb.jpg",
							Mobile:    "http://example.com/mobile.jpg",
							Tablet:    "http://example.com/tablet.jpg",
							Desktop:   "http://example.com/desktop.jpg",
						},
					},
				},
			},
		},
	}
}

func (r *mockOrderRepository) Find(ctx context.Context) ([]models.Order, error) {
	return r.orders, nil
}

func (r *mockOrderRepository) FindOne(ctx context.Context, id string) (*models.Order, error) {
	for _, order := range r.orders {
		if order.ID == id {
			return &order, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *mockOrderRepository) Create(ctx context.Context, order *models.Order) error {
	order.ID = uuid.New().String()
	r.orders = append(r.orders, *order)
	return nil
}

func (r *mockOrderRepository) Update(ctx context.Context, order *models.Order) error {
	for i, o := range r.orders {
		if o.ID == order.ID {
			r.orders[i] = *order
			return nil
		}
	}
	return sql.ErrNoRows
}

func (r *mockOrderRepository) Delete(ctx context.Context, id string) error {
	// Not implemented as per business requirements
	return sql.ErrNoRows
}

func (r *mockOrderRepository) CreateOrderItems(ctx context.Context, orderID string, items []models.OrderItem) error {
	// Mock implementation - just returns nil for success
	return nil
}

func (r *mockOrderRepository) GetOrderItems(ctx context.Context, orderID string) ([]models.OrderItem, error) {
	for _, order := range r.orders {
		if order.ID == orderID {
			return order.Items, nil
		}
	}
	return nil, sql.ErrNoRows
}

func TestOrderRepository_FindOne(t *testing.T) {
	repo := NewMockOrderRepository()
	ctx := context.Background()

	// Test existing order
	order, err := repo.FindOne(ctx, "test-order-1")
	assert.NoError(t, err)
	require.NotNil(t, order)
	assert.Equal(t, "test-order-1", order.ID)
	assert.Equal(t, 25.99, order.Total)
	assert.Len(t, order.Items, 1)

	// Test non-existing order
	order, err = repo.FindOne(ctx, "non-existing")
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestOrderRepository_Create(t *testing.T) {
	repo := NewMockOrderRepository()
	ctx := context.Background()

	newOrder := &models.Order{
		Total:     50.99,
		Discounts: 5.00,
		Items: []models.OrderItem{
			{
				ProductID: "test-product-2",
				Quantity:  3,
				Price:     15.99,
			},
		},
	}

	err := repo.Create(ctx, newOrder)
	assert.NoError(t, err)
	assert.NotEmpty(t, newOrder.ID)

	// Verify order was added
	orders, err := repo.Find(ctx)
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
}

func TestOrderRepository_CreateOrderItems(t *testing.T) {
	repo := NewMockOrderRepository()
	ctx := context.Background()

	items := []models.OrderItem{
		{
			ProductID: "test-product-1",
			Quantity:  2,
			Price:     10.99,
		},
		{
			ProductID: "test-product-2",
			Quantity:  1,
			Price:     15.99,
		},
	}

	err := repo.CreateOrderItems(ctx, "test-order-1", items)
	assert.NoError(t, err)
}

func TestOrderRepository_GetOrderItems(t *testing.T) {
	repo := NewMockOrderRepository()
	ctx := context.Background()

	items, err := repo.GetOrderItems(ctx, "test-order-1")
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "test-product-1", items[0].ProductID)
	assert.Equal(t, 2, items[0].Quantity)
	assert.Equal(t, 10.99, items[0].Price)

	// Test non-existing order
	items, err = repo.GetOrderItems(ctx, "non-existing")
	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestOrderRepository_Delete(t *testing.T) {
	repo := NewMockOrderRepository()
	ctx := context.Background()

	err := repo.Delete(ctx, "test-order-1")
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}
