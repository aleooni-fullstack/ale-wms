-- name: CreatePutAway :one
INSERT INTO put_aways (id, receipt_id, status, note, created_at, updated_at)
VALUES ($1, $2, $3, $4, now(), now())
RETURNING *;

-- name: GetPutAwayByID :one
SELECT * FROM put_aways
WHERE id = $1;

-- name: GetPutAwayByReceiptID :one
SELECT * FROM put_aways
WHERE receipt_id = $1;

-- name: UpdatePutAwayStatus :one
UPDATE put_aways
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: CreatePutAwayItem :one
INSERT INTO put_away_items (id, put_away_id, product_id, location_id, quantity, put_away, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, now(), now())
RETURNING *;

-- name: ListPutAwayItems :many
SELECT * FROM put_away_items
WHERE put_away_id = $1;

-- name: UpdatePutAwayItemPutAway :one
UPDATE put_away_items
SET put_away = $3, updated_at = now()
WHERE put_away_id = $1 AND product_id = $2
RETURNING *;