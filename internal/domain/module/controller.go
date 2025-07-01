package module

import (
	"github.com/MR5356/aurora/internal/config"
	"github.com/MR5356/aurora/internal/response"
	"github.com/MR5356/aurora/pkg/util/ginutil"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v61/github"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	return &Controller{
		service: GetService(),
	}
}

func (c *Controller) handleGithubAppInstall(ctx *gin.Context) {
	payload, err := github.ValidatePayload(ctx.Request, []byte(""))
	if err != nil {
		logrus.Errorf("github.ValidatePayload: %v", err)
		response.ErrorWithMsg(ctx, response.CodeParamsError, "invalid payload")
		return
	}

	eventType := github.WebHookType(ctx.Request)
	event, err := github.ParseWebHook(eventType, payload)
	if err != nil {
		logrus.Errorf("github.ParseWebHook: %v", err)
		response.ErrorWithMsg(ctx, response.CodeParamsError, "invalid webhook event")
		return
	}

	var installationID int64 = 0
	var action string = ""
	switch e := event.(type) {
	case *github.InstallationEvent:
		installationID = e.Installation.GetID()
		action = e.GetAction()
	case *github.InstallationRepositoriesEvent:
		installationID = e.Installation.GetID()
		action = e.GetAction()
	}

	logrus.Infof("installation: %+v", installationID)
	logrus.Infof("action: %s", action)
	if installationID == 0 {
		response.Success(ctx, nil)
		return
	}

	if err = c.service.UpdateGithubModule(ctx, action, installationID); err != nil {
		logrus.Errorf("UpdateGithubModule failed: %v", err)
		response.ErrorWithMsg(ctx, response.CodeServerError, "failed to update module")
		return
	} else {
		logrus.Infof("UpdateGithubModule success, event: %s, installationID: %d", eventType, installationID)
		response.Success(ctx, nil)
	}
}

func (c *Controller) handleGithubAppCallback(ctx *gin.Context) {
	installationIDStr := ctx.Query("installation_id")
	userID, ok := ctx.Get(config.ContextUserIDKey)
	if !ok {
		logrus.Error("user ID not found in context")
		response.ErrorWithMsg(ctx, response.CodeNoPermission, "user not authenticated")
		return
	}
	installationID, err := strconv.ParseInt(installationIDStr, 10, 64)
	if err != nil {
		logrus.Errorf("strconv.ParseInt: %v", err)
		response.ErrorWithMsg(ctx, response.CodeParamsError, "invalid installation ID")
		return
	}
	if c.service.RegisterInstallationID(ctx, installationID, userID.(string)) != nil {
		logrus.Errorf("RegisterInstallationID failed: %v", err)
		response.ErrorWithMsg(ctx, response.CodeServerError, "failed to register installation ID")
		return
	} else {
		logrus.Infof("RegisterInstallationID success, installationID: %d, user: %v", installationID, userID)
		response.Success(ctx, nil)
	}
}

func (c *Controller) handleListModules(ctx *gin.Context) {
	page, size := ginutil.GetPageParams(ctx)
	userID, ok := ctx.Get(config.ContextUserIDKey)
	if !ok {
		logrus.Error("user ID not found in context")
		response.ErrorWithMsg(ctx, response.CodeNoPermission, "user not authenticated")
		return
	}
	res, err := c.service.PageModule(ctx, page, size, userID.(string))
	if err != nil {
		logrus.Errorf("PageModule failed: %v", err)
		response.ErrorWithMsg(ctx, response.CodeServerError, "failed to list modules")
		return
	}
	response.Success(ctx, res)
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/module")
	api.POST("/github/app/install", c.handleGithubAppInstall)
	api.GET("/github/app/callback", c.handleGithubAppCallback)
	api.GET("/list", c.handleListModules)
}
