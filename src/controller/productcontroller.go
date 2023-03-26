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
	GetProduct(c *gin.Context)
}

type controller struct {
	service service.Service
}

func (ctrl *controller) GetProduct(c *gin.Context) {
	boolRes, restErr := ctrl.service.GetProduct()
	if restErr != nil {
		restErr.WrapPath("productcontroller.GetProduct/")
		log.Error(restErr)
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, boolRes)
}
