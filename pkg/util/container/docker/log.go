package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
)

func (c *Client) Logs(ctx context.Context, containerId string) (io.ReadCloser, error) {
	return c.client.ContainerLogs(ctx, containerId, container.LogsOptions{
		Follow: true,
		ShowStdout: true,
		ShowStderr: true,
	})
}