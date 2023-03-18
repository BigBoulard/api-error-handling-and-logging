package application

import (
	"github.com/BigBoulard/api-error-handling-and-logging/src/controller"
	"github.com/BigBoulard/api-error-handling-and-logging/src/httpclient"
	"github.com/BigBoulard/api-error-handling-and-logging/src/service"
)

func mapUrls() {
	controller := controller.NewController(
		service.NewService(
			httpclient.NewClient(),
		),
	)
	router.POST("/do-something", controller.DoSomething)
}
