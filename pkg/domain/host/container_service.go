package host

import (
	"context"
	"fmt"
	"io"

	"github.com/MR5356/aurora/pkg/util/container"
	"github.com/MR5356/aurora/pkg/util/container/containerd"
	"github.com/MR5356/aurora/pkg/util/container/docker"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	driverDontainerd = "containerd"
	driverDocker     = "docker"
)

func (s *Service) ListContainerNetwork(id uuid.UUID, driver string) ([]*container.Network, error) {
	if client, err := s.getContainerClient(id, driver); err != nil {
		return nil, err
	} else {
		return client.ListNetwork(context.TODO())
	}
}

func (s *Service) ListContainerImage(id uuid.UUID, driver string) ([]*container.Image, error) {
	if client, err := s.getContainerClient(id, driver); err != nil {
		return nil, err
	} else {
		return client.ListImage(context.TODO(), true)
	}
}

func (s *Service) ListContainer(id uuid.UUID, driver string) ([]*container.Container, error) {
	if client, err := s.getContainerClient(id, driver); err != nil {
		return nil, err
	} else {
		return client.ListContainer(context.TODO(), true)
	}
}

func (s *Service) GetContainerLogs(ctx context.Context, id uuid.UUID, containerId string, driver string) (io.ReadCloser, error) {
	if client, err := s.getContainerClient(id, driver); err != nil {
		return nil, err
	} else {
		return client.Logs(ctx, containerId)
	}
}

func (s *Service) ExecContainer(ctx *gin.Context, id uuid.UUID, containerId, driver, user, cmd string) error {
	logrus.Infof("come on")
	if client, err := s.getContainerClient(id, driver); err != nil {
		return err
	} else {
		return client.Terminal(ctx, containerId, user, cmd)
	}
}

func (s *Service) getContainerClient(id uuid.UUID, driver string) (container.Client, error) {
	// 优先在缓存中取客户端
	key := fmt.Sprintf("%s-%s", id.String(), driver)
	logrus.Debugf("get container client for %s", key)
	if client, ok := s.containerClientCache.Get(key); ok {
		return client, nil
	}

	if host, err := s.DetailHost(id); err != nil {
		logrus.Debugf("host id %s not found", id.String())
		return nil, err
	} else {
		var client container.Client
		var err error

		switch driver {
		case driverDontainerd:
			client, err = containerd.NewClientWithSSH(&host.HostInfo)
		case driverDocker:
			client, err = docker.NewClientWithSSHAndAPIVersion(&host.HostInfo, host.MetaInfo.Docker)
		default:
			return nil, fmt.Errorf("%s not support", driver)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to obtain the container. error: %w", err)
		}

		s.containerClientCache.Set(key, client)
		return client, err
	}
}
