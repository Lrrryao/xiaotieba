// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: user_role.sql

package db

import (
	"context"
)

const getUserRoles = `-- name: GetUserRoles :one
SELECT ur_id, user_id, user_name, role_id, role_name FROM users_roles
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetUserRoles(ctx context.Context, userID int32) (UsersRole, error) {
	row := q.db.QueryRowContext(ctx, getUserRoles, userID)
	var i UsersRole
	err := row.Scan(
		&i.UrID,
		&i.UserID,
		&i.UserName,
		&i.RoleID,
		&i.RoleName,
	)
	return i, err
}
