-- name: CreateZone :one
INSERT INTO zones (id, warehouse_id, code, name, active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, now(), now())
RETURNING *;

-- name: GetZoneByID :one
SELECT * FROM zones
WHERE id = $1;

-- name: GetZoneByCode :one
SELECT * FROM zones
WHERE warehouse_id = $1 AND code = $2;

-- name: ListZonesByWarehouse :many
SELECT * FROM zones
WHERE warehouse_id = $1 AND active = true
ORDER BY code
LIMIT $2 OFFSET $3;

-- name: UpdateZone :one
UPDATE zones
SET code = $2, name = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeactivateZone :exec
UPDATE zones
SET active = false, updated_at = now()
WHERE id = $1;