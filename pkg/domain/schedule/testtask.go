package schedule

import (
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/sirupsen/logrus"
)

type TestTask struct {
	params string
}

func (t *TestTask) Run() {
	if t.params == "" {
		panic("test task params is empty")
	}
	logrus.Infof("test task params: %s", t.params)
}

func (t *TestTask) SetParams(params string) {
	t.params = params
}

func init() {
	if err := database.GetDB().AutoMigrate(&Record{}, &Schedule{}); err != nil {
		logrus.Errorf("auto migrate failed, error: %v", err)
	}

	if err := GetExecutorManager().Register("test", func() Task {
		return &TestTask{}
	}); err != nil {
		logrus.Errorf("register test task failed, error: %v", err)
	}
}
