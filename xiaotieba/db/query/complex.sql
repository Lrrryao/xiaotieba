-- name: ListUser :many

SELECT users.id, users.name ,COUNT(*) AS UserCount
FROM users
INNER JOIN post ON users.id = post.user_id
GROUP BY users.id,users.name
HAVING COUNT(*) >= $1
ORDER BY UserCount DESC;
