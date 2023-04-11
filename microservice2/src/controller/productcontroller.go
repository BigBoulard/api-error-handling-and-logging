package controller

import (
	"net/http"

	"github.com/BigBoulard/api-error-handling-and-logging/microservice2/src/log"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice2/src/service"
	"github.com/gin-gonic/gin"
)

func NewController(service service.Service) ctrl {
	return &controller{
		service: service,
	}
}

type ctrl interface {
	GetTodos(c *gin.Context)
}

type controller struct {
	service service.Service
}

func (ctrl *controller) GetTodos(c *gin.Context) {
	todos, restErr := ctrl.service.GetTodos()
	if restErr != nil {
		restErr.WrapPath("productcontroller.GetTodos/")
		log.Error(restErr, "productcontroller", "GetTodos", "can't get todos")
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, todos)
}
