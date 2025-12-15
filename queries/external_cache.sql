-- name: GetExternalCache :one
SELECT cache_key, response_json, expires_at
FROM external_api_cache
WHERE cache_key = $1 AND expires_at > NOW();

-- name: UpsertExternalCache :exec
INSERT INTO external_api_cache (cache_key, response_json, expires_at, updated_at)
VALUES ($1, $2, $3, NOW())
    ON CONFLICT (cache_key)
DO UPDATE SET response_json = EXCLUDED.response_json,
           expires_at = EXCLUDED.expires_at,
           updated_at = NOW();

-- name: DeleteExpiredExternalCache :exec
DELETE FROM external_api_cache
WHERE expires_at <= NOW();
