package ginutil

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetPageParams(ctx *gin.Context) (page int, size int) {
	var err error
	pageStr := ctx.Query("page")
	sizeStr := ctx.Query("size")

	if page, err = strconv.Atoi(pageStr); err != nil {
		page = 1
	}

	if size, err = strconv.Atoi(sizeStr); err != nil {
		size = 10
	}

	if size > 50 {
		size = 50
	}
	return
}

func GetToken(ctx *gin.Context) string {
	tokenString := ctx.GetHeader("Authorization")
	if len(tokenString) == 0 {
		tokenString, _ = ctx.Cookie("token")
	}
	return tokenString
}
