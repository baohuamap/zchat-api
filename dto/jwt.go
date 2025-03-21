package dto

import "github.com/golang-jwt/jwt"

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.MapClaims
}
