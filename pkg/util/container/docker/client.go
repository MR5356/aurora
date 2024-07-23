package docker

import (
	"context"
	auContainer "github.com/MR5356/aurora/pkg/util/container"
	"github.com/MR5356/aurora/pkg/util/sshutil"
	"github.com/docker/cli/cli/command/formatter"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"net/http"
	"strings"
)

const (
	defaultDockerVersion = "1.40"
)

type Client struct {
	client *client.Client
}

func NewClientWithSSH(sshInfo *sshutil.HostInfo) (*Client, error) {
	helper, err := GetSSHConnectionHelper(sshInfo)
	if err != nil {
		return nil, err
	}

	cli, err := client.NewClientWithOpts(
		client.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				DialContext: helper.Dialer,
			},
		}),
		client.WithHost(helper.Host),
		client.WithDialContext(helper.Dialer),
		client.WithVersion(defaultDockerVersion),
	)
	if err != nil {
		return nil, err
	}

	client := &Client{
		client: cli,
	}

	if _, err = client.Version(context.TODO()); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Close() {
	_ = c.client.Close()
}

func (c *Client) ListContainer(ctx context.Context, all bool) ([]*auContainer.Container, error) {
	var res []*auContainer.Container
	if containers, err := c.client.ContainerList(ctx, container.ListOptions{
		All:     all,
		Filters: filters.NewArgs(),
	}); err != nil {
		return nil, err
	} else {
		for _, ct := range containers {
			mount := make([]auContainer.Mount, 0, len(ct.Mounts))
			for _, m := range ct.Mounts {
				mount = append(mount, auContainer.Mount{
					Source: m.Source,
					Dest:   m.Destination,
					Type:   string(m.Type),
				})
			}
			res = append(res, &auContainer.Container{
				ID:          ct.ID,
				Name:        strings.Join(formatter.StripNamePrefix(ct.Names), ","),
				Image:       ct.Image,
				ImageID:     ct.ImageID,
				Command:     ct.Command,
				Created:     ct.Created,
				Ports:       ct.Ports,
				Status:      ct.Status,
				State:       ct.State,
				NetworkMode: ct.HostConfig.NetworkMode,
				Mounts:      mount,
			})
		}
	}
	return res, nil
}

func (c *Client) ListImage(ctx context.Context, all bool) ([]*auContainer.Image, error) {
	var result []*auContainer.Image
	images, err := c.client.ImageList(ctx, image.ListOptions{
		All:     all,
		Filters: filters.Args{},
	})
	if err != nil {
		return nil, err
	}
	for _, i := range images {
		result = append(result, &auContainer.Image{
			ID:      i.ID,
			Labels:  i.Labels,
			Size:    i.Size,
			Name:    strings.Join(i.RepoTags, ","),
			Created: i.Created,
		})
	}
	return result, nil
}

func (c *Client) Version(ctx context.Context) (types.Version, error) {
	return c.client.ServerVersion(ctx)
}
