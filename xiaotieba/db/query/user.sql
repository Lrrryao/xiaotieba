-- name: CreateUser :one
INSERT INTO users (
  hash_password,
  name,
  power,
  email,
  phone
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE name = $1 LIMIT 1;

-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;



-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;