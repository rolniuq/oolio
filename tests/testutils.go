package testutils

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"oolio/internal/app/models"
)

// CreateTestProduct creates a test product with default values
func CreateTestProduct(t *testing.T) *models.Product {
	return &models.Product{
		ID:       uuid.New().String(),
		Name:     "Test Product",
		Price:    10.99,
		Category: "Waffle",
		Image: models.Image{
			Thumbnail: "http://example.com/thumb.jpg",
			Mobile:    "http://example.com/mobile.jpg",
			Tablet:    "http://example.com/tablet.jpg",
			Desktop:   "http://example.com/desktop.jpg",
		},
	}
}

// CreateTestProductWithCustomData creates a test product with custom values
func CreateTestProductWithCustomData(t *testing.T, name string, price float64, category string) *models.Product {
	product := CreateTestProduct(t)
	product.Name = name
	product.Price = price
	product.Category = category
	return product
}

// CreateTestOrder creates a test order with default values
func CreateTestOrder(t *testing.T) *models.Order {
	return &models.Order{
		ID:        uuid.New().String(),
		Total:     25.99,
		Discounts: 0.0,
		Items: []models.OrderItem{
			{
				ProductID: uuid.New().String(),
				Quantity:  2,
				Price:     10.99,
			},
		},
		Products: []models.Product{
			*CreateTestProduct(t),
		},
	}
}

// CreateTestOrderReq creates a test order request with default values
func CreateTestOrderReq(t *testing.T) *models.OrderReq {
	return &models.OrderReq{
		Items: []models.OrderItem{
			{
				ProductID: uuid.New().String(),
				Quantity:  2,
			},
		},
		CouponCode: "",
	}
}

// CreateTestOrderReqWithItems creates a test order request with custom items
func CreateTestOrderReqWithItems(t *testing.T, items []models.OrderItem, couponCode string) *models.OrderReq {
	return &models.OrderReq{
		Items:      items,
		CouponCode: couponCode,
	}
}

// CreateTestOrderItem creates a test order item with default values
func CreateTestOrderItem(t *testing.T) *models.OrderItem {
	return &models.OrderItem{
		ProductID: uuid.New().String(),
		Quantity:  1,
		Price:     10.99,
	}
}

// CreateTestOrderItemWithCustomData creates a test order item with custom values
func CreateTestOrderItemWithCustomData(t *testing.T, productID string, quantity int, price float64) *models.OrderItem {
	return &models.OrderItem{
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
	}
}

// CreateTestProducts creates multiple test products
func CreateTestProducts(t *testing.T, count int) []models.Product {
	products := make([]models.Product, count)
	for i := 0; i < count; i++ {
		products[i] = *CreateTestProductWithCustomData(
			t,
			"Test Product "+string(rune('A'+i)),
			float64(10+i),
			"Category "+string(rune('A'+i)),
		)
	}
	return products
}

// AssertProductEqual asserts that two products are equal
func AssertProductEqual(t *testing.T, expected, actual *models.Product) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.Price, actual.Price)
	require.Equal(t, expected.Category, actual.Category)
	require.Equal(t, expected.Image.Thumbnail, actual.Image.Thumbnail)
	require.Equal(t, expected.Image.Mobile, actual.Image.Mobile)
	require.Equal(t, expected.Image.Tablet, actual.Image.Tablet)
	require.Equal(t, expected.Image.Desktop, actual.Image.Desktop)
}

// AssertOrderEqual asserts that two orders are equal
func AssertOrderEqual(t *testing.T, expected, actual *models.Order) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Total, actual.Total)
	require.Equal(t, expected.Discounts, actual.Discounts)
	require.Equal(t, len(expected.Items), len(actual.Items))
	require.Equal(t, len(expected.Products), len(actual.Products))
}

// AssertOrderItemsEqual asserts that two order items slices are equal
func AssertOrderItemsEqual(t *testing.T, expected, actual []models.OrderItem) {
	require.Equal(t, len(expected), len(actual))
	for i := range expected {
		require.Equal(t, expected[i].ProductID, actual[i].ProductID)
		require.Equal(t, expected[i].Quantity, actual[i].Quantity)
		require.Equal(t, expected[i].Price, actual[i].Price)
	}
}

// TestContext creates a test context
func TestContext() context.Context {
	return context.Background()
}

// TestContextWithTimeout creates a test context with timeout
func TestContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30)
}

// Cleanup provides a cleanup function for tests
func Cleanup(t *testing.T, cleanup func()) {
	t.Cleanup(cleanup)
}

// RequireNoError requires that error is nil
func RequireNoError(t *testing.T, err error) {
	require.NoError(t, err)
}

// RequireError requires that error is not nil
func RequireError(t *testing.T, err error) {
	require.Error(t, err)
}

// AssertEqual asserts that two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}) {
	assert.Equal(t, expected, actual)
}

// AssertNotEqual asserts that two values are not equal
func AssertNotEqual(t *testing.T, expected, actual interface{}) {
	assert.NotEqual(t, expected, actual)
}

// AssertTrue asserts that condition is true
func AssertTrue(t *testing.T, condition bool) {
	assert.True(t, condition)
}

// AssertFalse asserts that condition is false
func AssertFalse(t *testing.T, condition bool) {
	assert.False(t, condition)
}

// AssertNotEmpty asserts that value is not empty
func AssertNotEmpty(t *testing.T, value interface{}) {
	assert.NotEmpty(t, value)
}

// AssertNil asserts that value is nil
func AssertNil(t *testing.T, value interface{}) {
	assert.Nil(t, value)
}

// AssertNotNil asserts that value is not nil
func AssertNotNil(t *testing.T, value interface{}) {
	assert.NotNil(t, value)
}
