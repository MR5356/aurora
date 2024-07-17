package host

import (
	"context"

	"github.com/MR5356/aurora/pkg/util/container"
	"github.com/MR5356/aurora/pkg/util/container/docker"
	"github.com/google/uuid"
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

func (s *Service) getContainerClient(id uuid.UUID) (container.Client, error){
	// 优先在缓存中取客户端
	if client, ok := s.containerClientCache.Get(id.String()); ok {
		return client, nil
	}

	if host, err := s.DetailHost(id); err != nil {
		return nil, err
	} else {
		if client, err := docker.NewClientWithSSH(&host.HostInfo); err != nil {
			return nil, err
		} else {
			s.containerClientCache.Set(id.String(), client)
			return client, err
		}
	}
}