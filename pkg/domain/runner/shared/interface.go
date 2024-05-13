package shared

import (
	"github.com/MR5356/aurora/pkg/domain/runner/proto"
	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"os/exec"
)

func Serve(task ITask) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins:         GetPlugins(task),
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}

func GetPlugins(task ...ITask) map[string]plugin.Plugin {
	if len(task) > 0 {
		return map[string]plugin.Plugin{
			"task": &TaskGRPCPlugin{Impl: task[0]},
		}
	} else {
		return map[string]plugin.Plugin{
			"task": &TaskGRPCPlugin{},
		}
	}
}

type ITask interface {
	// GetInfo Get the information of the task
	GetInfo() *proto.TaskInfo
	// GetParams Get the parameters required by the task
	GetParams() *proto.TaskParams
	// SetParams set task params, you can use it to set task configuration
	SetParams(params *proto.TaskParams)

	// Start start the task
	Start() error
	// Stop stop the task
	Stop() error
	// Pause pause the task
	Pause() error
	// Resume resume the task
	Resume() error
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
	return nil
}

func (u *UnimplementedITask) Stop() error {
	return nil
}

func (u *UnimplementedITask) Pause() error {
	return nil
}

func (u *UnimplementedITask) Resume() error {
	return nil
}

type Plugin struct {
	client *plugin.Client
	task   ITask
}

func GetPlugin(path string) (*Plugin, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  Handshake,
		Plugins:          GetPlugins(),
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
	task := anyTask.(ITask)
	//if !ok {
	//	return p, errors.New("task is not ITask")
	//}
	p.task = task
	return p, nil
}

func (p *Plugin) Close() {
	logrus.Info("close plugin")
	p.client.Kill()
}

func (p *Plugin) GetInfo() *proto.TaskInfo {
	return p.task.GetInfo()
}
