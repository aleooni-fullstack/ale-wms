-- name: CreatePicking :one
INSERT INTO pickings (id, order_id, status, note, created_at, updated_at)
VALUES ($1, $2, $3, $4, now(), now())
RETURNING *;

-- name: GetPickingByID :one
SELECT * FROM pickings
WHERE id = $1;

-- name: GetPickingByOrderID :one
SELECT * FROM pickings
WHERE order_id = $1;

-- name: UpdatePickingStatus :one
UPDATE pickings
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: CreatePickingItem :one
INSERT INTO picking_items (id, picking_id, product_id, location_id, quantity, picked, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, now(), now())
RETURNING *;

-- name: ListPickingItems :many
SELECT * FROM picking_items
WHERE picking_id = $1;

-- name: UpdatePickingItemPicked :one
UPDATE picking_items
SET picked = $3, updated_at = now()
WHERE picking_id = $1 AND product_id = $2
RETURNING *;