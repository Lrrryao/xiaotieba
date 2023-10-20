package db

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var ErrRecordNotFound = sql.ErrNoRows

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

// 将error值转换为postgres的错误码（字符串）
// 如果该错误不是PQ数据库的错误的话，就返回一个空字符串
func ErrorCode(err error) pq.ErrorCode {
	var pqErr *pq.Error

	if errors.As(err, &pqErr) {
		return pqErr.Code
	}
	return ""
}
