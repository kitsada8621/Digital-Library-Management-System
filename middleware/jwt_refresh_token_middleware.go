package middleware

import (
	"dlms/models"
	"dlms/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtRefreshAuth() gin.HandlerFunc {
	jwtService := services.JwtService()
	return func(ctx *gin.Context) {

		header := ctx.GetHeader("Authorization")
		if header == "" {
			ctx.JSON(http.StatusUnauthorized, models.ResponseJson{
				Success: false,
				Message: "authorization is required",
			})
			ctx.Abort()
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) == "Bearer" {
			ctx.JSON(http.StatusUnauthorized, models.ResponseJson{
				Success: false,
				Message: "invalid token",
			})
			ctx.Abort()
			return
		}

		claims, err := jwtService.VerifyRefreshToken(parts[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, models.ResponseJson{
				Success: false,
				Message: err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.Set("userId", claims.Subject)
		ctx.Next()
	}
}
