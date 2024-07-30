package docker

import (
	"io"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (c *Client) Terminal(ctx *gin.Context, containerId string, user, cmd string) error {
	logrus.Infof("4")
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{"webssh"},
	}
	logrus.Infof("5")

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.Errorf("upgrade error: %+v", err)
		return err
	}

	defer ws.Close()
	logrus.Infof("6")

	execIDResp, err := c.client.ContainerExecCreate(ctx, containerId, container.ExecOptions{
		User: user,
		Cmd: []string{cmd},
		AttachStdin: true,
		AttachStderr: true,
		AttachStdout: true,
	})

	if err != nil {
		logrus.Errorf("container exec create error: %+v", err)
		return err
	}

	logrus.Infof("1")

	hijackedResp, err := c.client.ContainerExecAttach(ctx, execIDResp.ID, container.ExecAttachOptions{
		Tty: true,
	})

	if err != nil {
		logrus.Errorf("container exec attach error: %+v", err)
		return err
	}
	defer hijackedResp.Close()
	logrus.Infof("2")

	done := make(chan struct{})

	go func() {
		_, err := io.Copy(hijackedResp.Conn, ws.UnderlyingConn())
		if err != nil {
			logrus.Errorf("failed to copy data to exec stdin: %+v", err)
		}
	}()

	go func() {
		_, err := stdcopy.StdCopy(ws.UnderlyingConn(), ws.UnderlyingConn(), hijackedResp.Reader)
		if err != nil {
			logrus.Errorf("failed to copy data from exec stdout/stderr: %v", err)
		}
		close(done)
	}()

	err = c.client.ContainerExecStart(ctx, execIDResp.ID, container.ExecStartOptions{
		Tty: true,
	})
	if err != nil {
		logrus.Errorf("container exec start error: %+v", err)
		return err
	}
	logrus.Infof("3")

	ws.SetCloseHandler(func(code int, text string) error {
		logrus.Infof("WebSocket closed with code: %d, text: %s", code, text)
		hijackedResp.Close()
		return nil
	})

	<- done
	return nil
}