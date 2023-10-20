-- name: CreateSession :one

INSERT INTO session (
  id,
  username,
  refresh_token,
  useragent,
  client_ip,
  isblocked,
  expire_at
  
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetSession :one
SELECT * FROM session
WHERE id = $1 LIMIT 1;



-- name: DeleteSession :exec
DELETE FROM session
WHERE id = $1;