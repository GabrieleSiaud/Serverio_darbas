-- +goose Up
-- +goose StatementBegin

CREATE TABLE reviews (
                         review_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         game_id   UUID NOT NULL REFERENCES games(game_id) ON DELETE CASCADE,
                         user_id   UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
                         rating    SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 5),
                         comment   TEXT,
                         created_at TIMESTAMPTZ DEFAULT NOW(),
                         updated_at TIMESTAMPTZ DEFAULT NOW(),
                         UNIQUE(user_id, game_id) -- vienas vartotojas vienam žaidimui
);

CREATE INDEX idx_reviews_user_id ON reviews(user_id);
CREATE INDEX idx_reviews_game_id ON reviews(game_id);

CREATE TRIGGER update_reviews_updated_at
    BEFORE UPDATE ON reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_reviews_updated_at ON reviews;
DROP INDEX IF EXISTS idx_reviews_user_id;
DROP INDEX IF EXISTS idx_reviews_game_id;
DROP TABLE IF EXISTS reviews CASCADE;

-- +goose StatementEnd
