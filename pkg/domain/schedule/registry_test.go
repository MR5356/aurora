package schedule

import (
	"github.com/MR5356/aurora/pkg/config"
	"testing"
)

var _ = config.New(config.WithDatabase("sqlite", ":memory:"))

func TestGetExecutorManager(t *testing.T) {
	m := GetExecutorManager()

	err := m.Register("test1", func() Task {
		return &TestTask{}
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = m.Register("test1", func() Task {
		return &TestTask{}
	})

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	_ = m.GetExecutors()

	exec, err := m.GetExecutor("test1")

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
