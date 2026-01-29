package repository

import (
	"context"

	"oolio/internal/app/models"
)

type BaseRepository[T any] interface {
	Find(ctx context.Context) ([]T, error)
	FindOne(ctx context.Context, id string) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
}

type ProductRepository interface {
	BaseRepository[models.Product]
}

type OrderRepository interface {
	BaseRepository[models.Order]
	CreateOrderItems(ctx context.Context, orderID string, items []models.OrderItem) error
	GetOrderItems(ctx context.Context, orderID string) ([]models.OrderItem, error)
}
