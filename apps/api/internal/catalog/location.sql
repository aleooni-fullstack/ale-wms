-- name: CreateLocation :one
INSERT INTO locations (id, zone_id, code, name, active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, now(), now())
RETURNING *;

-- name: GetLocationByID :one
SELECT * FROM locations
WHERE id = $1;

-- name: GetLocationByCode :one
SELECT * FROM locations
WHERE zone_id = $1 AND code = $2;

-- name: ListLocationsByZone :many
SELECT * FROM locations
WHERE zone_id = $1 AND active = true
ORDER BY code
LIMIT $2 OFFSET $3;

-- name: UpdateLocation :one
UPDATE locations
SET code = $2, name = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeactivateLocation :exec
UPDATE locations
SET active = false, updated_at = now()
WHERE id = $1;