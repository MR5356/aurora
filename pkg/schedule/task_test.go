package schedule

import (
	"github.com/google/uuid"
	"testing"
)

func TestNewWrapper(t *testing.T) {
	task := &TestTask{}

	task.SetParams("test")

	wrappedTask := NewWrapper(task, &Schedule{
		ID:         uuid.UUID{},
		Title:      "test",
		Desc:       "test",
		CronString: "*/5 * * * * *",
		TaskName:   "test",
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
		TaskName:   "test",
		Params:     "",
		Enabled:    true,
	})

	wrappedTask.Run()
}
