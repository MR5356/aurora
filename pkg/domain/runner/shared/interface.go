package shared

import (
	"github.com/MR5356/aurora/pkg/domain/runner/proto"
	"github.com/hashicorp/go-plugin"
)

func Serve(task ITask) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins:         GetPluginMap(task),
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}

func GetPluginMap(task ...ITask) map[string]plugin.Plugin {
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
