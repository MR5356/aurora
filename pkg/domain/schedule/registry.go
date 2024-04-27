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

func (m *Manager) Register(name string, task func() Task) error {
	if _, ok := m.tasks.Load(name); ok {
		return fmt.Errorf("task executor %s already registered", name)
	}
	m.tasks.Store(name, task)
	logrus.Infof("task executor %s registered", name)
	return nil
}

func (m *Manager) GetExecutor(name string) (func() Task, error) {
	if task, ok := m.tasks.Load(name); ok {
		return task.(func() Task), nil
	} else {
		return nil, fmt.Errorf("task executor %s not found", name)
	}
}

func (m *Manager) GetExecutors() []Executor {
	res := make([]Executor, 0)

	m.tasks.Range(func(key, value interface{}) bool {
		res = append(res, Executor{
			Name:        key.(string),
			DisplayName: key.(string),
		})
		return true
	})

	return res
}
