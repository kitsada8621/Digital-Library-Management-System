package routes

import (
	"dlms/controllers"
	"dlms/middleware"

	"github.com/gin-gonic/gin"
)

func InitBorrowRoute(r *gin.RouterGroup) {

	borrowController := controllers.NewBorrowController()
	borrowRoute := r.Group("borrow")
	borrowRoute.Use(middleware.JwtAuth())

	borrowRoute.POST("/create", middleware.Role([]string{"user"}), borrowController.BorrowBook)
	borrowRoute.PUT("/update/:id", middleware.Role([]string{"user"}), borrowController.UpdateBorrow)
	borrowRoute.PUT("/cancel/:id", middleware.Role([]string{"user"}), borrowController.CancelBorrowBook)
	borrowRoute.GET("/history", middleware.Role([]string{"user"}), borrowController.BorrowingHistory)

	borrowRoute.Use(middleware.Role([]string{"admin"}))
	borrowRoute.GET("/all", borrowController.GetBorrows)
	borrowRoute.GET("/:id", borrowController.BorrowDetails)
	borrowRoute.DELETE("/delete/:id", borrowController.DeleteBorrow)
	borrowRoute.PUT("/approve/:id", borrowController.ApproveBookBorrowing)

}
