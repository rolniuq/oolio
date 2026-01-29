package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"oolio/internal/app/handler"
	"oolio/internal/app/models"
)

// Mock service for testing
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductService) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) CreateProduct(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductService) UpdateProduct(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductService) DeleteProduct(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestProductHandler_ListProducts(t *testing.T) {
	mockService := &MockProductService{}
	handler := handler.NewProductHandler(mockService)
	ctx := context.Background()

	expectedProducts := []models.Product{
		{
			ID:       "test-1",
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
	}

	mockService.On("GetAllProducts", ctx).Return(expectedProducts, nil)

	// Setup Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/product", nil)

	handler.ListProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Product
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, expectedProducts[0].ID, response[0].ID)
	assert.Equal(t, expectedProducts[0].Name, response[0].Name)
	assert.Equal(t, expectedProducts[0].Price, response[0].Price)

	mockService.AssertExpectations(t)
}

func TestProductHandler_ListProducts_Error(t *testing.T) {
	mockService := &MockProductService{}
	handler := handler.NewProductHandler(mockService)
	ctx := context.Background()

	mockService.On("GetAllProducts", ctx).Return([]models.Product{}, assert.AnError)

	// Setup Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/product", nil)

	handler.ListProducts(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "error", response.Type)

	mockService.AssertExpectations(t)
}

func TestProductHandler_GetProduct(t *testing.T) {
	mockService := &MockProductService{}
	handler := handler.NewProductHandler(mockService)
	ctx := context.Background()

	expectedProduct := &models.Product{
		ID:       "test-1",
		Name:     "Test Product 1",
		Price:    10.99,
		Category: "Waffle",
		Image: models.Image{
			Thumbnail: "http://example.com/thumb.jpg",
			Mobile:    "http://example.com/mobile.jpg",
			Tablet:    "http://example.com/tablet.jpg",
			Desktop:   "http://example.com/desktop.jpg",
		},
	}

	mockService.On("GetProductByID", ctx, "test-1").Return(expectedProduct, nil)

	// Setup Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/product/test-1", nil)
	c.Params = gin.Params{{Key: "productId", Value: "test-1"}}

	handler.GetProduct(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Product
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct.ID, response.ID)
	assert.Equal(t, expectedProduct.Name, response.Name)
	assert.Equal(t, expectedProduct.Price, response.Price)

	mockService.AssertExpectations(t)
}

func TestProductHandler_GetProduct_NotFound(t *testing.T) {
	mockService := &MockProductService{}
	handler := handler.NewProductHandler(mockService)
	ctx := context.Background()

	mockService.On("GetProductByID", ctx, "not-found").Return(nil, assert.AnError)

	// Setup Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/product/not-found", nil)
	c.Params = gin.Params{{Key: "productId", Value: "not-found"}}

	handler.GetProduct(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "error", response.Type)

	mockService.AssertExpectations(t)
}

func TestProductHandler_GetProduct_EmptyID(t *testing.T) {
	handler := handler.NewProductHandler(&MockProductService{})

	// Setup Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/product/", nil)
	c.Params = gin.Params{{Key: "productId", Value: ""}}

	handler.GetProduct(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "error", response.Type)
}
