package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token is invalid")
	ErrInvalidToken = errors.New("token is expired")
)

type payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	CreateAt time.Time `json:"create_at"`
	ExpireAt time.Time `json:"expire_at"`
}

func NewPayload(username string, duration time.Duration) (*payload, error) {
	TokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &payload{
		ID:       TokenID, //Random出来的tokenID即为session的ID(前提是refreshToken,因为session里只存储了refreshToken)
		Username: username,
		CreateAt: time.Now(),
		ExpireAt: time.Now().Add(duration),
	}
	return payload, nil

}

//这里还有第二种NewPayload函数，而且其中使用的paseto.JSONToken是由paseto库设计的标准json token
//我们不在这个项目中是由这个NewPayload函数是因为，一旦按照该标准创建payload
//就无法使用其他token_maker类型，如jwt
// paseto.JSONToken defines standard token payload claims and allows for additional claims to be added.
// All of the standard claims are optional.
// func NewPayload2(username string, duration time.Duration) (*payload, error) {
// 	jsonToken := paseto.JSONToken{
//         Audience:   "test",
//         Issuer:     "test_service",
//         Jti:        "123",
//         Subject:    "test_subject",
//         IssuedAt:   now,
//         Expiration: exp,
//         NotBefore:  nbt,
//         }
// // Add custom claim    to the token
// jsonToken.Set("data", "this is a signed message")
// footer := "some footer"
// }

// 完成paseto库中标准jsonPayload的接口
func Valid(payload *payload) error {
	if time.Now().After(payload.ExpireAt) {
		return ErrExpiredToken
	}
	return nil
}
