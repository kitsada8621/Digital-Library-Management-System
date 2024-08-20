package main

import (
	"dlms/configs"
	"dlms/database"
	"dlms/pkg/utils"
	"dlms/routes"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {

	err := configs.InitDotEnv()
	if err != nil {
		panic(err)
	}

	_, err = configs.ConnectDB()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	r.Use(gin.Recovery())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("eqfield", utils.CustomEqField)
	}

	if err := database.EnsureSeederRoleData(); err != nil {
		panic(err.Error())
	}

	if err := database.EnsureSeederAdminData(); err != nil {
		panic(err.Error())
	}

	routes.InitRoute(r)

	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))
	r.Run(PORT)
}
