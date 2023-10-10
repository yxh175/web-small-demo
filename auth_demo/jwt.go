package auth_demo

import "github.com/dgrijalva/jwt-go"

type JWTClaims struct {
	ID       uint   `json:"id"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}
