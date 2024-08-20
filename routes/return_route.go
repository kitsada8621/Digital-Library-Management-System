package routes

import (
	"dlms/controllers"
	"dlms/middleware"

	"github.com/gin-gonic/gin"
)

func InitReturnBookRoute(r *gin.RouterGroup) {
	returnController := controllers.NewReturnController()

	returnRoute := r.Group("return")
	returnRoute.Use(middleware.JwtAuth())
	returnRoute.PUT("/book/:id", middleware.Role([]string{"user"}), returnController.ReturnBook)
	returnRoute.PUT("/book/approve/:id", middleware.Role([]string{"admin"}), returnController.ApproveBookReturn)
}
