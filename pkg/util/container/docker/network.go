package docker

import (
	"context"

	"github.com/MR5356/aurora/pkg/util/container"
	"github.com/docker/docker/api/types/network"
)

func (c *Client) ListNetwork(ctx context.Context) ([]*container.Network, error) {
	var res []*container.Network
	if networks, err := c.client.NetworkList(ctx, network.ListOptions{}); err != nil {
		return nil, err
	} else {
		for _, nw := range networks {
			res = append(res, &container.Network{
				ID: nw.ID,
				Name: nw.Name,
				Driver: nw.Driver,
				IPv6: nw.EnableIPv6,
				Internal: nw.Internal,
				Scope: nw.Scope,
			})
		}
		return res, nil
	}
}