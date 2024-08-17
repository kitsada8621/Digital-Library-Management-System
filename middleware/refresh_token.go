package middleware

import (
	"dlms/services"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthorizeRefreshToken() gin.HandlerFunc {
	jwtService := services.JwtService()
	return func(ctx *gin.Context) {

		header := ctx.GetHeader("Authorization")
		fmt.Printf("header: %s\n", header)
		if header == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				`success`: false,
				`message`: "authorization is required",
			})
			ctx.Abort()
			return
		}

		parts := strings.Split(header, " ")
		fmt.Println(parts)
		if len(parts) != 2 || strings.ToLower(parts[0]) == "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				`success`: false,
				`message`: "invalid token",
			})
			ctx.Abort()
			return
		}

		claims, err := jwtService.VerifyRefreshToken(parts[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				`success`: false,
				`message`: err,
			})
			ctx.Abort()
			return
		}

		ctx.Set("userId", claims.Subject)
		ctx.Next()
	}
}
