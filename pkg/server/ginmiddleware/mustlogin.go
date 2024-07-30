package ginmiddleware

import (
	"strings"

	"github.com/MR5356/aurora/pkg/domain/user"
	"github.com/MR5356/aurora/pkg/response"
	"github.com/MR5356/aurora/pkg/util/ginutil"
	"github.com/gin-gonic/gin"
)

func skipLogin(path string) bool {
	prefixes := []string{
		"/api/v1/user/info",
		"/api/v1/user/login",
		"/api/v1/user/logout",
		"/api/v1/user/callback",
		"/api/v1/user/oauth",
		"/api/v1/swagger",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func MustLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !skipLogin(ctx.Request.URL.Path) {
			u, err := user.GetJWTService().ParseToken(ginutil.GetToken(ctx))
			if err != nil {
				response.Error(ctx, response.CodeNotLogin)
				ctx.Abort()
				return
			}
			if u.IsBanned() {
				response.Error(ctx, response.CodeBanned)
				ctx.Abort()
				return
			}
			ctx.Set("user", u)
		}
		ctx.Next()
	}
}
