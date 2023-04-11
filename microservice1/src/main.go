package main

import (
	"fmt"

	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/conf"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/controllers/productscontroller"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/controllers/todoscontroller"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/httpclients/todosclient"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/log"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/repositories/productsrepo"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/services/productsservice"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/services/todosservice"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func main() {
	// load the env variables
	conf.LoadEnv()

	// instanciate the logger
	log.NewLogger()
	log.Info("main", "main", "starting application")

	// router conf
	gin.SetMode(conf.Env.GinMode)
	router = gin.New()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal(err, "application", "StartApplication", "can't set trusted proxies")
	}
	router.Use(log.Middleware())

	// router url mapping
	productController := productscontroller.NewController(
		productsservice.NewService(
			productsrepo.NewRepo(),
		),
	)
	router.GET("/products/:product-id", productController.GetProduct)

	todosController := todoscontroller.NewController(
		todosservice.NewService(
			todosclient.NewClient(),
		),
	)
	router.GET("/todos", todosController.GetTodos)

	// run launch
	err = router.Run(fmt.Sprintf("%s:%s", conf.Env.Host, conf.Env.Port))
	if err != nil {
		log.Fatal(err, "application", "StartApplication", "can't start web server")
	}
}
