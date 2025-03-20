-- name: CreateWine :one
INSERT INTO wines (
    id,
    name,
    color,
    producer,
    country,
    vintage
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteAllWines :exec
DELETE FROM wines;
