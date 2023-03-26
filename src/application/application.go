package application

import (
	"fmt"

	"github.com/BigBoulard/api-error-handling-and-logging/src/conf"
	"github.com/BigBoulard/api-error-handling-and-logging/src/log"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func StartApplication() {
	// load the env variables
	conf.LoadEnv()

	gin.SetMode("debug")

	router = gin.Default()
	err := router.SetTrustedProxies(nil)

	router.Use(log.Middleware())

	if err != nil {
		log.Fatal("application/main", err)
	}

	mapUrls()
	err = router.Run(fmt.Sprintf("%s:%s", conf.Env.Host, conf.Env.Port))
	if err != nil {
		log.Fatal("application/main", err)
	}
}
