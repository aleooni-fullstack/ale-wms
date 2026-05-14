-- name: CreateStockTransfer :one
INSERT INTO stock_transfers (id, product_id, from_location_id, to_location_id, quantity, status, note, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
RETURNING *;

-- name: GetStockTransferByID :one
SELECT * FROM stock_transfers
WHERE id = $1;

-- name: ListStockTransfersByProduct :many
SELECT * FROM stock_transfers
WHERE product_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListStockTransfersByFromLocation :many
SELECT * FROM stock_transfers
WHERE from_location_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListStockTransfersByToLocation :many
SELECT * FROM stock_transfers
WHERE to_location_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateStockTransferStatus :one
UPDATE stock_transfers
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING *;