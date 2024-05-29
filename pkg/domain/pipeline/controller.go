package pipeline

import (
	"github.com/MR5356/aurora/pkg/response"
	"github.com/gin-gonic/gin"
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

// @Summary	add workflow
// @Tags		pipeline
// @Param		workflow	body		WorkflowRequest	true	"workflow info"
// @Success	200			{object}	response.Response
// @Router		/pipeline/workflow [post]
// @Produce	json
func (c *Controller) handleAddWorkflow(ctx *gin.Context) {
	wfr := new(WorkflowRequest)
	if err := ctx.ShouldBindJSON(wfr); err != nil {
		logrus.Errorf("bind json failed, error: %v", err)
		response.Error(ctx, response.CodeParamsError)
		return
	}
	if err := c.service.AddWorkflow(wfr); err != nil {
		logrus.Errorf("add workflow failed, error: %v", err)
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

// @Summary	add workflow
// @Tags		pipeline
// @Param		workflow	body		WorkflowRequest	true	"workflow info"
// @Success	200			{object}	response.Response
// @Router		/pipeline/workflow [post]
// @Produce	json
func (c *Controller) handleUpdateWorkflow(ctx *gin.Context) {
	wfr := new(WorkflowRequest)
	if err := ctx.ShouldBindJSON(wfr); err != nil {
		logrus.Errorf("bind json failed, error: %v", err)
		response.Error(ctx, response.CodeParamsError)
		return
	}
	if err := c.service.UpdateWorkflow(wfr); err != nil {
		logrus.Errorf("add workflow failed, error: %v", err)
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

// @Summary	list workflow
// @Tags		pipeline
// @Success	200	{object}	response.Response{data=[]Workflow}
// @Router		/pipeline/workflow [get]
// @Produce	json
func (c *Controller) handleListWorkflow(ctx *gin.Context) {
	res, err := c.service.ListWorkflow(&Workflow{})
	if err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, res)
	}
}

// @Summary	get workflow
// @Tags		pipeline
// @Param		id	path		string	true	"workflow id"
// @Success	200	{object}	response.Response{data=Workflow}
// @Router		/pipeline/workflow/{id} [get]
// @Produce	json
func (c *Controller) handleGetWorkflow(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := c.service.GetWorkflow(&Workflow{ID: id})
	if err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, res)
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/pipeline")
	api.GET("/workflow", c.handleListWorkflow)
	api.POST("/workflow", c.handleAddWorkflow)
	api.PUT("/workflow", c.handleUpdateWorkflow)
	api.GET("/workflow/:id", c.handleGetWorkflow)
}
