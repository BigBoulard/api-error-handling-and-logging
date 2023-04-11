package todoscontroller

import (
	"net/http"

	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/log"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/services/todosservice"
	"github.com/gin-gonic/gin"
)

func NewController(todosService todosservice.Service) ctrl {
	return &controller{
		todosService: todosService,
	}
}

type ctrl interface {
	GetTodos(c *gin.Context)
}

type controller struct {
	todosService todosservice.Service
}

func (ctrl *controller) GetTodos(c *gin.Context) {
	todos, restErr := ctrl.todosService.GetTodos()
	if restErr != nil {
		restErr.WrapPath("todoscontroller.GetTodos/")
		log.Error(restErr, "can't get todos")
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, todos)
}
