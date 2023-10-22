// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error)
	CreatePost(ctx context.Context, arg CreatePostParams) (Post, error)
	CreateReply(ctx context.Context, arg CreateReplyParams) (Reply, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (VerifyEmail, error)
	DeleteComment(ctx context.Context, id int64) error
	DeletePost(ctx context.Context, id int64) error
	DeleteReply(ctx context.Context, id int64) error
	DeleteSession(ctx context.Context, id uuid.UUID) error
	DeleteUser(ctx context.Context, id int64) error
	GetComment(ctx context.Context, id int64) (Comment, error)
	GetCommentForUpdate(ctx context.Context, id int64) (Comment, error)
	GetPost(ctx context.Context, id int64) (Post, error)
	GetPostForUpdate(ctx context.Context, id int64) (Post, error)
	GetReply(ctx context.Context, id int64) (Reply, error)
	GetReplyForUpdate(ctx context.Context, id int64) (Reply, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	GetUserByUsername(ctx context.Context, name string) (User, error)
	GetUserForUpdate(ctx context.Context, id int64) (User, error)
	GetUserRoles(ctx context.Context, userID int32) (UsersRole, error)
	ListUser(ctx context.Context, dollar_1 interface{}) ([]ListUserRow, error)
}

var _ Querier = (*Queries)(nil)
