package plugin

import (
	"github.com/gin-gonic/gin"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	_ = group.Group("/plugin")

}
