package health

import (
	"encoding/json"
	"fmt"
	"github.com/MR5356/aurora/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

var (
	statisticClientCount = atomic.Int64{}
	statistic            *Statistics
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	c := &Controller{
		service: GetService(),
	}

	go c.cacheStatistics()
	return c
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
		st := ctx.Query("startTime")
		if startTime, err = time.Parse(time.RFC3339, st); err != nil {
			logrus.Errorf("bind query failed, error: %v", err)
			startTime = time.Now().Add(-10 * time.Minute)
		}
		et := ctx.Query("endTime")
		if endTime, err = time.Parse(time.RFC3339, et); err != nil {
			logrus.Errorf("bind query failed, error: %v", err)
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

// @Summary	get statistics
// @Tags		health
// @Success	200	{object}	response.Response{data=Statistics}
// @Router		/health/statistics [get]
// @Produce	json
func (c *Controller) handleGetStatistics(ctx *gin.Context) {
	if res, err := c.service.HealthStatistics(); err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, res)
	}
}

func (c *Controller) cacheStatistics() {
	for range time.Tick(time.Second) {
		if statisticClientCount.Load() > 0 {
			if res, err := c.service.HealthStatistics(); err != nil {
				logrus.Errorf("cache statistics failed, error: %v", err)
			} else {
				statistic = res
			}
		}
	}
}

// @Summary	get statistics with SSE
// @Tags		health
// @Success	200	{object}	response.Response
// @Router		/health/statistics/sse [get]
// @Produce	json
func (c *Controller) handleGetStatisticsWithSSE(ctx *gin.Context) {
	statisticClientCount.Add(1)
	defer statisticClientCount.Add(-1)
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	notify := ctx.Writer.CloseNotify()

	// send statistic right away
	resStr, _ := json.Marshal(statistic)
	_, err := fmt.Fprintf(ctx.Writer, "data: %s\n\n", resStr)
	if err != nil {
		logrus.Errorf("write response failed, error: %v", err)
		return
	}
	ctx.Writer.Flush()

	t := time.Tick(time.Second)

	for {
		select {
		case <-notify:
			logrus.Infof("Close SSE connection")
			return
		case <-t:
			resStr, _ = json.Marshal(statistic)
			_, err = fmt.Fprintf(ctx.Writer, "data: %s\n\n", resStr)
			if err != nil {
				logrus.Errorf("write response failed, error: %v", err)
				return
			}
			ctx.Writer.Flush()
		}
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/health")

	api.GET("/list", c.handleListHealth)
	api.POST("/add", c.handleAddHealth)
	api.PUT("/:id", c.handleUpdateHealth)
	api.DELETE("/:id", c.handleDeleteHealth)
	api.GET("/:id/detail", c.handleDetailHealth)
	api.GET("/types", c.handleGetHealthCheckTypes)
	api.GET("/statistics", c.handleGetStatistics)
	api.GET("/statistics/sse", c.handleGetStatisticsWithSSE)

	api.GET("/:id/record", c.handleGetTimeRangeRecord)
}
