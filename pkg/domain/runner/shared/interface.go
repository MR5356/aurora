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

	// GetWorkflow get the workflow
	GetWorkflow() *Workflow
	// DryRun dry run the task
	DryRun() error

	// Start start the task
	Start() error
	// Stop stop the task
	Stop() error
	// Pause pause the task
	Pause() error
	// Resume resume the task
	Resume() error
}

type Edge struct {
	Source ITask
	Target ITask
}

type Workflow struct {
	tasks []ITask
	edges []Edge
}

func (w *Workflow) AddTask(task ITask) {
	w.tasks = append(w.tasks, task)
}

func (w *Workflow) AddEdge(from, to ITask) {
	w.edges = append(w.edges, Edge{Source: from, Target: to})
}

func (w *Workflow) HasCycle() bool {
	visited := make(map[ITask]bool)
	recursionStack := make(map[ITask]bool)

	for _, task := range w.tasks {
		var chain []ITask
		if w.dfs(task, visited, recursionStack, &chain) {
			return true
		}
	}
	return false
}

func (w *Workflow) dfs(task ITask, visited map[ITask]bool, recursionStack map[ITask]bool, chain *[]ITask) bool {
	visited[task] = true
	recursionStack[task] = true
	defer delete(recursionStack, task)

	for _, edge := range w.edges {
		if edge.Source == task {
			*chain = append(*chain, edge.Target)
			if recursionStack[edge.Target] {
				// 发现环路
				*chain = append(*chain, edge.Target)
				return true
			}
			if !visited[edge.Target] {
				if w.dfs(edge.Target, visited, recursionStack, chain) {
					return true
				}
			}
		}
	}
	return false
}
