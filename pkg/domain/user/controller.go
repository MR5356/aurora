package user

import (
	"github.com/MR5356/aurora/pkg/response"
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

// @Summary	add user
// @Tags		user
// @Param		user	body		User	true	"user info"
// @Success	200		{object}	response.Response
// @Router		/user [post]
// @Produce	json
func (c *Controller) handleAddUser(ctx *gin.Context) {
	user := new(User)
	if err := ctx.ShouldBindJSON(user); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}
	if err := c.service.AddUser(user); err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/user")

	// add user
	api.POST("", c.handleAddUser)
}
