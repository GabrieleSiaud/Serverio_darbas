-- +goose Up
-- +goose StatementBegin

CREATE TABLE saved_games (
                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                             user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                             game_id UUID NOT NULL REFERENCES games(game_id) ON DELETE CASCADE,
                             created_at TIMESTAMPTZ DEFAULT NOW(),
                             updated_at TIMESTAMPTZ DEFAULT NOW(),
                             UNIQUE(user_id, game_id)
);

-- Indexai patogiam query
CREATE INDEX idx_saved_games_user_id ON saved_games(user_id);
CREATE INDEX idx_saved_games_game_id ON saved_games(game_id);

-- Trigger, kad updated_at atsinaujintų
CREATE TRIGGER update_saved_games_updated_at
    BEFORE UPDATE ON saved_games
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_saved_games_updated_at ON saved_games;
DROP INDEX IF EXISTS idx_saved_games_user_id;
DROP INDEX IF EXISTS idx_saved_games_game_id;
DROP TABLE IF EXISTS saved_games CASCADE;

-- +goose StatementEnd
