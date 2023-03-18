package controller

import (
	"net/http"
	"strconv"

	"github.com/BigBoulard/api-error-handling-and-logging/src/log"
	"github.com/BigBoulard/api-error-handling-and-logging/src/rest_errors"
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
	DoSomethingWithIntegerParam(c *gin.Context)
}

type controller struct {
	service service.Service
}

func (ctrl *controller) DoSomething(c *gin.Context) {
	boolRes, restErr := ctrl.service.DoSomething()
	if restErr != nil {
		log.Error(restErr, "ctrl", "DoSomething", restErr.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, boolRes)
}

func (ctrl *controller) DoSomethingWithIntegerParam(c *gin.Context) {
	i, err := strconv.ParseInt(c.Param("integer"), 10, 32)
	if err != nil {
		restErr := rest_errors.NewBadRequestError(err.Error())
		log.Error(err, "ctrl", "DoSomethingWithIntegerParam", restErr.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, i)
}
