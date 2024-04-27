package server

import "github.com/gin-gonic/gin"

type Controller interface {
	RegisterRoute(group *gin.RouterGroup)
}
