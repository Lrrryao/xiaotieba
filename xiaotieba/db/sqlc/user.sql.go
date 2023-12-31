// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  hash_password,
  name,
  power,
  email,
  phone
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING id, hash_password, name, power, email, created_at, phone, is_email_verified
`

type CreateUserParams struct {
	HashPassword string `json:"hash_password"`
	Name         string `json:"name"`
	Power        string `json:"power"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.HashPassword,
		arg.Name,
		arg.Power,
		arg.Email,
		arg.Phone,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.HashPassword,
		&i.Name,
		&i.Power,
		&i.Email,
		&i.CreatedAt,
		&i.Phone,
		&i.IsEmailVerified,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, hash_password, name, power, email, created_at, phone, is_email_verified FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.HashPassword,
		&i.Name,
		&i.Power,
		&i.Email,
		&i.CreatedAt,
		&i.Phone,
		&i.IsEmailVerified,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, hash_password, name, power, email, created_at, phone, is_email_verified FROM users
WHERE name = $1 LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, name string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, name)
	var i User
	err := row.Scan(
		&i.ID,
		&i.HashPassword,
		&i.Name,
		&i.Power,
		&i.Email,
		&i.CreatedAt,
		&i.Phone,
		&i.IsEmailVerified,
	)
	return i, err
}

const getUserForUpdate = `-- name: GetUserForUpdate :one
SELECT id, hash_password, name, power, email, created_at, phone, is_email_verified FROM users
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE
`

func (q *Queries) GetUserForUpdate(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserForUpdate, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.HashPassword,
		&i.Name,
		&i.Power,
		&i.Email,
		&i.CreatedAt,
		&i.Phone,
		&i.IsEmailVerified,
	)
	return i, err
}
