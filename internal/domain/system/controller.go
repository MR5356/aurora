package system

import (
	"github.com/MR5356/aurora/internal/response"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	return &Controller{
		service: GetService(),
	}
}

// @Summary	get statistic
// @Tags		system
// @Success	200	{object}	response.Response{data=[]Statistic}
// @Router		/system/statistic [get]
// @Produce	json
func (c *Controller) handleGetStatistic(ctx *gin.Context) {
	res, err := c.service.GetStatistic()
	if err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, res)
	}
}

func (c *Controller) handleGetVersion(ctx *gin.Context) {
	response.Success(ctx, c.service.GetVersionInfo())
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/system")

	api.GET("/statistic", c.handleGetStatistic)
	api.GET("/version", c.handleGetVersion)
}
