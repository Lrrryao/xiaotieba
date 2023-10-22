-- name: GetUserRoles :one
SELECT * FROM users_roles
WHERE user_id = $1 LIMIT 1;