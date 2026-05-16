-- name: GetStockBalance :one
SELECT * FROM stock_balances
WHERE product_id = $1 AND location_id = $2;

-- name: ListStockBalancesByProduct :many
SELECT * FROM stock_balances
WHERE product_id = $1
ORDER BY location_id;

-- name: ListStockBalancesByLocation :many
SELECT * FROM stock_balances
WHERE location_id = $1
ORDER BY product_id;

-- name: UpsertStockBalance :one
INSERT INTO stock_balances (id, product_id, location_id, quantity, reserved_quantity, updated_at)
VALUES ($1, $2, $3, $4, $5, now())
ON CONFLICT (product_id, location_id)
DO UPDATE SET quantity = $4, reserved_quantity = $5, updated_at = now()
RETURNING *;

-- name: UpdateReservedQuantity :one
UPDATE stock_balances
SET reserved_quantity = $3, updated_at = now()
WHERE product_id = $1 AND location_id = $2
RETURNING *;