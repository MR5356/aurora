package host

import (
	"github.com/MR5356/aurora/pkg/response"
	"github.com/MR5356/aurora/pkg/util/sshutil"
	"github.com/MR5356/aurora/pkg/util/terminal"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"net/http"
)

func (c *Controller) handleTerminal(ctx *gin.Context) {
	var hostInfo sshutil.HostInfo
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		hostInfo.Host = ctx.Query("host")
		hostInfo.Port = cast.ToUint16(ctx.Query("port"))
		hostInfo.Username = ctx.Query("username")
		hostInfo.Password = ctx.Query("password")
	} else {
		host, err := c.service.DetailHost(id)
		if err != nil {
			response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
			return
		}
		hostInfo = host.HostInfo
	}

	sshClient, err := sshutil.NewSSHClient(hostInfo)
	if err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		return
	}

	t := terminal.NewTerminal()
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{"webssh"},
	}

	webConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		return
	}

	webConn.SetCloseHandler(t.CloseHandler)

	t.Websocket = webConn
	t.Session, err = sshClient.GetSession()
	if err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		return
	}

	t.Stdin, err = t.Session.StdinPipe()
	if err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		return
	}

	sshOut := new(terminal.WsBufferWriter)
	t.Session.Stdout = sshOut
	t.Session.Stderr = sshOut
	t.Stdout = sshOut

	if err := t.Session.RequestPty("xterm-256color", 30, 120, terminal.Modes); err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		return
	}

	err = t.Session.Shell()
	if err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		return
	}

	go func() {
		if err := t.Session.Wait(); err != nil {
			logrus.Errorf("session wait error: %v", err)
		} else {
			logrus.Infof("session wait done")
		}
		logrus.Errorf("close websocket")
		_ = t.Websocket.Close()
		return
	}()
	go t.Send2SSH()
	go t.Send2Web()

}
