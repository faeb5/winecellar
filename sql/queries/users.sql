-- name: CreateUser :one
INSERT INTO users (
    id, 
    email,
    hashed_password
) VALUES (?, ?, ?)
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?
LIMIT 1;

-- name: DeleteAllUsers :exec
DELETE FROM users;
