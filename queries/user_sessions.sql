-- name: CreateSession :one
INSERT INTO user_sessions (user_id, session_token, jwt_token_id, device_info, ip_address, expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM user_sessions
WHERE session_token = $1;

-- name: ListSessionsByUser :many
SELECT * FROM user_sessions
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateSessionLastUsed :one
UPDATE user_sessions
SET last_used_at = NOW()
WHERE session_id = $1
    RETURNING *;

-- name: DeleteSession :exec
DELETE FROM user_sessions
WHERE session_id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM user_sessions
WHERE expires_at < NOW();
