package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Hashpassword(password string) (string, error) {
	hashpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("fail to hash word: %v", err)
	}
	return string(hashpassword), nil
}

func CheckPassWord(password string, hashedpassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
}
