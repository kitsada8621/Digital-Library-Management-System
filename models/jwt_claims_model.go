package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}
