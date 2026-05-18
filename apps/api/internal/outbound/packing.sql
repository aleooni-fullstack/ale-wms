-- name: CreatePacking :one
INSERT INTO packings (id, order_id, picking_id, status, note, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, now(), now())
RETURNING *;

-- name: GetPackingByID :one
SELECT * FROM packings
WHERE id = $1;

-- name: GetPackingByOrderID :one
SELECT * FROM packings
WHERE order_id = $1;

-- name: UpdatePackingStatus :one
UPDATE packings
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING *;