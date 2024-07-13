package user

import (
	"github.com/MR5356/aurora/pkg/response"
	"github.com/MR5356/aurora/pkg/util/ginutil"
	"github.com/gin-gonic/gin"
)

func MustAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		u, err := GetJWTService().ParseToken(ginutil.GetToken(ctx))
		if err != nil {
			response.Error(ctx, response.CodeNotLogin)
			ctx.Abort()
			return
		}

		if !u.IsAdmin() {
			response.Error(ctx, response.CodeNoPermission)
			ctx.Abort()
			return
		}
		ctx.Next()
	}

}
