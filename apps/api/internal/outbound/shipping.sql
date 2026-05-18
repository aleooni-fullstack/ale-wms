-- name: CreateShipping :one
INSERT INTO shippings (id, order_id, packing_id, status, tracking_code, note, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, now(), now())
RETURNING *;

-- name: GetShippingByID :one
SELECT * FROM shippings
WHERE id = $1;

-- name: GetShippingByOrderID :one
SELECT * FROM shippings
WHERE order_id = $1;

-- name: UpdateShippingStatus :one
UPDATE shippings
SET status = $2, tracking_code = $3, updated_at = now()
WHERE id = $1
RETURNING *;