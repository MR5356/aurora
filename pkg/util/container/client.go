package container

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
)

type Client interface {
	ContainerList(ctx context.Context, all bool) ([]*Container, error) // Show all containers (default shows just running)
	ImageList(ctx context.Context, all bool) ([]image.Summary, error)  // Show all images (default hides intermediate images)
}

type Container struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Image       string             `json:"image"`
	ImageID     string             `json:"imageId"`
	Command     string             `json:"command"`
	Created     int64              `json:"created"`
	Ports       []types.Port       `json:"ports"`
	Status      string             `json:"status"`
	State       string             `json:"state"`
	NetworkMode string             `json:"networkMode"`
	Mounts      []types.MountPoint `json:"mounts"`
}
