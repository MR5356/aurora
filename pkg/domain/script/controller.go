package script

import (
	"github.com/MR5356/aurora/pkg/response"
	"github.com/MR5356/aurora/pkg/util/ginutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/script")

	api.POST("", c.handleAddScript)
	api.PUT("", c.handleUpdateScript)
	api.PUT("/batch/delete", c.handleDeleteScript)
	api.GET("/page", c.handlePageScript)
	api.GET("/:id/detail", c.handleDetailScript)
}
