-- name: GetProducts :many
SELECT id, name, price, category, thumbnail_url, mobile_url, tablet_url, desktop_url, created_at, updated_at
FROM products
ORDER BY name;

-- name: GetProductByID :one
SELECT id, name, price, category, thumbnail_url, mobile_url, tablet_url, desktop_url, created_at, updated_at
FROM products
WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO products (name, price, category, thumbnail_url, mobile_url, tablet_url, desktop_url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, price, category, thumbnail_url, mobile_url, tablet_url, desktop_url, created_at, updated_at;

-- name: UpdateProduct :one
UPDATE products 
SET name = $2, price = $3, category = $4, thumbnail_url = $5, mobile_url = $6, tablet_url = $7, desktop_url = $8, updated_at = NOW()
WHERE id = $1
RETURNING id, name, price, category, thumbnail_url, mobile_url, tablet_url, desktop_url, created_at, updated_at;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;