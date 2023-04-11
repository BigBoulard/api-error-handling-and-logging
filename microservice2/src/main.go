package main

import (
	"fmt"

	"github.com/BigBoulard/api-error-handling-and-logging/microservice2/src/conf"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice2/src/controller"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice2/src/httpclient"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice2/src/log"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice2/src/service"
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
	gin.SetMode("debug")
	router = gin.New()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal(err, "application", "StartApplication", "can't set trusted proxies")
	}
	router.Use(log.Middleware())

	// router url mapping
	controller := controller.NewController(
		service.NewService(
			httpclient.NewTodosClient(),
		),
	)
	router.GET("/todos", controller.GetTodos)

	// router launch
	err = router.Run(fmt.Sprintf("%s:%s", conf.Env.Host, conf.Env.Port))
	if err != nil {
		log.Fatal(err, "application", "StartApplication", "can't start web server")
	}
}
