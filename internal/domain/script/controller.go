package script

import (
	"github.com/MR5356/aurora/internal/response"
	"github.com/MR5356/aurora/pkg/util/ginutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	return &Controller{
		service: GetService(),
	}
}

// @Summary	add script
// @Tags		script
// @Param		script	body		Script	true	"script info"
// @Success	200		{object}	response.Response
// @Router		/script [post]
// @Produce	json
func (c *Controller) handleAddScript(ctx *gin.Context) {
	script := new(Script)

	if err := ctx.ShouldBindJSON(script); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}
	if err := c.service.AddScript(script); err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

// @Summary	update script
// @Tags		script
// @Param		script	body		Script	true	"script info"
// @Success	200		{object}	response.Response
// @Router		/script [put]
// @Produce	json
func (c *Controller) handleUpdateScript(ctx *gin.Context) {
	script := new(Script)

	if err := ctx.ShouldBindJSON(script); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}
	if err := c.service.UpdateScript(script); err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

// @Summary	delete script
// @Tags		script
// @Param		script	body		[]uuid.UUID	true	"script ids"
// @Success	200		{object}	response.Response
// @Router		/script/batch/delete [put]
// @Produce	json
func (c *Controller) handleDeleteScript(ctx *gin.Context) {
	ids := make([]uuid.UUID, 0)

	if err := ctx.ShouldBindJSON(&ids); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}

	if err := c.service.BatchDeleteScript(ids); err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

// @Summary	page script
// @Tags		script
// @Success	200		{object}	response.Response
// @Param		page	query		int	false	"page number"
// @Param		size	query		int	false	"page size"
// @Router		/script/page [get]
// @Produce	json
func (c *Controller) handlePageScript(ctx *gin.Context) {
	page, size := ginutil.GetPageParams(ctx)

	if res, err := c.service.PageScript(page, size, &Script{}); err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
	} else {
		response.Success(ctx, res)
	}
}

// @Summary	detail script
// @Tags		script
// @Success	200	{object}	response.Response
// @Param		id	path		string	true	"script id"
// @Router		/script/{id}/detail [get]
// @Produce	json
func (c *Controller) handleDetailScript(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		if res, err := c.service.DetailScript(id); err != nil {
			response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		} else {
			response.Success(ctx, res)
		}
	}
}

// @Summary	run script on hosts
// @Tags		script
// @Param		script	body		RunScriptParams	true	"script info"
// @Success	200		{object}	response.Response
// @Router		/script/exec [post]
// @Produce	json
func (c *Controller) handleRunScriptOnHosts(ctx *gin.Context) {
	rsp := new(RunScriptParams)

	if err := ctx.ShouldBindJSON(rsp); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}
	if err := c.service.RunScriptOnHosts(rsp); err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

// @Summary	page script record
// @Tags		script
// @Success	200		{object}	response.Response
// @Param		page	query		int	false	"page number"
// @Param		size	query		int	false	"page size"
// @Router		/script/exec/record/page [get]
// @Produce	json
func (c *Controller) handlePageScriptRecord(ctx *gin.Context) {
	page, size := ginutil.GetPageParams(ctx)
	if res, err := c.service.PageRecord(page, size, &Record{}); err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
	} else {
		response.Success(ctx, res)
	}
}

// @Summary	stop script
// @Tags		script
// @Param		id	path		string	true	"script id"
// @Success	200	{object}	response.Response
// @Router		/script/{id}/stop [delete]
// @Produce	json
func (c *Controller) handleStopScript(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		if err := c.service.StopScript(id); err != nil {
			response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

// @Summary	get job logs
// @Tags		script
// @Param		id	path		string	true	"script id"
// @Success	200	{object}	response.Response
// @Router		/script/exec/{id}/log [get]
// @Produce	json
func (c *Controller) handleGetJobLogs(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		if res, err := c.service.GetJobLog(id); err != nil {
			response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		} else {
			response.Success(ctx, res)
		}
	}
}

// @Summary	get script content
// @Tags		script
// @Success	200	{object}	response.Response
// @Param		id	path		string	true	"script id"
// @Router		/script/{id}/{title} [get]
// @Produce	json
func (c *Controller) handleGetScriptContent(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		if res, err := c.service.GetScriptFile(id); err != nil {
			response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		} else {
			ctx.String(http.StatusOK, res)
		}
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/script")

	api.POST("", c.handleAddScript)
	api.PUT("", c.handleUpdateScript)
	api.PUT("/batch/delete", c.handleDeleteScript)
	api.GET("/page", c.handlePageScript)
	api.GET("/:id/detail", c.handleDetailScript)
	api.GET("/:id/:title", c.handleGetScriptContent)

	exec := api.Group("/exec")
	exec.POST("", c.handleRunScriptOnHosts)
	exec.GET("/record/page", c.handlePageScriptRecord)
	exec.DELETE("/:id", c.handleStopScript)
	exec.GET("/:id/log", c.handleGetJobLogs)

}
