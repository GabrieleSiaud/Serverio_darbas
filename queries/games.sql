-- name: CreateGame :one
INSERT INTO games (title, description, release_date)
VALUES ($1, $2, $3)
    RETURNING *;

-- name: GetGameByID :one
SELECT * FROM games
WHERE id = $1;

-- name: ListGames :many
SELECT * FROM games
ORDER BY created_at DESC;

-- name: UpdateGame :one
UPDATE games
SET title = $2,
    description = $3,
    release_date = $4,
    updated_at = NOW()
WHERE id = $1
    RETURNING *;

-- name: DeleteGame :exec
DELETE FROM games
WHERE id = $1;
