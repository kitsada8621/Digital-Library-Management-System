package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IJwtService interface {
	GenerateToken(sub string) (*string, error)
	VerifyToken(tokenString string) (*jwt.RegisteredClaims, error)
	RefreshToken(sub string) (*string, error)
	VerifyRefreshToken(tokenString string) (*jwt.RegisteredClaims, error)
}
type JwtServiceImpl struct{}

func JwtService() IJwtService {
	return &JwtServiceImpl{}
}

func (s *JwtServiceImpl) GenerateToken(sub string) (*string, error) {

	expireTime, err := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRES"))
	if err != nil {
		return nil, err
	}

	claims := &jwt.RegisteredClaims{
		Subject:   sub,
		Issuer:    os.Getenv("APP_NAME"),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expireTime))),
	}

	jwtKey := os.Getenv("JWT_ACCESS_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (s *JwtServiceImpl) VerifyToken(tokenString string) (*jwt.RegisteredClaims, error) {

	jwtKey := os.Getenv("JWT_ACCESS_SECRET")
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 	return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		// }

		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
func (s *JwtServiceImpl) RefreshToken(sub string) (*string, error) {

	expireTime, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRES"))
	if err != nil {
		return nil, err
	}

	claims := &jwt.RegisteredClaims{
		Subject:   sub,
		Issuer:    os.Getenv("APP_NAME"),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expireTime))),
	}

	jwtKey := os.Getenv("JWT_REFRESH_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (s *JwtServiceImpl) VerifyRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {

	jwtKey := os.Getenv("JWT_REFRESH_SECRET")
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
