package schedule

import (
	"github.com/MR5356/aurora/pkg/config"
	"github.com/google/uuid"
	"testing"
)

var _ = config.New(config.WithDatabase("sqlite", ":memory:"))

func TestNewWrapper(t *testing.T) {
	task := &TestTask{}

	task.SetParams("test")

	wrappedTask := NewWrapper(task, &Schedule{
		ID:         uuid.UUID{},
		Title:      "test",
		Desc:       "test",
		CronString: "*/5 * * * * *",
		Executor:   "test",
		Params:     "",
		Enabled:    true,
	})

	wrappedTask.Run()
}

func TestNewWrapperWithError(t *testing.T) {
	task := &TestTask{}

	wrappedTask := NewWrapper(task, &Schedule{
		ID:         uuid.UUID{},
		Title:      "test",
		Desc:       "test",
		CronString: "*/5 * * * * *",
		Executor:   "test",
		Params:     "",
		Enabled:    true,
	})

	wrappedTask.Run()
}
