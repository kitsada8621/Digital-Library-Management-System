package middleware

import (
	"dlms/models"
	"dlms/services"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAuth() gin.HandlerFunc {
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

		claims, err := jwtService.VerifyToken(parts[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, models.ResponseJson{
				Success: false,
				Message: err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.Set("userId", claims.Subject)
		roleJson, err := json.Marshal(claims.Roles)
		ctx.Set("roles", string(roleJson))
		ctx.Next()
	}
}
