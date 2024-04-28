package schedule

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	manager *Manager
	once    sync.Once
)

type Manager struct {
	tasks sync.Map
}

func GetExecutorManager() *Manager {
	once.Do(func() {
		manager = &Manager{
			tasks: sync.Map{},
		}
	})

	return manager
}

func (m *Manager) Register(executor Executor, task func() Task) error {
	if _, ok := m.tasks.Load(executor.Name); ok {
		return fmt.Errorf("task executor %s already registered", executor.Name)
	}
	executor.task = task
	m.tasks.Store(executor.Name, executor)
	logrus.Infof("task executor %s registered", executor.Name)
	return nil
}

func (m *Manager) GetExecutor(name string) (func() Task, error) {
	if task, ok := m.tasks.Load(name); ok {
		return task.(Executor).task, nil
	} else {
		return nil, fmt.Errorf("task executor %s not found", name)
	}
}

func (m *Manager) GetExecutors() []Executor {
	res := make([]Executor, 0)

	m.tasks.Range(func(key, value interface{}) bool {
		res = append(res, value.(Executor))
		return true
	})

	return res
}
