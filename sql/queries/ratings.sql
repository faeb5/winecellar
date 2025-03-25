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

-- name: GetRatingByID :one
SELECT *
FROM ratings
WHERE ID = ?;

-- name: UpdateRatingByID :one
UPDATE ratings
SET
    rating = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE ID = ?
RETURNING *;
