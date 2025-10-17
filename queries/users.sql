-- name: CreateUser :one
INSERT INTO users (name, surname, username, email, password)
VALUES ($1, $2, $3, $4, $5)
    RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET name = $2,
    surname = $3,
    username = $4,
    email = $5,
    password = $6,
    updated_at = NOW()
WHERE id = $1
    RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login = NOW()
WHERE id = $1;

-- name: DeleteSessionByToken :exec
DELETE FROM user_sessions
WHERE session_token = $1;