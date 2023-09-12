-- name: CreateReply :one
INSERT INTO reply(
   content,
   user_id ,
   post_id,
   comment_id
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetReply :one
SELECT * FROM reply
WHERE id = $1 LIMIT 1;

-- name: GetReplyForUpdate :one
SELECT * FROM reply
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;


-- name: DeleteReply :exec
DELETE FROM reply
WHERE id = $1;