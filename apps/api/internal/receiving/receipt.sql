-- name: CreateReceipt :one
INSERT INTO receipts (id, purchase_order_id, status, note, created_at, updated_at)
VALUES ($1, $2, $3, $4, now(), now())
RETURNING *;

-- name: GetReceiptByID :one
SELECT * FROM receipts
WHERE id = $1;

-- name: GetReceiptByPurchaseOrderID :one
SELECT * FROM receipts
WHERE purchase_order_id = $1;

-- name: UpdateReceiptStatus :one
UPDATE receipts
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: CreateReceiptItem :one
INSERT INTO receipt_items (id, receipt_id, product_id, expected_quantity, received_quantity, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, now(), now())
RETURNING *;

-- name: ListReceiptItems :many
SELECT * FROM receipt_items
WHERE receipt_id = $1;

-- name: UpdateReceiptItemReceivedQuantity :one
UPDATE receipt_items
SET received_quantity = $3, updated_at = now()
WHERE receipt_id = $1 AND product_id = $2
RETURNING *;