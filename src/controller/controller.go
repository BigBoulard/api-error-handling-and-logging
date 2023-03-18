package controller

import (
	"net/http"

	"github.com/BigBoulard/api-error-handling-and-logging/src/log"
	"github.com/BigBoulard/api-error-handling-and-logging/src/service"
	"github.com/gin-gonic/gin"
)

func NewController(service service.Service) ctrl {
	return &controller{
		service: service,
	}
}

type ctrl interface {
	DoSomething(c *gin.Context)
}

type controller struct {
	service service.Service
}

func (ctrl *controller) DoSomething(c *gin.Context) {
	boolRes, restErr := ctrl.service.DoSomething()
	if restErr != nil {
		log.Error(restErr, "ctrl", "DoSomething", restErr.Causes())
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, boolRes)
}
