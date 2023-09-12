// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: reply.sql

package db

import (
	"context"
)

const createReply = `-- name: CreateReply :one
INSERT INTO reply(
   content,
   user_id ,
   post_id,
   comment_id
) VALUES (
  $1, $2, $3, $4
) RETURNING id, content, user_id, post_id, comment_id, created_at
`

type CreateReplyParams struct {
	Content   string `json:"content"`
	UserID    int64  `json:"user_id"`
	PostID    int64  `json:"post_id"`
	CommentID int64  `json:"comment_id"`
}

func (q *Queries) CreateReply(ctx context.Context, arg CreateReplyParams) (Reply, error) {
	row := q.db.QueryRowContext(ctx, createReply,
		arg.Content,
		arg.UserID,
		arg.PostID,
		arg.CommentID,
	)
	var i Reply
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.UserID,
		&i.PostID,
		&i.CommentID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteReply = `-- name: DeleteReply :exec
DELETE FROM reply
WHERE id = $1
`

func (q *Queries) DeleteReply(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteReply, id)
	return err
}

const getReply = `-- name: GetReply :one
SELECT id, content, user_id, post_id, comment_id, created_at FROM reply
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetReply(ctx context.Context, id int64) (Reply, error) {
	row := q.db.QueryRowContext(ctx, getReply, id)
	var i Reply
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.UserID,
		&i.PostID,
		&i.CommentID,
		&i.CreatedAt,
	)
	return i, err
}

const getReplyForUpdate = `-- name: GetReplyForUpdate :one
SELECT id, content, user_id, post_id, comment_id, created_at FROM reply
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE
`

func (q *Queries) GetReplyForUpdate(ctx context.Context, id int64) (Reply, error) {
	row := q.db.QueryRowContext(ctx, getReplyForUpdate, id)
	var i Reply
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.UserID,
		&i.PostID,
		&i.CommentID,
		&i.CreatedAt,
	)
	return i, err
}
