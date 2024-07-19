package host

import (
	"bufio"
	"fmt"

	"github.com/MR5356/aurora/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

// @Summary	get container log
// @Tags		container
// @Param		id	path		string	true	"host id"
// @Param		cid	path		string	true	"container id"
// @Success	200	{object}	response.Response
// @Router		/host/container/{id}/container/{cid}/log [get]
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
			if logs, err := c.service.GetContainerLogs(ctx, id, cid); err != nil {
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
