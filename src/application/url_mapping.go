package application

import (
	"github.com/BigBoulard/api-error-handling-and-logging/src/controller"
	"github.com/BigBoulard/api-error-handling-and-logging/src/repository"
	"github.com/BigBoulard/api-error-handling-and-logging/src/service"
)

func mapUrls() {
	controller := controller.NewController(
		service.NewService(
			repository.NewRepo(),
		),
	)
	router.GET("/products/:product-id", controller.GetProduct) // return a 404 from the persistance layer
}
