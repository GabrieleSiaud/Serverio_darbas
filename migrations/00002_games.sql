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
INSERT INTO games (id, title, description, release_date, created_at, updated_at)
VALUES (
           uuid_generate_v4(),
           'The Witcher 3: Wild Hunt',
           'An open-world RPG where players follow Geralt of Rivia in a story-rich adventure full of choices and consequences.',
           '2015-05-19',
           NOW(),
           NOW()
       );

INSERT INTO games (id, title, description, release_date, created_at, updated_at)
VALUES (
           uuid_generate_v4(),
           'Half-Life 2',
           'A groundbreaking first-person shooter with physics-based gameplay and an immersive sci-fi story.',
           '2004-11-16',
           NOW(),
           NOW()
       );

INSERT INTO games (id, title, description, release_date, created_at, updated_at)
VALUES (
           uuid_generate_v4(),
           'Minecraft',
           'A sandbox game where players build, explore, and survive in blocky procedurally generated worlds.',
           '2011-11-18',
           NOW(),
           NOW()
       );

INSERT INTO games (id, title, description, release_date, created_at, updated_at)
VALUES (
           uuid_generate_v4(),
           'Cyberpunk 2077',
           'A futuristic open-world RPG set in Night City, offering deep storytelling and character customization.',
           '2020-12-10',
           NOW(),
           NOW()
       );

INSERT INTO games (id, title, description, release_date, created_at, updated_at)
VALUES (
           uuid_generate_v4(),
           'League of Legends',
           'A popular MOBA where players compete in strategic 5v5 battles with diverse champions.',
           '2009-10-27',
           NOW(),
           NOW()
       );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_games_updated_at ON games;
DROP INDEX IF EXISTS idx_games_title;
DROP TABLE IF EXISTS games;
-- +goose StatementEnd




