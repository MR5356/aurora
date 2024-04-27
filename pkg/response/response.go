package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func New(ctx *gin.Context, httpCode int, code string, msg string, data any) {
	ctx.JSON(httpCode, &Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

func Success(ctx *gin.Context, data any) {
	New(ctx, http.StatusOK, CodeSuccess, MessageSuccess, data)
}

func Error(ctx *gin.Context, code string, msg string) {
	New(ctx, http.StatusOK, code, msg, nil)
}
