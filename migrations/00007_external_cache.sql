-- +goose Up
-- +goose StatementBegin

CREATE TABLE external_api_cache (
                                    cache_key TEXT PRIMARY KEY,
                                    response_json JSONB NOT NULL,
                                    expires_at TIMESTAMPTZ NOT NULL,
                                    created_at TIMESTAMPTZ DEFAULT NOW(),
                                    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_external_api_cache_expires ON external_api_cache(expires_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS external_api_cache;
-- +goose StatementEnd
