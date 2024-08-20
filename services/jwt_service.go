package services

import (
	"dlms/models"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IJwtService interface {
	GenerateToken(user models.User) (*string, error)
	VerifyToken(tokenString string) (*models.JwtClaims, error)
	RefreshToken(user models.User) (*string, error)
	VerifyRefreshToken(tokenString string) (*models.JwtClaims, error)
}
type JwtServiceImpl struct {
	roleService IRoleService
}

func JwtService() IJwtService {
	return &JwtServiceImpl{
		roleService: RoleService(),
	}
}

func (s *JwtServiceImpl) GenerateToken(user models.User) (*string, error) {

	expireTime, err := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRES"))
	if err != nil {
		return nil, err
	}

	roles := []string{}
	for _, roleId := range user.Roles {
		if role, err := s.roleService.FindById(roleId.Hex()); err == nil {
			roles = append(roles, role.RoleName)
		}
	}

	claims := models.JwtClaims{
		roles,
		jwt.RegisteredClaims{
			Subject:   user.ID.Hex(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expireTime))),
			Issuer:    os.Getenv("APP_NAME"),
		},
	}

	jwtKey := os.Getenv("JWT_ACCESS_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (s *JwtServiceImpl) VerifyToken(tokenString string) (*models.JwtClaims, error) {

	jwtKey := os.Getenv("JWT_ACCESS_SECRET")
	token, err := jwt.ParseWithClaims(tokenString, &models.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtKey), nil
	})
	if err != nil {
		jwtError := handleJWTError(err)
		if jwtError != "" {
			return nil, fmt.Errorf(jwtError)
		}
	} else if claims, ok := token.Claims.(*models.JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
func (s *JwtServiceImpl) RefreshToken(user models.User) (*string, error) {

	expireTime, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRES"))
	if err != nil {
		return nil, err
	}

	roles := []string{}
	for _, roleId := range user.Roles {
		if role, err := s.roleService.FindById(roleId.Hex()); err == nil {
			roles = append(roles, role.RoleName)
		}
	}

	claims := models.JwtClaims{
		roles,
		jwt.RegisteredClaims{
			Subject:   user.ID.Hex(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expireTime))),
			Issuer:    os.Getenv("APP_NAME"),
		},
	}

	jwtKey := os.Getenv("JWT_REFRESH_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (s *JwtServiceImpl) VerifyRefreshToken(tokenString string) (*models.JwtClaims, error) {

	jwtKey := os.Getenv("JWT_REFRESH_SECRET")
	token, err := jwt.ParseWithClaims(tokenString, &models.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtKey), nil
	})

	if err != nil {
		jwtError := handleJWTError(err)
		if jwtError != "" {
			return nil, fmt.Errorf(jwtError)
		}
	} else if claims, ok := token.Claims.(*models.JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func handleJWTError(err error) string {
	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		log.Println("The token is expired")
		return "Token has expired, please login again"
	case errors.Is(err, jwt.ErrTokenMalformed):
		log.Println("Malformed token")
		return "Malformed token"
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		log.Println("The token is not valid yet")
		return "Token not valid yet"
	default:
		log.Println("Couldn't handle this token:", err)
		return "Invalid token"
	}
}
