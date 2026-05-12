-- name: CreateStockMovement :one
INSERT INTO stock_movements (id, product_id, location_id, type, quantity, note, created_at)
VALUES ($1, $2, $3, $4, $5, $6, now())
RETURNING *;

-- name: GetStockMovementByID :one
SELECT * FROM stock_movements
WHERE id = $1;

-- name: ListStockMovementsByProduct :many
SELECT * FROM stock_movements
WHERE product_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListStockMovementsByLocation :many
SELECT * FROM stock_movements
WHERE location_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;