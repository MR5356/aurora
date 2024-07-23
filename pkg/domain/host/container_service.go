package host

import (
	"context"
	"errors"
	"io"

	"github.com/MR5356/aurora/pkg/util/container"
	"github.com/MR5356/aurora/pkg/util/container/containerd"
	"github.com/MR5356/aurora/pkg/util/container/docker"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (s *Service) ListContainerNetwork(id uuid.UUID) ([]*container.Network, error) {
	if client, err := s.getContainerClient(id); err != nil {
		return nil, err
	} else {
		return client.ListNetwork(context.TODO())
	}
}

func (s *Service) ListContainerImage(id uuid.UUID) ([]*container.Image, error) {
	if client, err := s.getContainerClient(id); err != nil {
		return nil, err
	} else {
		return client.ListImage(context.TODO(), true)
	}
}

func (s *Service) ListContainer(id uuid.UUID) ([]*container.Container, error) {
	if client, err := s.getContainerClient(id); err != nil {
		return nil, err
	} else {
		return client.ListContainer(context.TODO(), true)
	}
}

func (s *Service) GetContainerLogs(ctx context.Context, id uuid.UUID, containerId string) (io.ReadCloser, error) {
	if client, err := s.getContainerClient(id); err != nil {
		return nil, err
	} else {
		return client.Logs(ctx, containerId)
	}
}

func (s *Service) getContainerClient(id uuid.UUID) (container.Client, error) {
	// 优先在缓存中取客户端
	logrus.Debugf("get container client for %s", id.String())
	if client, ok := s.containerClientCache.Get(id.String()); ok {
		return client, nil
	}

	if host, err := s.DetailHost(id); err != nil {
		logrus.Debugf("host id %s not found", id.String())
		return nil, err
	} else {
		var client container.Client
		var err error
		logrus.Debugf("try to use docker driver")
		if client, err = docker.NewClientWithSSH(&host.HostInfo); err != nil {
			logrus.Debugf("docker driver error, try to use containerd")
			if client, err = containerd.NewClientWithSSH(&host.HostInfo); err != nil {
				logrus.Debugf("container driver error")
				return nil, errors.New("failed to obtain the container. Please check whether Docker or containerd is installed")
			}
		}
		s.containerClientCache.Set(id.String(), client)
		return client, err
	}
}
