-- +goose Up
-- +goose StatementBegin

CREATE TABLE roles (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE user_roles (
                            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                            role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
                            PRIMARY KEY (user_id, role_id)
);

-- Seed roles
INSERT INTO roles (name) VALUES ('admin'), ('moderator'), ('user');

-- Assign admin role to admin@example.com
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.email = 'admin@example.com' AND r.name = 'admin';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS roles;

-- +goose StatementEnd
