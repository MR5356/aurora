package containerd

import (
	"context"
	"github.com/MR5356/aurora/pkg/util/container"
	"github.com/containerd/containerd/namespaces"
)

func (c *Client) ImageList(ctx context.Context, all bool) ([]*container.Image, error) {
	var result []*container.Image
	if ns, err := c.client.NamespaceService().List(ctx); err != nil {
		return nil, err
	} else {
		for _, n := range ns {
			if images, err := c.getImageByNamespace(ctx, n); err != nil {
				return nil, err
			} else {
				result = append(result, images...)
			}
		}
	}

	return result, nil
}

func (c *Client) ContainerList(ctx context.Context, all bool) ([]*container.Container, error) {
	var result []*container.Container
	if ns, err := c.client.NamespaceService().List(ctx); err != nil {
		return nil, err
	} else {
		for _, n := range ns {
			if containers, err := c.getContainerByNamespace(ctx, n); err != nil {
				return nil, err
			} else {
				result = append(result, containers...)
			}
		}
	}
	return result, nil
}

func (c *Client) getImageByNamespace(ctx context.Context, namespace string) ([]*container.Image, error) {
	var result []*container.Image
	ctx = namespaces.WithNamespace(ctx, namespace)
	images, err := c.client.ListImages(ctx)
	if err != nil {
		return nil, err
	}
	for _, i := range images {
		result = append(result, &container.Image{
			ID:      i.Target().Digest.Hex(),
			Labels:  i.Labels(),
			Size:    i.Target().Size,
			Name:    i.Name(),
			Created: i.Metadata().CreatedAt.Unix(),
		})
	}
	return result, nil
}

func (c *Client) getContainerByNamespace(ctx context.Context, namespace string) ([]*container.Container, error) {
	var result []*container.Container
	ctx = namespaces.WithNamespace(ctx, namespace)
	cs, err := c.client.Containers(ctx)
	if err != nil {
		return nil, err
	}
	for _, c := range cs {
		result = append(result, &container.Container{
			ID:   c.ID(),
			Name: c.ID(),
		})
	}
	return result, nil
}
