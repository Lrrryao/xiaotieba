// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: comment.sql

package db

import (
	"context"
)

const createComment = `-- name: CreateComment :one
INSERT INTO comment(
  content,
  user_id,
  post_id
) VALUES (
  $1, $2, $3
) RETURNING id, content, user_id, post_id, created_at
`

type CreateCommentParams struct {
	Content string `json:"content"`
	UserID  int64  `json:"user_id"`
	PostID  int64  `json:"post_id"`
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error) {
	row := q.db.QueryRowContext(ctx, createComment, arg.Content, arg.UserID, arg.PostID)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.UserID,
		&i.PostID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteComment = `-- name: DeleteComment :exec
DELETE FROM comment
WHERE id = $1
`

func (q *Queries) DeleteComment(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteComment, id)
	return err
}

const getComment = `-- name: GetComment :one
SELECT id, content, user_id, post_id, created_at FROM comment
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetComment(ctx context.Context, id int64) (Comment, error) {
	row := q.db.QueryRowContext(ctx, getComment, id)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.UserID,
		&i.PostID,
		&i.CreatedAt,
	)
	return i, err
}

const getCommentForUpdate = `-- name: GetCommentForUpdate :one
SELECT id, content, user_id, post_id, created_at FROM comment
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE
`

func (q *Queries) GetCommentForUpdate(ctx context.Context, id int64) (Comment, error) {
	row := q.db.QueryRowContext(ctx, getCommentForUpdate, id)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.UserID,
		&i.PostID,
		&i.CreatedAt,
	)
	return i, err
}
