package token

import (
	"time"
)

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, *payload, error)
	VerifyToken(token string) (*payload, error)
}
