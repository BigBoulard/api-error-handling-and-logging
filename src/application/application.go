package application

import (
	"log"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func StartApplication() {
	gin.SetMode("debug")

	router = gin.Default()
	err := router.SetTrustedProxies(nil) // Context.ClientIP() returns the client IP directly
	if err != nil {
		log.Fatal(err, "gw", "app", "router.SetTrustedProxies(nil)")
	}

	mapUrls()
	err = router.Run("localhost:9090")
	if err != nil {
		log.Fatal("Unable to start Gin server")
	}
}
