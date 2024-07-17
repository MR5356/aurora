package host

import (
	"github.com/MR5356/aurora/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary	list container network
// @Tags		container
// @Param		id	path		string	true	"host id"
// @Success	200	{object}	response.Response
// @Router		/host/container/{id}/network [get]
// @Produce	json
func (c *Controller) handleListNetwork(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		if res, err := c.service.ListContainerNetwork(id); err != nil {
			response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		} else {
			response.Success(ctx, res)
		}
	}
}

// @Summary	list container container
// @Tags		container
// @Param		id	path		string	true	"host id"
// @Success	200	{object}	response.Response
// @Router		/host/container/{id}/container [get]
// @Produce	json
func (c *Controller) handleListContainer(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		if res, err := c.service.ListContainer(id); err != nil {
			response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		} else {
			response.Success(ctx, res)
		}
	}
}

// @Summary	list container image
// @Tags		container
// @Param		id	path		string	true	"host id"
// @Success	200	{object}	response.Response
// @Router		/host/container/{id}/image [get]
// @Produce	json
func (c *Controller) handleListImage(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		if res, err := c.service.ListContainerImage(id); err != nil {
			response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		} else {
			response.Success(ctx, res)
		}
	}
}
