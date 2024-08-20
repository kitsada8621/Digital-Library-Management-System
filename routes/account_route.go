package routes

import (
	"dlms/controllers"
	"dlms/middleware"

	"github.com/gin-gonic/gin"
)

func InitAccountRoute(r *gin.RouterGroup) {

	accountController := controllers.NewAccountController()

	r.POST("/login", accountController.Login)
	r.POST("/register", accountController.Register)
	r.GET("/profile", middleware.JwtAuth(), accountController.Profile)
	r.POST("/refresh/token", middleware.JwtRefreshAuth(), accountController.RefreshToken)
}
