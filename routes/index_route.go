package routes

import (
	"github.com/gin-gonic/gin"
)

func InitRoute(r *gin.Engine) {

	api := r.Group("api")
	InitAccountRoute(api)
	InitCategoryRoute(api)
	IntiBookRoute(api)
	InitBorrowRoute(api)
	InitReturnBookRoute(api)
	InitUserRoute(api)
	InitRoleRoute(api)
}
