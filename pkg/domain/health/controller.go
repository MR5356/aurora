package health

import (
	"github.com/MR5356/aurora/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	return &Controller{
		service: GetService(),
	}
}

// @Summary	list health
// @Tags		health
// @Success	200	{object}	response.Response{data=[]Health}
// @Router		/health/list [get]
// @Produce	json
func (c *Controller) handleListHealth(ctx *gin.Context) {
	if res, err := c.service.ListHealth(&Health{}); err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, res)
	}
}

// @Summary	add health
// @Tags		health
// @Param		health	body		Health	true	"health info"
// @Success	200		{object}	response.Response
// @Router		/health/add [post]
// @Produce	json
func (c *Controller) handleAddHealth(ctx *gin.Context) {
	health := new(Health)
	if err := ctx.ShouldBindJSON(health); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}
	if err := c.service.AddHealth(health); err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

// @Summary	update health
// @Tags		health
// @Param		health	body		Health	true	"health info"
// @Param		id		path		string	true	"health id"
// @Success	200		{object}	response.Response
// @Router		/health/{id} [put]
// @Produce	json
func (c *Controller) handleUpdateHealth(ctx *gin.Context) {
	health := new(Health)
	if err := ctx.ShouldBindJSON(health); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}

	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	} else {
		health.ID = id
	}

	if err := c.service.UpdateHealth(health); err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

// @Summary	delete health
// @Tags		health
// @Param		id	path		string	true	"health id"
// @Success	200	{object}	response.Response
// @Router		/health/{id} [delete]
// @Produce	json
func (c *Controller) handleDeleteHealth(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	} else {
		if err := c.service.DeleteHealth(&Health{ID: id}); err != nil {
			response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

// @Summary	detail health
// @Tags		health
// @Param		id	path		string	true	"health id"
// @Success	200	{object}	response.Response{data=Health}
// @Router		/health/{id} [get]
// @Produce	json
func (c *Controller) handleDetailHealth(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	} else {
		if res, err := c.service.DetailHealth(id); err != nil {
			response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
		} else {
			response.Success(ctx, res)
		}
	}
}

// @Summary	get time range record
// @Tags		health
// @Param		id			path		string	true	"health id"
// @Param		startTime	query		string	false	"start time"
// @Param		endTime		query		string	false	"end time"
// @Success	200			{object}	response.Response{data=[]Record}
// @Router		/health/{id}/record [get]
// @Produce	json
func (c *Controller) handleGetTimeRangeRecord(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	} else {
		var startTime, endTime time.Time
		if err := ctx.ShouldBindQuery(&startTime); err != nil {
			startTime = time.Now().Add(-10 * time.Minute)
		}
		if err := ctx.ShouldBindQuery(&endTime); err != nil {
			endTime = time.Now()
		}
		logrus.Infof("startTime: %v, endTime: %v", startTime, endTime)
		if res, err := c.service.GetTimeRangeRecord(id, startTime, endTime); err != nil {
			response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
		} else {
			response.Success(ctx, res)
		}
	}
}

// @Summary	get health check types
// @Tags		health
// @Success	200	{object}	response.Response{data=[]CheckType}
// @Router		/health/types [get]
// @Produce	json
func (c *Controller) handleGetHealthCheckTypes(ctx *gin.Context) {
	response.Success(ctx, c.service.GetHealthCheckTypes())
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/health")

	api.GET("/list", c.handleListHealth)
	api.POST("/add", c.handleAddHealth)
	api.PUT("/:id", c.handleUpdateHealth)
	api.DELETE("/:id", c.handleDeleteHealth)
	api.GET("/:id/detail", c.handleDetailHealth)
	api.GET("/types", c.handleGetHealthCheckTypes)

	api.GET("/:id/record", c.handleGetTimeRangeRecord)
}
