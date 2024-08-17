package routes

import (
	"github.com/gin-gonic/gin"
)

func InitRoute(r *gin.Engine) {

	api := r.Group("api")
	InitAccountRoute(api)
	InitCategoryRoute(api)

}
