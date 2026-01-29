package services

import (
	"context"
	"fmt"

	"oolio/internal/app/models"
	"oolio/internal/app/repository"
)

type ProductService interface {
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) error
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id string) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	products, err := s.repo.Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}

	return products, nil
}

func (s *productService) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	if id == "" {
		return nil, fmt.Errorf("product ID cannot be empty")
	}

	product, err := s.repo.FindOne(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by ID %s: %w", id, err)
	}

	return product, nil
}

func (s *productService) CreateProduct(ctx context.Context, product *models.Product) error {
	if err := s.validateProduct(product); err != nil {
		return fmt.Errorf("product validation failed: %w", err)
	}

	err := s.repo.Create(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (s *productService) UpdateProduct(ctx context.Context, product *models.Product) error {
	if err := s.validateProduct(product); err != nil {
		return fmt.Errorf("product validation failed: %w", err)
	}

	if product.ID == "" {
		return fmt.Errorf("product ID is required for update")
	}

	err := s.repo.Update(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("product ID cannot be empty")
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (s *productService) validateProduct(product *models.Product) error {
	if product == nil {
		return fmt.Errorf("product cannot be nil")
	}

	if product.Name == "" {
		return fmt.Errorf("product name is required")
	}

	if product.Price <= 0 {
		return fmt.Errorf("product price must be greater than 0")
	}

	if product.Category == "" {
		return fmt.Errorf("product category is required")
	}

	return nil
}
