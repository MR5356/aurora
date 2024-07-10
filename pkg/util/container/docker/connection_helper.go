package docker

import (
	"context"
	"net"
)

type ConnectionHelper interface {
	GetDialer() func(ctx context.Context, network, addr string) (net.Conn, error)
	GetHost() string
}
