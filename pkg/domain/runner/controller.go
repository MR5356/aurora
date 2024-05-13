package runner

import (
	"github.com/MR5356/aurora/pkg/config"
	"github.com/MR5356/aurora/pkg/domain/runner/shared"
	"github.com/MR5356/aurora/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

// @Summary	register plugin
// @Tags		runner
// @Accept		multipart/form-data
// @Produce	application/json
// @Param		file formData file true "file"
// @Router		/runner/plugin/register [post]
// @Success	200 {object} response.Response
func (c *Controller) handleRegisterPlugin(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		logrus.Errorf("register plugin failed, error: %v", err)
		response.Error(ctx, response.CodeParamsError)
		return
	}
	dst := config.Current().Server.PluginPath + "/" + file.Filename

	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}

	err = os.Chmod(dst, 0755)
	if err != nil {
		logrus.Errorf("chmod plugin failed, error: %v", err)
		response.Error(ctx, response.CodeParamsError)
		return
	}

	p, err := shared.GetPlugin(dst)
	defer p.Close()
	if err != nil {
		logrus.Errorf("get plugin info failed, error: %v", err)
		response.Error(ctx, response.CodeParamsError)
		return
	}
	response.Success(ctx, p.GetTask().GetInfo())
}

func (c *Controller) handleRunner(ctx *gin.Context) {
	token, ok := ctx.GetQuery("token")
	if !ok {
		response.Error(ctx, response.CodeParamsError)
	}
	logrus.Infof("get token: %s", token)

	upgrager := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{"runner"},
	}

	webConn, err := upgrager.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.Errorf("upgrade failed, error: %v", err)
		response.Error(ctx, response.CodeParamsError)
	}
	defer webConn.Close()

	closeCh := make(chan struct{})
	webConn.SetCloseHandler(func(code int, text string) error {
		logrus.Infof("close connection")
		closeCh <- struct{}{}
		return nil
	})

	go func() {
		for {
			_, data, err := webConn.ReadMessage()
			if err != nil {
				if !strings.Contains(err.Error(), "close 1000 (normal)") {
					logrus.Errorf("read message failed, error: %+v", err)
				}
				return
			}
			logrus.Infof("receive message: %s", string(data))
		}
	}()

	select {
	case <-closeCh:
		logrus.Infof("close connection")
		return
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/runner")

	api.GET("/conn", c.handleRunner)

	pluginApi := api.Group("/plugin")
	pluginApi.POST("/register", c.handleRegisterPlugin)
}
