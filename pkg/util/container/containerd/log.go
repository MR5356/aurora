package containerd

import (
	"context"
	"io"
)

func (c *Client) Logs(ctx context.Context, containerId string) (io.ReadCloser, error) {
	return nil, nil
}