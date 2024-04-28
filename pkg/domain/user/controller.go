package user

import (
	"github.com/MR5356/aurora/pkg/config"
	"github.com/MR5356/aurora/pkg/response"
	"github.com/MR5356/aurora/pkg/util/ginutil"
	"github.com/gin-gonic/gin"
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
			userinfo, err := c.service.GetUserInfo(oAuthName, code)
			if err != nil {
				response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
				return
			}
			user := new(User)
			user.ID = userinfo.ID
			user.Username = userinfo.Username
			user.Nickname = userinfo.Nickname
			user.Email = userinfo.Email
			user.Phone = userinfo.Phone
			user.Avatar = userinfo.Avatar

			token, err := GetJWTService().CreateToken(user)
			if err != nil {
				response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
				return
			}
			ctx.SetCookie("token", token, int(config.Current().JWT.Expire.Seconds()), "", "", false, false)
			ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		}
	}
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
		response.ErrorWithMsg(ctx, response.CodeParamsError, err.Error())
	} else {
		response.Success(ctx, user)
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/user")

	// add user
	//api.POST("", c.handleAddUser)

	// get user info
	api.GET("/info", c.handleUserInfo)

	// get oauth url
	api.GET("/oauth", c.handleGetOauthUrl)

	// callback
	api.GET("/callback", c.handleCallback)
}
