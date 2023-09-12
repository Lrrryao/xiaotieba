// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: post.sql

package db

import (
	"context"
)

const createPost = `-- name: CreatePost :one
INSERT INTO post (
  user_id,
  titles,
  content
 
) VALUES (
  $1, $2, $3
) RETURNING id, user_id, titles, content, created_at
`

type CreatePostParams struct {
	UserID  int64  `json:"user_id"`
	Titles  string `json:"titles"`
	Content string `json:"content"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost, arg.UserID, arg.Titles, arg.Content)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Titles,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}

const deletePost = `-- name: DeletePost :exec
DELETE FROM post
WHERE id = $1
`

func (q *Queries) DeletePost(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletePost, id)
	return err
}

const getPost = `-- name: GetPost :one
SELECT id, user_id, titles, content, created_at FROM post
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPost(ctx context.Context, id int64) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPost, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Titles,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}

const getPostForUpdate = `-- name: GetPostForUpdate :one
SELECT id, user_id, titles, content, created_at FROM post
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE
`

func (q *Queries) GetPostForUpdate(ctx context.Context, id int64) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostForUpdate, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Titles,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}
