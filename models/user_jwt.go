package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserJwt struct {
	SID  int64  `json:"sid"`
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func (u *UserJwt) GenerateJwt(secret []byte) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": u,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	})
	s, err := t.SignedString(secret)
	if err != nil {
		return "", err
	}
	return s, nil
}
