package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"oolio/internal/app/models"
	"oolio/internal/database/sqlc"

	"github.com/google/uuid"
)

type productRepository struct {
	db  *sql.DB
	qtx *sqlc.Queries
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{
		db:  db,
		qtx: sqlc.New(db),
	}
}

func (r *productRepository) Find(ctx context.Context) ([]models.Product, error) {
	dbProducts, err := r.qtx.GetProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return r.mapSQLCToModels(dbProducts), nil
}

func (r *productRepository) FindOne(ctx context.Context, id string) (*models.Product, error) {
	productUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	dbProduct, err := r.qtx.GetProductByID(ctx, productUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	product := r.mapSQLCToModel(dbProduct)
	return &product, nil
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	params := sqlc.CreateProductParams{
		Name:         product.Name,
		Price:        fmt.Sprintf("%.2f", product.Price),
		Category:     product.Category,
		ThumbnailUrl: stringToNullString(product.Image.Thumbnail),
		MobileUrl:    stringToNullString(product.Image.Mobile),
		TabletUrl:    stringToNullString(product.Image.Tablet),
		DesktopUrl:   stringToNullString(product.Image.Desktop),
	}

	dbProduct, err := r.qtx.CreateProduct(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	// Update the product with the generated ID
	product.ID = dbProduct.ID.String()
	return nil
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	productUUID, err := uuid.Parse(product.ID)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	params := sqlc.UpdateProductParams{
		ID:           productUUID,
		Name:         product.Name,
		Price:        fmt.Sprintf("%.2f", product.Price),
		Category:     product.Category,
		ThumbnailUrl: stringToNullString(product.Image.Thumbnail),
		MobileUrl:    stringToNullString(product.Image.Mobile),
		TabletUrl:    stringToNullString(product.Image.Tablet),
		DesktopUrl:   stringToNullString(product.Image.Desktop),
	}

	_, err = r.qtx.UpdateProduct(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("product not found")
		}
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	productUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	err = r.qtx.DeleteProduct(ctx, productUUID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (r *productRepository) mapSQLCToModels(dbProducts []sqlc.Product) []models.Product {
	products := make([]models.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		products[i] = r.mapSQLCToModel(dbProduct)
	}
	return products
}

func (r *productRepository) mapSQLCToModel(dbProduct sqlc.Product) models.Product {
	return models.Product{
		ID:       dbProduct.ID.String(),
		Name:     dbProduct.Name,
		Price:    parseFloat(dbProduct.Price),
		Category: dbProduct.Category,
		Image: models.Image{
			Thumbnail: nullStringToString(dbProduct.ThumbnailUrl),
			Mobile:    nullStringToString(dbProduct.MobileUrl),
			Tablet:    nullStringToString(dbProduct.TabletUrl),
			Desktop:   nullStringToString(dbProduct.DesktopUrl),
		},
	}
}

func parseFloat(s string) float64 {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return 0.0
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func stringToNullString(s string) sql.NullString {
	if strings.TrimSpace(s) == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
