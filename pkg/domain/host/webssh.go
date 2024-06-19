package host

import (
	"github.com/MR5356/aurora/pkg/response"
	"github.com/MR5356/aurora/pkg/util/sshutil"
	"github.com/MR5356/aurora/pkg/util/terminal"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (c *Controller) handleTerminal(ctx *gin.Context) {
	if id, err := uuid.Parse(ctx.Param("id")); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		host, err := c.service.DetailHost(id)
		if err != nil {
			response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
			return
		}

		sshClient, err := sshutil.NewSSHClient(host.HostInfo)
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
}
