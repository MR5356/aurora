package schedule

import (
	"github.com/MR5356/aurora/pkg/config"
	"testing"
)

var _ = config.New(config.WithDatabase("sqlite", ":memory:"))

func TestGetExecutorManager(t *testing.T) {
	m := GetExecutorManager()

	err := m.Register(Executor{Name: "test", DisplayName: "test executor"}, func() Task {
		return &TestTask{}
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = m.Register(Executor{Name: "test", DisplayName: "test executor"}, func() Task {
		return &TestTask{}
	})

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	executors := m.GetExecutors()

	if len(executors) != 1 {
		t.Errorf("expected 1 tasks, got %d", len(executors))
	}

	exec, err := m.GetExecutor("test")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if exec == nil {
		t.Errorf("expected task, got nil")
	}

	exec, err = m.GetExecutor("test2")

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if exec != nil {
		t.Errorf("expected nil task")
	}
}
