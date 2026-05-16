-- name: CreateStockReservation :one
INSERT INTO stock_reservations (id, product_id, location_id, quantity, status, reference, note, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
RETURNING *;

-- name: GetStockReservationByID :one
SELECT * FROM stock_reservations
WHERE id = $1;

-- name: ListStockReservationsByProduct :many
SELECT * FROM stock_reservations
WHERE product_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListStockReservationsByLocation :many
SELECT * FROM stock_reservations
WHERE location_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListStockReservationsByReference :many
SELECT * FROM stock_reservations
WHERE reference = $1
ORDER BY created_at DESC;

-- name: UpdateStockReservationStatus :one
UPDATE stock_reservations
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING *;