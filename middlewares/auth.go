package middlewares

import (
	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}
