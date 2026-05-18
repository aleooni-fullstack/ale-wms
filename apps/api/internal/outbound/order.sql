-- name: CreateOrder :one
INSERT INTO orders (id, reference, status, note, created_at, updated_at)
VALUES ($1, $2, $3, $4, now(), now())
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1;

-- name: GetOrderByReference :one
SELECT * FROM orders
WHERE reference = $1;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (id, order_id, product_id, location_id, quantity, created_at)
VALUES ($1, $2, $3, $4, $5, now())
RETURNING *;

-- name: ListOrderItems :many
SELECT * FROM order_items
WHERE order_id = $1;