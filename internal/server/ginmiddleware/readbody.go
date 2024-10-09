package ginmiddleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
)

func ReadBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err == nil {
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		ctx.Next()
	}
}
