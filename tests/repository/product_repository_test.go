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

// Mock repository for testing
type mockProductRepository struct {
	products []models.Product
}

// Ensure mock implements the interface
var _ repository.ProductRepository = (*mockProductRepository)(nil)

func NewMockProductRepository() repository.ProductRepository {
	return &mockProductRepository{
		products: []models.Product{
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
			{
				ID:       "test-product-2",
				Name:     "Test Product 2",
				Price:    15.99,
				Category: "Waffle",
				Image: models.Image{
					Thumbnail: "http://example.com/thumb2.jpg",
					Mobile:    "http://example.com/mobile2.jpg",
					Tablet:    "http://example.com/tablet2.jpg",
					Desktop:   "http://example.com/desktop2.jpg",
				},
			},
		},
	}
}

func (r *mockProductRepository) Find(ctx context.Context) ([]models.Product, error) {
	return r.products, nil
}

func (r *mockProductRepository) FindOne(ctx context.Context, id string) (*models.Product, error) {
	for _, product := range r.products {
		if product.ID == id {
			return &product, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *mockProductRepository) Create(ctx context.Context, product *models.Product) error {
	product.ID = uuid.New().String()
	r.products = append(r.products, *product)
	return nil
}

func (r *mockProductRepository) Update(ctx context.Context, product *models.Product) error {
	for i, p := range r.products {
		if p.ID == product.ID {
			r.products[i] = *product
			return nil
		}
	}
	return sql.ErrNoRows
}

func (r *mockProductRepository) Delete(ctx context.Context, id string) error {
	for i, product := range r.products {
		if product.ID == id {
			r.products = append(r.products[:i], r.products[i+1:]...)
			return nil
		}
	}
	return sql.ErrNoRows
}

func TestProductRepository_Find(t *testing.T) {
	repo := NewMockProductRepository()
	ctx := context.Background()

	products, err := repo.Find(ctx)
	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Test Product 1", products[0].Name)
}

func TestProductRepository_FindOne(t *testing.T) {
	repo := NewMockProductRepository()
	ctx := context.Background()

	// Test existing product
	product, err := repo.FindOne(ctx, "test-product-1")
	assert.NoError(t, err)
	require.NotNil(t, product)
	assert.Equal(t, "Test Product 1", product.Name)
	assert.Equal(t, 10.99, product.Price)

	// Test non-existing product
	product, err = repo.FindOne(ctx, "non-existing")
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestProductRepository_Create(t *testing.T) {
	repo := NewMockProductRepository()
	ctx := context.Background()

	newProduct := &models.Product{
		Name:     "New Product",
		Price:    25.99,
		Category: "Waffle",
		Image: models.Image{
			Thumbnail: "http://example.com/new-thumb.jpg",
			Mobile:    "http://example.com/new-mobile.jpg",
			Tablet:    "http://example.com/new-tablet.jpg",
			Desktop:   "http://example.com/new-desktop.jpg",
		},
	}

	err := repo.Create(ctx, newProduct)
	assert.NoError(t, err)
	assert.NotEmpty(t, newProduct.ID)

	// Verify product was added
	products, err := repo.Find(ctx)
	assert.NoError(t, err)
	assert.Len(t, products, 3)
}

func TestProductRepository_Update(t *testing.T) {
	repo := NewMockProductRepository()
	ctx := context.Background()

	// Get existing product
	product, err := repo.FindOne(ctx, "test-product-1")
	require.NoError(t, err)
	require.NotNil(t, product)

	// Update product
	product.Name = "Updated Product"
	product.Price = 99.99

	err = repo.Update(ctx, product)
	assert.NoError(t, err)

	// Verify update
	updatedProduct, err := repo.FindOne(ctx, "test-product-1")
	assert.NoError(t, err)
	assert.Equal(t, "Updated Product", updatedProduct.Name)
	assert.Equal(t, 99.99, updatedProduct.Price)
}

func TestProductRepository_Update_NotFound(t *testing.T) {
	repo := NewMockProductRepository()
	ctx := context.Background()

	product := &models.Product{
		ID:       "non-existing",
		Name:     "Updated Product",
		Price:    99.99,
		Category: "Waffle",
	}

	err := repo.Update(ctx, product)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestProductRepository_Delete(t *testing.T) {
	repo := NewMockProductRepository()
	ctx := context.Background()

	err := repo.Delete(ctx, "test-product-1")
	assert.NoError(t, err)

	// Verify deletion
	products, err := repo.Find(ctx)
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "test-product-2", products[0].ID)
}

func TestProductRepository_Delete_NotFound(t *testing.T) {
	repo := NewMockProductRepository()
	ctx := context.Background()

	err := repo.Delete(ctx, "non-existing")
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}
