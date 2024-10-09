package user

import (
	"github.com/MR5356/aurora/internal/config"
	"github.com/MR5356/aurora/internal/response"
	"github.com/MR5356/aurora/pkg/util/ginutil"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
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

// @Summary	get oauth url
// @Tags		user
// @Param		oauth		query		string	true	"auth type"
// @Param		redirectURL	query		string	false	"redirect url"
// @Success	200			{object}	response.Response{data=string}
// @Router		/user/oauth [get]
// @Produce	json
func (c *Controller) handleGetOauthUrl(ctx *gin.Context) {
	if oAuthName, ok := ctx.GetQuery("oauth"); !ok {
		response.Error(ctx, response.CodeParamsError)
		return
	} else {
		redirectURL, ok := ctx.GetQuery("redirectURL")
		if !ok {
			redirectURL = "/"
		}

		if res, err := c.service.GetOAuthURL(oAuthName, redirectURL); err != nil {
			response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
		} else {
			response.Success(ctx, res)
		}
	}
}

// @Summary	list user
// @Tags		user
// @Success	200	{object}	response.Response{data=[]User}
// @Router		/user/list [get]
// @Produce	json
func (c *Controller) handleListUser(ctx *gin.Context) {
	if users, err := c.service.ListUser(new(User)); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		response.Success(ctx, users)
	}
}

func (c *Controller) handleCallback(ctx *gin.Context) {
	if oAuthName, ok := ctx.GetQuery("oauth"); !ok {
		response.Error(ctx, response.CodeParamsError)
		return
	} else {
		if code, ok := ctx.GetQuery("code"); !ok {
			response.Error(ctx, response.CodeParamsError)
			return
		} else {
			redirectURL, ok := ctx.GetQuery("state")
			if !ok {
				redirectURL = "/"
			}

			token, err := c.service.GetOAuthToken(oAuthName, code)
			if err != nil {
				response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
				return
			}

			ctx.SetCookie("token", token, int(config.Current().JWT.Expire.Seconds()), "", "", false, false)
			ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		}
	}
}

// @Summary	logout
// @Tags		user
// @Success	200	{object}	response.Response
// @Router		/user/logout [get]
// @Produce	json
func (c *Controller) handleLogout(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "", "", false, false)
	response.Success(ctx, nil)
}

// @Summary	get user info
// @Tags		user
// @Success	200	{object}	response.Response{data=User}
// @Router		/user/info [get]
// @Produce	json
func (c *Controller) handleUserInfo(ctx *gin.Context) {
	token := ginutil.GetToken(ctx)

	user, err := GetJWTService().ParseToken(token)
	if err != nil {
		logrus.Errorf("parse token failed, error: %v", err)
		response.Success(ctx, nil)
	} else {
		response.Success(ctx, user)
	}
}

// @Summary	get available oauth
// @Tags		user
// @Success	200	{object}	response.Response{data=[]oauth.AvailableOAuth}
// @Router		/user/oauth/all [get]
// @Produce	json
func (c *Controller) handleGetAvailableOauth(ctx *gin.Context) {
	response.Success(ctx, c.service.GetAvailableOAuth())
}

// @Summary	ban user
// @Tags		user
// @Param		id	body		Uid	true	"uid"
// @Success	200	{object}	response.Response
// @Router		/user/ban [post]
// @Produce	json
func (c *Controller) handleBanUser(ctx *gin.Context) {
	var uid Uid
	if err := ctx.ShouldBindJSON(&uid); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}

	if currentUser, ok := ctx.Get("user"); !ok {
		response.Error(ctx, response.CodeNotLogin)
		return
	} else {
		if currentUser.(*User).ID == uid.ID {
			response.ErrorWithMsg(ctx, response.CodeParamsError, "cannot ban yourself")
			return
		}
	}

	if err := c.service.SetUserStatus(&User{ID: uid.ID}, StatusBan); err != nil {
		logrus.Infof("set user status failed, error: %v", err)
		response.Error(ctx, response.CodeParamsError)
	} else {
		response.Success(ctx, nil)
	}
}

// @Summary	unban user
// @Tags		user
// @Param		id	body		Uid	true	"uid"
// @Success	200	{object}	response.Response
// @Router		/user/unban [post]
// @Produce	json
func (c *Controller) handleUnbanUser(ctx *gin.Context) {
	var uid Uid
	if err := ctx.ShouldBindJSON(&uid); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}

	if err := c.service.SetUserStatus(&User{ID: uid.ID}, StatusActive); err != nil {
		response.Error(ctx, response.CodeParamsError)
	} else {
		response.Success(ctx, nil)
	}
}

// @Summary	login
// @Tags		user
// @Success	200	{object}	response.Response
// @Router		/user/login [post]
// @Param		user	body	LoginRequest	true	"user"
// @Produce	json
func (c *Controller) handleLogin(ctx *gin.Context) {
	user := &LoginRequest{}

	if err := ctx.ShouldBindJSON(user); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}

	token, err := c.service.Login(user)
	if err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
	} else {
		ctx.SetCookie("token", token, int(config.Current().JWT.Expire.Seconds()), "", "", false, false)
		//ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		response.Success(ctx, nil)
	}
}

// @Summary	reset password
// @Tags		user
// @Success	200	{object}	response.Response
// @Router		/user/reset [post]
// @Param		user	body	ResetPasswordRequest	true	"user"
// @Produce	json
func (c *Controller) handleResetPassword(ctx *gin.Context) {
	user := &ResetPasswordRequest{}

	if err := ctx.ShouldBindJSON(user); err != nil {
		response.Error(ctx, response.CodeParamsError)
		return
	}

	isAdmin := false
	if currentUser, ok := ctx.Get("user"); ok {
		isAdmin = currentUser.(*User).IsAdmin()
	}
	if err := c.service.ResetPassword(user, isAdmin); err != nil {
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/user")

	admin := api.Group("")
	admin.Use(MustAdmin())

	// add user
	//api.POST("", c.handleAddUser)

	// list user
	admin.GET("list", c.handleListUser)

	// ban
	admin.POST("/ban", c.handleBanUser)

	// unban
	admin.POST("/unban", c.handleUnbanUser)

	// reset password
	admin.POST("/reset", c.handleResetPassword)

	// get user info
	api.GET("/info", c.handleUserInfo)

	// login
	api.POST("/login", c.handleLogin)

	// get available oauth
	api.GET("/oauth/all", c.handleGetAvailableOauth)

	// get oauth url
	api.GET("/oauth", c.handleGetOauthUrl)

	// callback
	api.GET("/callback", c.handleCallback)

	// logout
	api.GET("/logout", c.handleLogout)
}
