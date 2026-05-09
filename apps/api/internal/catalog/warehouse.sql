-- name: CreateWarehouse :one
INSERT INTO warehouses (id, code, name, address, active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, now(), now())
RETURNING *;

-- name: GetWarehouseByID :one
SELECT * FROM warehouses
WHERE id = $1;

-- name: GetWarehouseByCode :one
SELECT * FROM warehouses
WHERE code = $1;

-- name: ListWarehouses :many
SELECT * FROM warehouses
WHERE active = true
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: UpdateWarehouse :one
UPDATE warehouses
SET code = $2, name = $3, address = $4, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeactivateWarehouse :exec
UPDATE warehouses
SET active = false, updated_at = now()
WHERE id = $1;