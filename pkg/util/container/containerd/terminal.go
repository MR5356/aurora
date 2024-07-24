package containerd

import (
	"github.com/gin-gonic/gin"
)

func (c *Client) Terminal(ctx *gin.Context, containerId string, user, cmd string) error {
	return nil
}