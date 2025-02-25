-- name: CreateWine :one
INSERT INTO wines (
    id,
    name,
    color,
    wine_maker,
    country,
    vintage
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteAllWines :exec
DELETE FROM wines;
