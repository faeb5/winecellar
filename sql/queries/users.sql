-- name: CreateUser :one
INSERT INTO users (
    id, 
    username,
    email,
    hashed_password
) VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ?;

-- name: DeleteAllUsers :exec
DELETE FROM users;
