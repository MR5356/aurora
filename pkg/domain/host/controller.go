package host

import (
	"github.com/MR5356/aurora/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	return &Controller{
		service: GetService(),
	}
}

// @Summary	list host group
// @Tags		host
// @Success	200	{object}	response.Response{data=[]Group}
// @Router		/host/group/list [get]
// @Produce	json
func (c *Controller) handleListGroup(ctx *gin.Context) {
	if res, err := c.service.ListGroup(&Group{}); err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, res)
	}
}

// @Summary	add host group
// @Tags		host
// @Param		group	body		Group	true	"group info"
// @Success	200		{object}	response.Response
// @Router		/host/group/add [post]
// @Produce	json
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

// @Summary	update host group
// @Tags		host
// @Param		group	body		Group	true	"group info"
// @Param		id		path		string	true	"group id"
// @Success	200		{object}	response.Response
// @Router		/host/group/{id} [put]
// @Produce	json
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

// @Summary	delete host group
// @Tags		host
// @Param		id	path		string	true	"group id"
// @Success	200	{object}	response.Response
// @Router		/host/group/{id} [delete]
// @Produce	json
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

// @Summary	add host
// @Tags		host
// @Param		host	body		Host	true	"host info"
// @Success	200		{object}	response.Response
// @Router		/host/add [post]
// @Produce	json
func (c *Controller) handleAddHost(ctx *gin.Context) {
	host := new(Host)

	if err := ctx.ShouldBindJSON(host); err != nil {
		logrus.Errorf("bind json failed, error: %v", err)
		response.Error(ctx, response.CodeParamsError)
		return
	} else {
		if err := c.service.AddHost(host); err != nil {
			response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

// @Summary	delete host
// @Tags		host
// @Param		id	path		string	true	"host id"
// @Success	200	{object}	response.Response
// @Router		/host/{id} [delete]
// @Produce	json
func (c *Controller) handleDeleteHost(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		if err := c.service.DeleteHost(id); err != nil {
			response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

// @Summary	update host
// @Tags		host
// @Param		host	body		Host	true	"host info"
// @Param		id		path		string	true	"host id"
// @Success	200		{object}	response.Response
// @Router		/host/{id} [put]
// @Produce	json
func (c *Controller) handleUpdateHost(ctx *gin.Context) {
	host := new(Host)

	if err := ctx.ShouldBindJSON(host); err != nil {
		logrus.Errorf("bind json failed, error: %v", err)
		response.Error(ctx, response.CodeParamsError)
		return
	} else {
		if id, err := uuid.Parse(ctx.Param("id")); err != nil {
			response.Error(ctx, response.CodeParamsError)
		} else {
			host.ID = id
			if err := c.service.UpdateHost(host); err != nil {
				response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
			} else {
				response.Success(ctx, nil)
			}
		}
	}
}

// @Summary	list host
// @Tags		host
// @Success	200			{object}	response.Response{data=[]Host}
// @Param		group_id	query		string	false	"group id"
// @Router		/host/list [get]
// @Produce	json
func (c *Controller) handleListHost(ctx *gin.Context) {
	groupId, err := uuid.Parse(ctx.Query("group_id"))
	host := &Host{}
	if err == nil {
		host.GroupId = groupId
	}
	res, err := c.service.ListHost(host)
	if err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, res)
	}
}

// @Summary	detail host
// @Tags		host
// @Param		id	path		string	true	"host id"
// @Success	200	{object}	response.Response{data=Host}
// @Router		/host/{id}/detail [get]
// @Produce	json
func (c *Controller) handleDetailHost(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		if res, err := c.service.DetailHost(id); err != nil {
			response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
		} else {
			response.Success(ctx, res)
		}
	}
}

func (c *Controller) RegisterRoute(engine *gin.RouterGroup) {
	host := engine.Group("host")
	host.POST("/add", c.handleAddHost)
	host.DELETE("/:id", c.handleDeleteHost)
	host.PUT("/:id", c.handleUpdateHost)
	host.GET("/list", c.handleListHost)
	host.GET("/:id/detail", c.handleDetailHost)

	group := host.Group("group")
	group.GET("/list", c.handleListGroup)
	group.POST("/add", c.handleAddGroup)
	group.PUT("/:id", c.handleUpdateGroup)
	group.DELETE("/:id", c.handleDeleteGroup)

	term := host.Group("terminal")
	term.GET(":id", c.handleTerminal)

	container := host.Group("container")
	container.GET("/:id/network", c.handleListNetwork)
	container.GET("/:id/image", c.handleListImage)
	container.GET("/:id/container", c.handleListContainer)
	container.GET("/:id/container/:cid/log", c.handleGetContainerLogs)
}
