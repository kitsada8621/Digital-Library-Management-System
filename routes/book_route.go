package routes

import (
	"dlms/controllers"
	"dlms/middleware"

	"github.com/gin-gonic/gin"
)

func IntiBookRoute(r *gin.RouterGroup) {

	bookController := controllers.NewBookController()
	bookRoute := r.Group("book")
	bookRoute.Use(middleware.JwtAuth())
	bookRoute.POST("/most/borrowed", bookController.MostBorrowedBooks)
	bookRoute.POST("/all", middleware.Role([]string{"admin", "user", "author"}), bookController.GetBooks)
	bookRoute.POST("/create", middleware.Role([]string{"author"}), bookController.CreateBook)
	bookRoute.GET("/:id", middleware.Role([]string{"admin", "user", "author"}), bookController.BookDetails)
	bookRoute.GET("/:id/edit", middleware.Role([]string{"author"}), bookController.EditBook)
	bookRoute.PUT("/update/:id", middleware.Role([]string{"author"}), bookController.UpdateBook)
	bookRoute.DELETE("/delete/:id", middleware.Role([]string{"admin", "author"}), bookController.DeleteBook)

}
