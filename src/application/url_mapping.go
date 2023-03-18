package application

import (
	"github.com/BigBoulard/api-error-handling-and-logging/src/controller"
	"github.com/BigBoulard/api-error-handling-and-logging/src/dbclient"
	"github.com/BigBoulard/api-error-handling-and-logging/src/service"
)

func mapUrls() {
	controller := controller.NewController(
		service.NewService(
			dbclient.NewClient(),
		),
	)
	router.POST("/do-something", controller.DoSomething) // return a 404 from the persistance layer

	router.POST("/do-something-with-integer-param/:integer", controller.DoSomethingWithIntegerParam) // return a 400 from the controller (parse error)

}
