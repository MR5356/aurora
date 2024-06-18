package host

import (
	"github.com/MR5356/aurora/pkg/response"
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

//	@Summary	list host group
//	@Tags		host
//	@Success	200	{object}	response.Response{data=[]Group}
//	@Router		/host/group/list [get]
//	@Produce	json
func (c *Controller) handleListGroup(ctx *gin.Context) {
	if res, err := c.service.ListGroup(&Group{}); err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, res)
	}
}

//	@Summary	add host group
//	@Tags		host
//	@Param		group	body		Group	true	"group info"
//	@Success	200		{object}	response.Response
//	@Router		/host/group/add [post]
//	@Produce	json
func (c *Controller) handleAddGroup(ctx *gin.Context) {
	group := new(Group)
	if err := ctx.ShouldBindJSON(group); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}
	if err := c.service.AddGroup(group); err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

//	@Summary	update host group
//	@Tags		host
//	@Param		group	body		Group	true	"group info"
//	@Param		id		path		string	true	"group id"
//	@Success	200		{object}	response.Response
//	@Router		/host/group/{id} [put]
//	@Produce	json
func (c *Controller) handleUpdateGroup(ctx *gin.Context) {
	group := new(Group)
	if err := ctx.ShouldBindJSON(group); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}

	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	} else {
		group.ID = id
	}

	if err := c.service.UpdateGroup(group); err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

//	@Summary	delete host group
//	@Tags		host
//	@Param		id	path		string	true	"group id"
//	@Success	200	{object}	response.Response
//	@Router		/host/group/{id} [delete]
//	@Produce	json
func (c *Controller) handleDeleteGroup(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		if err := c.service.DeleteGroup(id); err != nil {
			response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

func (c *Controller) RegisterRoute(engine *gin.RouterGroup) {
	host := engine.Group("host")

	group := host.Group("group")
	group.GET("/list", c.handleListGroup)
	group.POST("/add", c.handleAddGroup)
	group.PUT("/:id", c.handleUpdateGroup)
	group.DELETE("/:id", c.handleDeleteGroup)
}
