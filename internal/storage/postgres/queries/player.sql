-- name: GetPlayer :one
SELECT * FROM players WHERE id = $1;

-- name: CreatePlayer :one
INSERT INTO players (
    user_id
)
VALUES (
    $1
)
RETURNING *;

-- name: SavePlayer :exec
UPDATE players
SET
    x = $2,
    y = $3,
    hp = $4,
    mana = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetPlayerByUserId :one
SELECT * FROM players WHERE user_id = $1 LIMIT 1;

-- name: GetPlayerByID :one
SELECT * FROM players WHERE id = $1;

-- name: GetPlayersByUserId :many
SELECT * FROM players WHERE user_id = $1;

-- name: GetAllPlayers :many
SELECT * FROM players;