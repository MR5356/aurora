package schedule

import "github.com/sirupsen/logrus"

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
