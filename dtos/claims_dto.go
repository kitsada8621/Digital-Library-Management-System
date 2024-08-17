package dtos

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}
