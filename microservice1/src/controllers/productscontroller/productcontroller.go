package productscontroller

import (
	"net/http"

	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/log"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/services/productsservice"
	"github.com/gin-gonic/gin"
)

func NewController(productsService productsservice.Service) ctrl {
	return &controller{
		productsService: productsService,
	}
}

type ctrl interface {
	GetProduct(c *gin.Context)
}

type controller struct {
	productsService productsservice.Service
}

func (ctrl *controller) GetProduct(c *gin.Context) {
	boolRes, restErr := ctrl.productsService.GetProduct()
	if restErr != nil {
		restErr.WrapPath("productcontroller.GetProduct/")
		log.Error(restErr, "can't get product")
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, boolRes)
}
