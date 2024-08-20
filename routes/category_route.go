package routes

import (
	"dlms/controllers"
	"dlms/middleware"

	"github.com/gin-gonic/gin"
)

func InitCategoryRoute(r *gin.RouterGroup) {
	categoryController := controllers.NewCategoryController()
	cateRoute := r.Group("category")
	cateRoute.Use(middleware.JwtAuth())
	cateRoute.GET("/all", middleware.Role([]string{"user", "author", "admin"}), categoryController.GetCategories)
	cateRoute.Use(middleware.Role([]string{"author", "admin"}))
	cateRoute.POST("/", categoryController.CreateCategory)
	cateRoute.PUT("/:id", categoryController.UpdateCategory)
	cateRoute.DELETE("/:id", categoryController.DeleteCategory)

}
