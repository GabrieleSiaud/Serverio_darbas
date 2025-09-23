-- name: SaveGame :one
INSERT INTO saved_games (user_id, game_id)
VALUES ($1, $2)
    ON CONFLICT (user_id, game_id) DO NOTHING
RETURNING *;

-- name: GetSavedGame :one
SELECT * FROM saved_games
WHERE id = $1;

-- name: ListSavedGamesByUser :many
SELECT sg.*, g.title, g.description, g.release_date
FROM saved_games sg
         JOIN games g ON sg.game_id = g.id
WHERE sg.user_id = $1
ORDER BY sg.created_at DESC;

-- name: DeleteSavedGame :exec
DELETE FROM saved_games
WHERE user_id = $1 AND game_id = $2;
