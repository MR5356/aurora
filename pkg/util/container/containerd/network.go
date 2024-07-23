package containerd

import (
	"context"

	"github.com/MR5356/aurora/pkg/util/container"
)

func (c *Client) ListNetwork(ctx context.Context) ([]*container.Network, error) {
	// TODO: Get containerd cni
	return nil, nil
}
