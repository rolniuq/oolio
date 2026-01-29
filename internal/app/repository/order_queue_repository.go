package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"oolio/internal/app/models"
)

type OrderQueueRepository interface {
	AddToQueue(ctx context.Context, item *models.OrderQueueItem) error
	GetPendingItems(ctx context.Context, batchSize int) ([]*models.OrderQueueItem, error)
	UpdateItem(ctx context.Context, item *models.OrderQueueItem) error
	MarkAsProcessing(ctx context.Context, itemID string) error
	MarkAsCompleted(ctx context.Context, itemID string, order *models.Order) error
	MarkAsFailed(ctx context.Context, itemID string, errorMsg string) error
	GetQueueStats(ctx context.Context) (map[string]int, error)
	GetOrderFromQueue(ctx context.Context, itemID string) (*models.OrderQueueItem, error)
	GetAllOrders(ctx context.Context) ([]*models.OrderQueueItem, error)
}

type orderQueueRepository struct {
	db *sql.DB
}

func NewOrderQueueRepository(db *sql.DB) OrderQueueRepository {
	return &orderQueueRepository{db: db}
}

func (r *orderQueueRepository) AddToQueue(ctx context.Context, item *models.OrderQueueItem) error {
	orderReqJSON, err := json.Marshal(item.OrderReq)
	if err != nil {
		return fmt.Errorf("failed to marshal order request: %w", err)
	}

	query := `
		INSERT INTO order_queue (id, order_req, status, created_at, updated_at, retry_count)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = r.db.ExecContext(ctx, query, item.ID, orderReqJSON, item.Status, item.CreatedAt, item.UpdatedAt, item.RetryCount)
	if err != nil {
		return fmt.Errorf("failed to insert into order queue: %w", err)
	}

	return nil
}

func (r *orderQueueRepository) GetPendingItems(ctx context.Context, batchSize int) ([]*models.OrderQueueItem, error) {
	query := `
		SELECT id, order_req, status, created_at, updated_at, error, order_data, retry_count
		FROM order_queue
		WHERE status = 'pending' 
		OR (status = 'failed' AND retry_count < 3)
		ORDER BY created_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`

	rows, err := r.db.QueryContext(ctx, query, batchSize)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending items: %w", err)
	}
	defer rows.Close()

	var items []*models.OrderQueueItem
	for rows.Next() {
		item := &models.OrderQueueItem{}
		var orderReqJSON []byte
		var orderData []byte
		var error sql.NullString

		err := rows.Scan(
			&item.ID,
			&orderReqJSON,
			&item.Status,
			&item.CreatedAt,
			&item.UpdatedAt,
			&error,
			&orderData,
			&item.RetryCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan queue item: %w", err)
		}

		if err := json.Unmarshal(orderReqJSON, &item.OrderReq); err != nil {
			return nil, fmt.Errorf("failed to unmarshal order request: %w", err)
		}

		if error.Valid {
			item.Error = error.String
		}

		if len(orderData) > 0 {
			var order models.Order
			if err := json.Unmarshal(orderData, &order); err == nil {
				item.Order = &order
			}
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating queue items: %w", err)
	}

	return items, nil
}

func (r *orderQueueRepository) UpdateItem(ctx context.Context, item *models.OrderQueueItem) error {
	orderDataJSON := []byte("{}")
	if item.Order != nil {
		var err error
		orderDataJSON, err = json.Marshal(item.Order)
		if err != nil {
			return fmt.Errorf("failed to marshal order data: %w", err)
		}
	}

	query := `
		UPDATE order_queue 
		SET status = $2, updated_at = $3, error = $4, order_data = $5, retry_count = $6
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, item.ID, item.Status, item.UpdatedAt, item.Error, orderDataJSON, item.RetryCount)
	if err != nil {
		return fmt.Errorf("failed to update queue item: %w", err)
	}

	return nil
}

func (r *orderQueueRepository) MarkAsProcessing(ctx context.Context, itemID string) error {
	query := `UPDATE order_queue SET status = 'processing', updated_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), itemID)
	return err
}

func (r *orderQueueRepository) MarkAsCompleted(ctx context.Context, itemID string, order *models.Order) error {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	query := `
		UPDATE order_queue 
		SET status = 'completed', updated_at = $1, order_data = $2, error = NULL
		WHERE id = $3
	`
	_, err = r.db.ExecContext(ctx, query, time.Now(), orderJSON, itemID)
	return err
}

func (r *orderQueueRepository) MarkAsFailed(ctx context.Context, itemID string, errorMsg string) error {
	query := `
		UPDATE order_queue 
		SET status = 'failed', updated_at = $1, error = $2, retry_count = retry_count + 1
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, time.Now(), errorMsg, itemID)
	return err
}

func (r *orderQueueRepository) GetQueueStats(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT status, COUNT(*) 
		FROM order_queue 
		GROUP BY status
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan queue stats: %w", err)
		}
		stats[status] = count
	}

	return stats, nil
}

func (r *orderQueueRepository) GetOrderFromQueue(ctx context.Context, itemID string) (*models.OrderQueueItem, error) {
	query := `
		SELECT id, order_req, status, created_at, updated_at, error, order_data, retry_count
		FROM order_queue
		WHERE id = $1
	`

	var item models.OrderQueueItem
	var orderReqJSON []byte
	var orderData []byte
	var error sql.NullString

	err := r.db.QueryRowContext(ctx, query, itemID).Scan(
		&item.ID,
		&orderReqJSON,
		&item.Status,
		&item.CreatedAt,
		&item.UpdatedAt,
		&error,
		&orderData,
		&item.RetryCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order from queue: %w", err)
	}

	if err := json.Unmarshal(orderReqJSON, &item.OrderReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order request: %w", err)
	}

	if error.Valid {
		item.Error = error.String
	}

	if len(orderData) > 0 {
		var order models.Order
		if err := json.Unmarshal(orderData, &order); err == nil {
			item.Order = &order
		}
	}

	return &item, nil
}

func (r *orderQueueRepository) GetAllOrders(ctx context.Context) ([]*models.OrderQueueItem, error) {
	query := `
		SELECT id, order_req, status, created_at, updated_at, error, order_data, retry_count
		FROM order_queue
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders: %w", err)
	}
	defer rows.Close()

	var orders []*models.OrderQueueItem
	for rows.Next() {
		var item models.OrderQueueItem
		var orderReqJSON []byte
		var orderData []byte
		var error sql.NullString

		err := rows.Scan(
			&item.ID,
			&orderReqJSON,
			&item.Status,
			&item.CreatedAt,
			&item.UpdatedAt,
			&error,
			&orderData,
			&item.RetryCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		if err := json.Unmarshal(orderReqJSON, &item.OrderReq); err != nil {
			return nil, fmt.Errorf("failed to unmarshal order request: %w", err)
		}

		if error.Valid {
			item.Error = error.String
		}

		if len(orderData) > 0 {
			var order models.Order
			if err := json.Unmarshal(orderData, &order); err == nil {
				item.Order = &order
			}
		}

		orders = append(orders, &item)
	}

	return orders, nil
}
