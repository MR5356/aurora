package schedule

import (
	"github.com/MR5356/aurora/pkg/domain/authentication"
	"github.com/MR5356/aurora/pkg/domain/user"
	"github.com/MR5356/aurora/pkg/response"
	"github.com/MR5356/aurora/pkg/server/ginmiddleware/datafilter"
	"github.com/MR5356/aurora/pkg/util/ginutil"
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

// @Summary	add schedule
// @Tags		schedule
// @Param		schedule	body		Schedule	true	"schedule info"
// @Success	200			{object}	response.Response
// @Router		/schedule [post]
// @Produce	json
func (c *Controller) handleAddSchedule(ctx *gin.Context) {
	schedule := new(Schedule)
	if err := ctx.ShouldBindJSON(schedule); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}
	if err := c.service.AddSchedule(schedule); err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
	} else {
		u, err := user.GetJWTService().ParseToken(ginutil.GetToken(ctx))
		if err != nil {
			logrus.Errorf("parse token failed, error: %v", err)
			response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
			return
		}
		if ok, err := authentication.GetPermission().AddPolicyForRoleInDomain(AuthDomain, u.ID, schedule.ID.String(), ActionOwner); err != nil || !ok {
			logrus.Errorf("add policy for role in domain failed, error: %v", err)
			response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
			return
		}
		response.Success(ctx, nil)
	}
}

// @Summary	update schedule
// @Tags		schedule
// @Param		schedule	body		Schedule	true	"schedule info"
// @Success	200			{object}	response.Response
// @Router		/schedule/{id} [put]
// @Produce	json
func (c *Controller) handleUpdateSchedule(ctx *gin.Context) {
	schedule := new(Schedule)
	if err := ctx.ShouldBindJSON(schedule); err != nil {
		logrus.Errorf("bind json failed, error: %v", err)
		response.Error(ctx, response.CodeParamsError)
		return
	}
	if err := c.service.UpdateSchedule(schedule); err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

// @Summary	page schedule
// @Tags		schedule
// @Success	200		{object}	response.Response
// @Param		page	query		int	false	"page number"
// @Param		size	query		int	false	"size number"
// @Router		/schedule/page [get]
// @Produce	json
func (c *Controller) handlePageSchedule(ctx *gin.Context) {
	page, size := ginutil.GetPageParams(ctx)

	if res, err := c.service.PageSchedule(page, size, &Schedule{}); err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, res)
	}
}

// @Summary	page schedule record
// @Tags		schedule
// @Success	200		{object}	response.Response
// @Param		page	query		int	false	"page number"
// @Param		size	query		int	false	"size number"
// @Router		/schedule/record/page [get]
// @Produce	json
func (c *Controller) handlePageScheduleRecord(ctx *gin.Context) {
	page, size := ginutil.GetPageParams(ctx)

	if res, err := c.service.PageScheduleRecord(page, size, &Record{}); err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, res)
	}
}

// @Summary	detail schedule
// @Tags		schedule
// @Param		id	path		string	true	"schedule id"
// @Success	200	{object}	response.Response{data=Schedule}
// @Router		/schedule/{id}/detail [get]
// @Produce	json
func (c *Controller) handleDetailSchedule(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		if res, err := c.service.DetailSchedule(id); err != nil {
			response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		} else {
			response.Success(ctx, res)
		}
	}
}

// @Summary	delete schedule
// @Tags		schedule
// @Param		id	path		string	true	"schedule id"
// @Success	200	{object}	response.Response
// @Router		/schedule/{id} [delete]
// @Produce	json
func (c *Controller) handleDeleteSchedule(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		if err := c.service.DeleteSchedule(id); err != nil {
			response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

// @Summary	get task executors
// @Tags		schedule
// @Success	200	{object}	response.Response
// @Router		/schedule/executors [get]
// @Produce	json
func (c *Controller) handleGetTaskExecutors(ctx *gin.Context) {
	response.Success(ctx, c.service.GetTaskExecutors())
}

// @Summary	list schedule
// @Tags		schedule
// @Success	200	{object}	response.Response{data=[]Schedule}
// @Router		/schedule/list [get]
// @Produce	json
func (c *Controller) handleListSchedule(ctx *gin.Context) {
	res, err := c.service.scheduleDB.List(&Schedule{})
	if err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, res)
	}
}

// @Summary	batch delete schedule
// @Tags		schedule
// @Param		ids	body		[]uuid.UUID	true	"schedule ids"
// @Success	200	{object}	response.Response
// @Router		/schedule/batch/delete [put]
// @Produce	json
func (c *Controller) handleBatchDeleteSchedule(ctx *gin.Context) {
	ids := make([]uuid.UUID, 0)

	if err := ctx.ShouldBindJSON(&ids); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}
	if err := c.service.BatchDeleteSchedule(ids); err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

// @Summary	batch enable schedule
// @Tags		schedule
// @Param		ids	body		[]uuid.UUID	true	"schedule ids"
// @Success	200	{object}	response.Response
// @Router		/schedule/batch/enable [delete]
// @Produce	json
func (c *Controller) handleBatchEnableSchedule(ctx *gin.Context) {
	ids := make([]uuid.UUID, 0)

	if err := ctx.ShouldBindJSON(&ids); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}
	if err := c.service.BatchSetScheduleEnable(ids, true); err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

// @Summary	batch enable schedule
// @Tags		schedule
// @Param		ids	body		[]uuid.UUID	true	"schedule ids"
// @Success	200	{object}	response.Response
// @Router		/schedule/batch/disable [delete]
// @Produce	json
func (c *Controller) handleBatchDisableSchedule(ctx *gin.Context) {
	ids := make([]uuid.UUID, 0)

	if err := ctx.ShouldBindJSON(&ids); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}
	if err := c.service.BatchSetScheduleEnable(ids, false); err != nil {
		response.ErrorWithMsg(ctx, response.CodeServerError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/schedule")

	// FIXME: disable global filter
	//api.Use(datafilter.AutomationFilter())
	datafilter.RegisterFilter([]datafilter.Filter{
		{
			Function: c.handleDeleteSchedule,
			IsBefore: true,
			Action:   []string{ActionAdmin, ActionOwner},
			Domain:   AuthDomain,
		},
		{
			Function: c.handleBatchDeleteSchedule,
			IsBefore: true,
			Action:   []string{ActionAdmin, ActionOwner},
			Domain:   AuthDomain,
		},
		{
			Function: c.handleBatchEnableSchedule,
			IsBefore: true,
			Action:   []string{ActionAdmin, ActionOwner},
			Domain:   AuthDomain,
		},
		{
			Function: c.handleBatchDisableSchedule,
			IsBefore: true,
			Action:   []string{ActionAdmin, ActionOwner},
			Domain:   AuthDomain,
		},
		{
			Function: c.handleUpdateSchedule,
			IsBefore: true,
			Action:   []string{ActionAdmin, ActionOwner},
			Domain:   AuthDomain,
		},
		{
			Function: c.handleDetailSchedule,
			IsBefore: true,
			Action:   []string{ActionAdmin, ActionOwner, ActionUser},
			Domain:   AuthDomain,
		},
		{
			Function: c.handlePageSchedule,
			IsBefore: false,
			Action:   []string{ActionAdmin, ActionOwner, ActionUser},
			Domain:   AuthDomain,
		},
		{
			Function: c.handleListSchedule,
			IsBefore: false,
			Action:   []string{ActionAdmin, ActionOwner, ActionUser},
			Domain:   AuthDomain,
		},
		//{
		//	Function: c.handlePageScheduleRecord,
		//	IsBefore: false,
		//	Action:   []string{ActionAdmin, ActionOwner, ActionUser},
		//	Domain:   AuthDomain,
		//},
	})

	// list schedule
	api.GET("/list", c.handleListSchedule)

	// page schedule
	api.GET("/page", c.handlePageSchedule)

	// detail schedule
	api.GET("/:id/detail", c.handleDetailSchedule)

	// add schedule
	api.POST("", c.handleAddSchedule)

	// update schedule
	api.PUT("/:id", c.handleUpdateSchedule)

	// delete schedule
	api.DELETE("/:id", c.handleDeleteSchedule)

	// batch delete schedule
	api.PUT("/batch/delete", c.handleBatchDeleteSchedule)

	// batch enable schedule
	api.PUT("/batch/enable", c.handleBatchEnableSchedule)

	// batch disable schedule
	api.PUT("/batch/disable", c.handleBatchDisableSchedule)

	// page schedule record
	api.GET("/record/page", c.handlePageScheduleRecord)

	// get task executors
	api.GET("/executors", c.handleGetTaskExecutors)
}
