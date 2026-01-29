package repository

import (
	"context"
	"database/sql"
	"fmt"

	"oolio/internal/app/models"
	"oolio/internal/database/sqlc"

	"github.com/google/uuid"
)

type orderRepository struct {
	db  *sql.DB
	qtx *sqlc.Queries
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{
		db:  db,
		qtx: sqlc.New(db),
	}
}

func (r *orderRepository) Find(ctx context.Context) ([]models.Order, error) {
	// This would need to be implemented with a new query in SQLC
	// For now, returning empty slice
	return []models.Order{}, nil
}

func (r *orderRepository) FindOne(ctx context.Context, id string) (*models.Order, error) {
	orderUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	dbOrder, err := r.qtx.GetOrderByID(ctx, orderUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Get order items
	orderItems, err := r.qtx.GetOrderItemsByOrderID(ctx, uuid.NullUUID{UUID: orderUUID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	order := r.mapSQLCToModel(dbOrder, orderItems)
	return &order, nil
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	params := sqlc.CreateOrderParams{
		Total:     fmt.Sprintf("%.2f", order.Total),
		Discounts: stringToNullString(fmt.Sprintf("%.2f", order.Discounts)),
		Status:    stringToNullString("pending"),
	}

	dbOrder, err := r.qtx.CreateOrder(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Update the order with the generated ID
	order.ID = dbOrder.ID.String()

	// Create order items
	if len(order.Items) > 0 {
		err = r.CreateOrderItems(ctx, order.ID, order.Items)
		if err != nil {
			return fmt.Errorf("failed to create order items: %w", err)
		}
	}

	return nil
}

func (r *orderRepository) Update(ctx context.Context, order *models.Order) error {
	orderUUID, err := uuid.Parse(order.ID)
	if err != nil {
		return fmt.Errorf("invalid order ID: %w", err)
	}

	_, err = r.qtx.UpdateOrderStatus(ctx, sqlc.UpdateOrderStatusParams{
		ID:     orderUUID,
		Status: stringToNullString("completed"),
	})
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

func (r *orderRepository) Delete(ctx context.Context, id string) error {
	// Order deletion would need to be implemented with proper cascade handling
	// For now, not implemented as per business requirements
	return fmt.Errorf("order deletion not implemented")
}

func (r *orderRepository) CreateOrderItems(ctx context.Context, orderID string, items []models.OrderItem) error {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return fmt.Errorf("invalid order ID: %w", err)
	}

	for _, item := range items {
		productUUID, err := uuid.Parse(item.ProductID)
		if err != nil {
			return fmt.Errorf("invalid product ID: %w", err)
		}

		// We need to get the current product price at time of order
		// For now, using a placeholder price - in real implementation this would come from product service
		params := sqlc.CreateOrderItemsParams{
			OrderID:     uuid.NullUUID{UUID: orderUUID, Valid: true},
			ProductID:   uuid.NullUUID{UUID: productUUID, Valid: true},
			Quantity:    int32(item.Quantity),
			PriceAtTime: "0.00", // This should be the actual product price at time of order
		}

		_, err = r.qtx.CreateOrderItems(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}
	}

	return nil
}

func (r *orderRepository) GetOrderItems(ctx context.Context, orderID string) ([]models.OrderItem, error) {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	dbOrderItems, err := r.qtx.GetOrderItemsByOrderID(ctx, uuid.NullUUID{UUID: orderUUID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	items := make([]models.OrderItem, len(dbOrderItems))
	for i, dbItem := range dbOrderItems {
		productID := ""
		if dbItem.ProductID.Valid {
			productID = dbItem.ProductID.UUID.String()
		}
		items[i] = models.OrderItem{
			ProductID: productID,
			Quantity:  int(dbItem.Quantity),
			Price:     parseFloat(dbItem.PriceAtTime),
		}
	}

	return items, nil
}

func (r *orderRepository) mapSQLCToModel(dbOrder sqlc.Order, dbOrderItems []sqlc.GetOrderItemsByOrderIDRow) models.Order {
	orderItems := make([]models.OrderItem, len(dbOrderItems))
	for i, dbItem := range dbOrderItems {
		productID := ""
		if dbItem.ProductID.Valid {
			productID = dbItem.ProductID.UUID.String()
		}
		orderItems[i] = models.OrderItem{
			ProductID: productID,
			Quantity:  int(dbItem.Quantity),
			Price:     parseFloat(dbItem.PriceAtTime),
		}
	}

	return models.Order{
		ID:        dbOrder.ID.String(),
		Total:     parseFloat(dbOrder.Total),
		Discounts: parseFloat(nullStringToString(dbOrder.Discounts)),
		Items:     orderItems,
	}
}
