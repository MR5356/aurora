package shared

import (
	"errors"
	"github.com/MR5356/aurora/pkg/domain/runner/proto"
	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"os/exec"
)

type Plugin struct {
	client *plugin.Client
	task   ITask
}

func GetPlugin(path string) (*Plugin, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  Handshake,
		Plugins:          GetPluginMap(),
		Cmd:              exec.Command("sh", "-c", path),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	p := &Plugin{
		client: client,
	}

	rpcClient, err := client.Client()
	if err != nil {
		logrus.Errorf("failed to get client: %v", err)
		return p, err
	}

	anyTask, err := rpcClient.Dispense("task")
	if err != nil {
		logrus.Errorf("failed to dispense task: %v", err)
		return p, err
	}
	task, ok := anyTask.(ITask)
	if !ok {
		return p, errors.New("task is not ITask")
	}
	p.task = task
	return p, nil
}

func (p *Plugin) Close() {
	logrus.Debugf("close plugin")
	p.client.Kill()
}

func (p *Plugin) GetTask() ITask {
	return p.task
}

type UnimplementedITask struct{}

func (u *UnimplementedITask) GetInfo() *proto.TaskInfo {
	return &proto.TaskInfo{}
}

func (u *UnimplementedITask) GetParams() *proto.TaskParams {
	return &proto.TaskParams{}
}

func (u *UnimplementedITask) SetParams(params *proto.TaskParams) {}

func (u *UnimplementedITask) Start() error {
	return errors.New("method Start not implemented")
}

func (u *UnimplementedITask) Stop() error {
	return errors.New("method Stop not implemented")
}

func (u *UnimplementedITask) Pause() error {
	return errors.New("method Pause not implemented")
}

func (u *UnimplementedITask) Resume() error {
	return errors.New("method Resume not implemented")
}

func (u *UnimplementedITask) GetWorkflow() *Workflow {
	return nil
}

func (u *UnimplementedITask) DryRun() error {
	return nil
}
