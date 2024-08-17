package routes

import (
	"dlms/controllers"
	"dlms/middleware"

	"github.com/gin-gonic/gin"
)

func InitCategoryRoute(r *gin.RouterGroup) {
	categoryController := controllers.NewCategoryController()

	r.GET("/categories", middleware.Authorize(), categoryController.GetCategories)

	api := r.Group("category")
	api.POST("/", middleware.Authorize(), categoryController.CreateCategory)
	api.PUT("/:id", middleware.Authorize(), categoryController.UpdateCategory)
	api.DELETE("/:id", middleware.Authorize(), categoryController.DeleteCategory)

}
