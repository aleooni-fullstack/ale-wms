-- name: CreatePurchaseOrder :one
INSERT INTO purchase_orders (id, reference, supplier, status, note, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, now(), now())
RETURNING *;

-- name: GetPurchaseOrderByID :one
SELECT * FROM purchase_orders
WHERE id = $1;

-- name: GetPurchaseOrderByReference :one
SELECT * FROM purchase_orders
WHERE reference = $1;

-- name: ListPurchaseOrders :many
SELECT * FROM purchase_orders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdatePurchaseOrderStatus :one
UPDATE purchase_orders
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: CreatePurchaseOrderItem :one
INSERT INTO purchase_order_items (id, purchase_order_id, product_id, quantity, created_at)
VALUES ($1, $2, $3, $4, now())
RETURNING *;

-- name: ListPurchaseOrderItems :many
SELECT * FROM purchase_order_items
WHERE purchase_order_id = $1;