// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	PostID    int64     `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Titles    string    `json:"titles"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Power struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Reply struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	PostID    int64     `json:"post_id"`
	CommentID int64     `json:"comment_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Role struct {
	ID       int64  `json:"id"`
	RoleName string `json:"role_name"`
}

type RolesPower struct {
	RpID      int64  `json:"rp_id"`
	RoleID    int32  `json:"role_id"`
	RoleName  string `json:"role_name"`
	PowerID   int32  `json:"power_id"`
	PowerName string `json:"power_name"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	ExpireAt     time.Time `json:"expire_at"`
	CreatedAt    time.Time `json:"created_at"`
	ClientIp     string    `json:"client_ip"`
	Useragent    string    `json:"useragent"`
	Isblocked    bool      `json:"isblocked"`
	RefreshToken string    `json:"refresh_token"`
}

type User struct {
	ID              int64     `json:"id"`
	HashPassword    string    `json:"hash_password"`
	Name            string    `json:"name"`
	Power           string    `json:"power"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
	Phone           string    `json:"phone"`
	IsEmailVerified bool      `json:"is_email_verified"`
}

type UsersRole struct {
	UrID     int32  `json:"ur_id"`
	UserID   int32  `json:"user_id"`
	UserName string `json:"user_name"`
	RoleID   int32  `json:"role_id"`
	RoleName string `json:"role_name"`
}

type VerifyEmail struct {
	ID              int64          `json:"id"`
	Username        sql.NullString `json:"username"`
	Email           sql.NullString `json:"email"`
	SecretCode      string         `json:"secret_code"`
	IsSecretUsed    bool           `json:"is_secret_used"`
	SecretCreatedAt time.Time      `json:"secret_created_at"`
	SecretExpiredAt time.Time      `json:"secret_expired_at"`
}
