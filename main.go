package main

import (
	"dlms/configs"
	"dlms/routes"
	"dlms/utils"
	"fmt"
	"os"

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

	// defer func() {
	// 	if err := client.Disconnect(context.Background()); err != nil {
	// 		panic(err)
	// 	} else {
	// 		fmt.Println("db disconnect")
	// 	}
	// }()

	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("eqfield", utils.CustomEqField)
	}

	routes.InitRoute(r)

	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))
	r.Run(PORT)
}
