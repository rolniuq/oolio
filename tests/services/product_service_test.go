package services

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"oolio/internal/app/models"
	"oolio/internal/app/services"
)

// Mock repository for testing
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Find(ctx context.Context) ([]models.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) FindOne(ctx context.Context, id string) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) Create(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Update(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestProductService_GetAllProducts(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := services.NewProductService(mockRepo)
	ctx := context.Background()

	expectedProducts := []models.Product{
		{
			ID:       "test-1",
			Name:     "Test Product 1",
			Price:    10.99,
			Category: "Waffle",
		},
	}

	mockRepo.On("Find", ctx).Return(expectedProducts, nil)

	products, err := service.GetAllProducts(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProductByID(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := services.NewProductService(mockRepo)
	ctx := context.Background()

	expectedProduct := &models.Product{
		ID:       "test-1",
		Name:     "Test Product 1",
		Price:    10.99,
		Category: "Waffle",
	}

	mockRepo.On("FindOne", ctx, "test-1").Return(expectedProduct, nil)

	product, err := service.GetProductByID(ctx, "test-1")

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProductByID_EmptyID(t *testing.T) {
	service := services.NewProductService(&MockProductRepository{})
	ctx := context.Background()

	product, err := service.GetProductByID(ctx, "")

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "product ID cannot be empty")
}

func TestProductService_GetProductByID_NotFound(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := services.NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindOne", ctx, "not-found").Return(nil, sql.ErrNoRows)

	product, err := service.GetProductByID(ctx, "not-found")

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "failed to get product")
	mockRepo.AssertExpectations(t)
}

func TestProductService_CreateProduct(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := services.NewProductService(mockRepo)
	ctx := context.Background()

	product := &models.Product{
		Name:     "New Product",
		Price:    25.99,
		Category: "Waffle",
		Image: models.Image{
			Thumbnail: "http://example.com/thumb.jpg",
			Mobile:    "http://example.com/mobile.jpg",
			Tablet:    "http://example.com/tablet.jpg",
			Desktop:   "http://example.com/desktop.jpg",
		},
	}

	mockRepo.On("Create", ctx, product).Return(nil)

	err := service.CreateProduct(ctx, product)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_CreateProduct_ValidationError(t *testing.T) {
	service := services.NewProductService(&MockProductRepository{})
	ctx := context.Background()

	// Test with nil product
	err := service.CreateProduct(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product validation failed")

	// Test with empty name
	product := &models.Product{
		Name:     "",
		Price:    25.99,
		Category: "Waffle",
	}
	err = service.CreateProduct(ctx, product)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product name is required")

	// Test with invalid price
	product.Name = "Valid Name"
	product.Price = -1
	err = service.CreateProduct(ctx, product)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product price must be greater than 0")

	// Test with empty category
	product.Price = 25.99
	product.Category = ""
	err = service.CreateProduct(ctx, product)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product category is required")
}

func TestProductService_UpdateProduct(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := services.NewProductService(mockRepo)
	ctx := context.Background()

	product := &models.Product{
		ID:       "test-1",
		Name:     "Updated Product",
		Price:    99.99,
		Category: "Waffle",
		Image: models.Image{
			Thumbnail: "http://example.com/thumb.jpg",
			Mobile:    "http://example.com/mobile.jpg",
			Tablet:    "http://example.com/tablet.jpg",
			Desktop:   "http://example.com/desktop.jpg",
		},
	}

	mockRepo.On("Update", ctx, product).Return(nil)

	err := service.UpdateProduct(ctx, product)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateProduct_EmptyID(t *testing.T) {
	service := services.NewProductService(&MockProductRepository{})
	ctx := context.Background()

	product := &models.Product{
		Name:     "Updated Product",
		Price:    99.99,
		Category: "Waffle",
	}

	err := service.UpdateProduct(ctx, product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product ID is required")
}

func TestProductService_DeleteProduct(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := services.NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("Delete", ctx, "test-1").Return(nil)

	err := service.DeleteProduct(ctx, "test-1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_DeleteProduct_EmptyID(t *testing.T) {
	service := services.NewProductService(&MockProductRepository{})
	ctx := context.Background()

	err := service.DeleteProduct(ctx, "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product ID cannot be empty")
}
