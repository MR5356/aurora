package terminal

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerTerminal struct {
	config     Config
	dc         DockerConfig
	client     *client.Client
	hijacked   types.HijackedResponse
	responseId string
	ctx        context.Context
}

type DockerConfig struct {
	ContainerId string
	Endpoint    []string
}

func NewDockerTerminal(dc DockerConfig, config Config) *DockerTerminal {
	if dc.Endpoint == nil || len(dc.Endpoint) == 0 {
		dc.Endpoint = []string{"/bin/sh"}
	}
	return &DockerTerminal{
		config: config,
		dc:     dc,
		ctx:    context.Background(),
	}
}

func (t *DockerTerminal) Start() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	t.client = cli

	execConfig := container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          t.dc.Endpoint,
	}

	response, err := t.client.ContainerExecCreate(t.ctx, t.dc.ContainerId, execConfig)
	if err != nil {
		return err
	}
	t.responseId = response.ID

	hijack, err := t.client.ContainerExecAttach(t.ctx, response.ID, container.ExecAttachOptions{
		Tty: true,
	})
	if err != nil {
		return err
	}

	t.hijacked = hijack

	return t.Resize(t.config.Cols, t.config.Rows)
}

func (t *DockerTerminal) Close() error {
	if t.hijacked.Conn != nil {
		t.hijacked.Close()
	}
	return nil
}

func (t *DockerTerminal) Write(p []byte) (n int, err error) {
	return t.hijacked.Conn.Write(p)
}

func (t *DockerTerminal) Read(p []byte) (n int, err error) {
	return t.hijacked.Conn.Read(p)
}

func (t *DockerTerminal) Resize(cols, rows uint32) error {
	return t.client.ContainerExecResize(t.ctx, t.responseId, container.ResizeOptions{
		Height: uint(rows),
		Width:  uint(cols),
	})
}
