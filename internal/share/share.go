package share

import (
	"os"

	"github.com/golang-jwt/jwt/v4"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

type MyError struct {
	Message string
	Error   string
}

type Claims struct {
	UserId               int `json:"uid"`
	jwt.RegisteredClaims `json:"claims"`
}
