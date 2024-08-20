package middleware

import (
	"dlms/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Role(roles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		roleJson := ctx.GetString("roles")

		var userRoles []string
		if roleJson != "" {
			if err := json.Unmarshal([]byte(roleJson), &userRoles); err != nil {
				fmt.Println("\n\nErr encode role: ", err.Error())
			}
		}

		fmt.Println("Claims Role: ", userRoles)
		fmt.Println("Roles: ", roles)

		permission := false
		for _, role := range roles {
			for _, userRole := range userRoles {
				if role == userRole {
					permission = true
					break
				}
			}
		}

		if !permission {
			ctx.JSON(http.StatusForbidden, models.ResponseJson{
				Success: false,
				Message: "Access is denied",
			})
			ctx.Abort()
		}

		ctx.Next()
	}
}
