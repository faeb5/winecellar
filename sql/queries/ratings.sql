-- name: CreateRating :one
INSERT INTO ratings (
    id,
    wine_id,
    user_id,
    rating
) VALUES (?, ?, ?, ?)
RETURNING *;

-- name: DeleteAllRatings :exec
DELETE FROM ratings;
