-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users ( 
    username, 
    password_hash,
    created_at
) 
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET
    username = $2,
    password_hash = $3
WHERE id = $1;