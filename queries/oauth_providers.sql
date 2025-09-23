-- name: LinkOAuthProvider :one
INSERT INTO oauth_providers (user_id, provider, provider_user_id, provider_username, provider_email, access_token, refresh_token, token_expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    ON CONFLICT (provider, provider_user_id) DO UPDATE
                                                    SET access_token = EXCLUDED.access_token,
                                                    refresh_token = EXCLUDED.refresh_token,
                                                    token_expires_at = EXCLUDED.token_expires_at,
                                                    updated_at = NOW()
                                                    RETURNING *;

-- name: GetOAuthProviderByUser :one
SELECT * FROM oauth_providers
WHERE user_id = $1 AND provider = $2;

-- name: GetOAuthProviderByExternalID :one
SELECT * FROM oauth_providers
WHERE provider = $1 AND provider_user_id = $2;

-- name: DeleteOAuthProvider :exec
DELETE FROM oauth_providers
WHERE user_id = $1 AND provider = $2;
