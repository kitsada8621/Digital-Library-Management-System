package routes

import (
	"dlms/controllers"
	"dlms/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoleRoute(r *gin.RouterGroup) {
	roleController := controllers.NewRoleController()
	roleRoute := r.Group("role")
	roleRoute.Use(middleware.JwtAuth())
	roleRoute.GET("/all", middleware.Role([]string{"admin"}), roleController.GetRoles)
	roleRoute.POST("/assign", middleware.Role([]string{"admin"}), roleController.AssignRole)
}
