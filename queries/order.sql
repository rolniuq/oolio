-- name: CreateOrder :one
INSERT INTO orders (total, discounts, status)
VALUES ($1, $2, $3)
RETURNING id, total, discounts, status, created_at, updated_at;

-- name: GetOrderByID :one
SELECT id, total, discounts, status, created_at, updated_at
FROM orders
WHERE id = $1;

-- name: CreateOrderItems :many
INSERT INTO order_items (order_id, product_id, quantity, price_at_time)
VALUES ($1, $2, $3, $4)
RETURNING id, order_id, product_id, quantity, price_at_time, created_at;

-- name: GetOrderItemsByOrderID :many
SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price_at_time, oi.created_at,
       p.name, p.category, p.thumbnail_url, p.mobile_url, p.tablet_url, p.desktop_url
FROM order_items oi
JOIN products p ON oi.product_id = p.id
WHERE oi.order_id = $1;

-- name: UpdateOrderStatus :one
UPDATE orders 
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, total, discounts, status, created_at, updated_at;