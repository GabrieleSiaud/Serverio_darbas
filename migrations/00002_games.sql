-- +goose Up
-- +goose StatementBegin

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Games table
CREATE TABLE games (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       title VARCHAR(255) NOT NULL,
                       description TEXT,
                       release_date DATE,
                       created_at TIMESTAMPTZ DEFAULT NOW(),
                       updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create index on title for faster searching
CREATE INDEX idx_games_title ON games(title);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

-- Create trigger for updated_at
CREATE TRIGGER update_games_updated_at
    BEFORE UPDATE ON games
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO games (id, title, description, release_date, created_at, updated_at)
VALUES (
           uuid_generate_v4(),
           'Mortal Kombat',
           'A classic fighting game franchise featuring brutal combat and iconic characters.',
           '1992-10-08',
           NOW(),
           NOW()
       );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop trigger
DROP TRIGGER IF EXISTS update_games_updated_at ON games;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop index
DROP INDEX IF EXISTS idx_games_title;

-- Drop games table
DROP TABLE IF EXISTS games CASCADE;

-- +goose StatementEnd
