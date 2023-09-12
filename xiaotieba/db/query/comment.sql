-- name: CreateComment :one
INSERT INTO comment(
  content,
  user_id,
  post_id
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetComment :one
SELECT * FROM comment
WHERE id = $1 LIMIT 1;

-- name: GetCommentForUpdate :one
SELECT * FROM comment
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;


-- name: DeleteComment :exec
DELETE FROM comment
WHERE id = $1;