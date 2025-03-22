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
