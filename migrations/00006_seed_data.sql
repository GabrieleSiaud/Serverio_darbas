-- +goose Up
INSERT INTO users (name, surname, username, email, password)
VALUES
    ('Admin', 'Admin', 'admin', 'admin@example.com', 'admin123'),
    ('John', 'Doe', 'johndoe', 'john@example.com', 'password123');

INSERT INTO games (title, description, release_date)
VALUES
    ('Diablo III', 'Action RPG', '2012-05-15'),
    ('Overwatch', 'Team Shooter', '2016-05-24');

INSERT INTO saved_games (user_id, game_id)
SELECT u.id, g.id FROM users u, games g WHERE u.username='johndoe' AND g.title='Diablo III';

INSERT INTO reviews (user_id, game_id, rating, comment)
SELECT u.id, g.id, 5, 'Awesome game!' FROM users u, games g WHERE u.username='johndoe' AND g.title='Diablo III';
-- +goose Down
DELETE FROM reviews;
DELETE FROM saved_games;
DELETE FROM games;
DELETE FROM users;
