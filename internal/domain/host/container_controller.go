package host

import (
	"bufio"
	"fmt"

	"github.com/MR5356/aurora/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// @Summary	list container network
// @Tags		container
// @Param		id		path		string	true	"host id"
// @Param		driver	path		string	true	"container driver"
// @Success	200		{object}	response.Response
// @Router		/host/container/{id}/{driver}/network [get]
// @Produce	json
func (c *Controller) handleListNetwork(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		if res, err := c.service.ListContainerNetwork(id, ctx.Param("driver")); err != nil {
			response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		} else {
			response.Success(ctx, res)
		}
	}
}

// @Summary	list container container
// @Tags		container
// @Param		id		path		string	true	"host id"
// @Param		driver	path		string	true	"container driver"
// @Success	200		{object}	response.Response
// @Router		/host/container/{id}/{driver}/container [get]
// @Produce	json
func (c *Controller) handleListContainer(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		if res, err := c.service.ListContainer(id, ctx.Param("driver")); err != nil {
			response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		} else {
			response.Success(ctx, res)
		}
	}
}

// @Summary	list container image
// @Tags		container
// @Param		id		path		string	true	"host id"
// @Param		driver	path		string	true	"container driver"
// @Success	200		{object}	response.Response
// @Router		/host/container/{id}/{driver}/image [get]
// @Produce	json
func (c *Controller) handleListImage(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		if res, err := c.service.ListContainerImage(id, ctx.Param("driver")); err != nil {
			response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		} else {
			response.Success(ctx, res)
		}
	}
}

// @Summary	get container log
// @Tags		container
// @Param		id		path		string	true	"host id"
// @Param		cid		path		string	true	"container id"
// @Param		driver	path		string	true	"container driver"
// @Success	200		{object}	response.Response
// @Router		/host/container/{id}/{driver}/container/{cid}/log [get]
// @Produce	json
func (c *Controller) handleGetContainerLogs(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		cid := ctx.Param("cid")
		if len(cid) == 0 {
			response.Error(ctx, response.CodeParamsError)
		} else {
			if logs, err := c.service.GetContainerLogs(ctx, id, cid, ctx.Param("driver")); err != nil {
				response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
			} else {
				defer logs.Close()

				notify := ctx.Writer.CloseNotify()

				scanner := bufio.NewScanner(logs)
				for scanner.Scan() {
					select {
					case <-notify:
						logrus.Infof("Close container logs for %s", cid)
						return
					default:
						line := scanner.Bytes()
						_, err = fmt.Fprintf(ctx.Writer, "data: %s\n\n", string(line))
						if err != nil {
							logrus.Errorf("write response failed, error: %+v", err)
							return
						}
						ctx.Writer.Flush()
					}
				}
			}
		}
	}
}

func (c *Controller) handleExecTerminal(ctx *gin.Context) {
	logrus.Infof("1111")
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		logrus.Infof("123")

		response.Error(ctx, response.CodeParamsError)
	} else {
		cid := ctx.Param("cid")
		if len(cid) == 0 {
			logrus.Infof("124")
			response.Error(ctx, response.CodeParamsError)
		} else {
			user := ctx.Query("user")
			cmd := ctx.Query("cmd")
			if len(cmd) == 0 {
				cmd = "/bin/sh"
			}
			if err := c.service.ExecContainer(ctx, id, cid, ctx.Param("driver"), user, cmd); err != nil {
				logrus.Infof("125")
				response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
			}
		}
	}
}
