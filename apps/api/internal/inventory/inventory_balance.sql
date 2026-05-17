-- name: CreateInventoryBalance :one
INSERT INTO inventory_balances (id, location_id, status, note, created_at, updated_at)
VALUES ($1, $2, $3, $4, now(), now())
RETURNING *;

-- name: GetInventoryBalanceByID :one
SELECT * FROM inventory_balances
WHERE id = $1;

-- name: ListInventoryBalancesByLocation :many
SELECT * FROM inventory_balances
WHERE location_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateInventoryBalanceStatus :one
UPDATE inventory_balances
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: CreateInventoryBalanceItem :one
INSERT INTO inventory_balance_items (id, inventory_balance_id, product_id, system_quantity, counted_quantity, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, now(), now())
RETURNING *;

-- name: GetInventoryBalanceItem :one
SELECT * FROM inventory_balance_items
WHERE inventory_balance_id = $1 AND product_id = $2;

-- name: ListInventoryBalanceItems :many
SELECT * FROM inventory_balance_items
WHERE inventory_balance_id = $1
ORDER BY product_id;

-- name: UpdateInventoryBalanceItemCountedQuantity :one
UPDATE inventory_balance_items
SET counted_quantity = $3, updated_at = now()
WHERE inventory_balance_id = $1 AND product_id = $2
RETURNING *;