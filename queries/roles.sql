-- name: ListUserRoles :many
SELECT r.name
FROM user_roles ur
         JOIN roles r ON r.id = ur.role_id
WHERE ur.user_id = $1;

-- name: AddUserRole :exec
INSERT INTO user_roles (user_id, role_id)
SELECT $1, r.id
FROM roles r
WHERE r.name = $2
    ON CONFLICT DO NOTHING;
