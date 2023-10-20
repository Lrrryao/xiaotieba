package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	//chacha20poly是一个算法
)

// PasetoMaker is a PASETO token maker
// paseto之所以被封装进PasetoMaker结构体里， 是因为它封装好了加密、解码等功能
type paseto_maker struct {
	symmetricKey []byte
	paseto       *paseto.V2
}

func NewPaseto_Maker(key string) (Maker, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key: length of the key is not %d", chacha20poly1305.KeySize)
	}
	maker := &paseto_maker{
		symmetricKey: []byte(key),
		paseto:       paseto.NewV2(),
	}
	return maker, nil
}

func (maker *paseto_maker) CreateToken(username string, duration time.Duration) (string, *payload, error) {
	Payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}

	//paseto_maker自带的函数，用于生成token
	token, err := maker.paseto.Encrypt(maker.symmetricKey, Payload, nil)

	return token, Payload, err
}

func (maker *paseto_maker) VerifyToken(token string) (*payload, error) {
	payload := &payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if err := Valid(payload); err != nil {
		return nil, ErrExpiredToken
	}

	return payload, nil
}

// footer 可以包含任何你认为有用的元数据，例如令牌的版本、颁发者、颁发时间、过期时间或其他自定义信息。这些元数据可以帮助你更好地理解和管理令牌，
// 但它们不会被解密，只有令牌的有效载荷（payload）部分才会被解密。
// footer := map[string]interface{}{
//     "version":  "1.0",          // 令牌版本
//     "issuedAt": time.Now().Unix(), // 颁发时间
// }
// 然后，你可以使用这个 footer 信息在验证令牌时进行参考或记录。

// 总之，footer 是用于存储不需要被解密但需要与令牌一起传输的元数据的字段。根据你的需求，你可以选择使用它来附加有关令牌的额外信息。
