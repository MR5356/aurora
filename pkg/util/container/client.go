package container

import (
	"context"
	"github.com/docker/docker/api/types"
)

type Client interface {
	ContainerList(ctx context.Context, all bool) ([]*Container, error) // Show all containers (default shows just running)
	ImageList(ctx context.Context, all bool) ([]*Image, error)         // Show all images (default hides intermediate images)
	Close()
}

type Container struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Image       string       `json:"image"`
	ImageID     string       `json:"imageId"`
	Command     string       `json:"command"`
	Created     int64        `json:"created"`
	Ports       []types.Port `json:"ports"`
	Status      string       `json:"status"`
	State       string       `json:"state"`
	NetworkMode string       `json:"networkMode"`
	Mounts      []Mount      `json:"mounts"`
	Runtime     string       `json:"runtime"`
}

type Mount struct {
	Source string `json:"source"`
	Dest   string `json:"dest"`
	Type   string `json:"type"`
}

type Image struct {
	ID      string            `json:"id"`
	Labels  map[string]string `json:"labels"`
	Size    int64             `json:"size"`
	Name    string            `json:"name"`
	Created int64             `json:"created"`
}
