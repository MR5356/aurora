package container

import (
	"context"

	"github.com/docker/docker/api/types"
)

type Client interface {
	ListContainer(ctx context.Context, all bool) ([]*Container, error) // Show all containers (default shows just running)
	ListImage(ctx context.Context, all bool) ([]*Image, error)         // Show all images (default hides intermediate images)
	ListNetwork(ctx context.Context) ([]*Network, error)
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

type Network struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Driver   string `json:"driver"`
	IPv6     bool   `json:"ipv6"`
	Internal bool   `json:"internal"`
	Scope    string `json:"scope"`
}
