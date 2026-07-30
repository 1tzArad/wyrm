-- name: GetPlayer :one
SELECT * FROM players WHERE id = $1;

-- name: CreatePlayer :one
INSERT INTO players (
    id,
    user_id
)
VALUES (
    $1,
    $2
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