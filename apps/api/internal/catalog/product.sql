-- name: CreateProduct :one
INSERT INTO products (id, sku, name, description, unit, active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, now(), now())
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1;

-- name: GetProductBySKU :one
SELECT * FROM products
WHERE sku = $1;

-- name: ListProducts :many
SELECT * FROM products
WHERE active = true
ORDER BY name;

-- name: UpdateProduct :one
UPDATE products
SET sku = $2, name = $3, description = $4, unit = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeactivateProduct :exec
UPDATE products
SET active = false, updated_at = now()
WHERE id = $1;