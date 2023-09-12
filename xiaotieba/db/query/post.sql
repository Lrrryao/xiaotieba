-- name: CreatePost :one
INSERT INTO post (
  user_id,
  titles,
  content
 
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetPost :one
SELECT * FROM post
WHERE id = $1 LIMIT 1;


-- name: GetPostForUpdate :one
SELECT * FROM post
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;


-- name: DeletePost :exec
DELETE FROM post
WHERE id = $1;