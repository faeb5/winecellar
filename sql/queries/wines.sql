-- name: CreateWine :one
INSERT INTO wines (
    id,
    name,
    color,
    producer,
    country,
    vintage,
    created_by
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteAllWines :exec
DELETE FROM wines;

-- name: GetWineByProducerAndNameAndVintage :one
SELECT *
FROM wines
WHERE producer = ?
    AND name = ?
    AND vintage = ?;

-- name: DeleteWine :exec
DELETE FROM wines WHERE id = ?;

-- name: GetWineByID :one
SELECT *
FROM wines
WHERE id = ?;

-- name: UpdateWineByID :one
UPDATE wines SET
    color = ?,
    name = ?,
    producer = ?,
    country = ?,
    vintage = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: GetAllWines :many
SELECT *
FROM wines;
