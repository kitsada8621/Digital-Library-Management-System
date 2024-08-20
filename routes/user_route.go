package routes

import (
	"dlms/controllers"
	"dlms/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRoute(r *gin.RouterGroup) {
	userController := controllers.NewUserController()
	userRoute := r.Group("user")
	userRoute.Use(middleware.JwtAuth())
	userRoute.GET("/all", middleware.Role([]string{"admin"}), userController.GetUsers)
}
