package schedule

import (
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/MR5356/aurora/pkg/middleware/eventbus"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/sirupsen/logrus"
)

const (
	TaskStatusRunning = "running"
	TaskStatusError   = "error"
	TaskStatusSuccess = "success"
)

type Task interface {
	// Run task entry
	Run()
	// SetParams set task params, you can use it to set task configuration
	SetParams(params string)
}

type WrappedTask struct {
	Task

	schedule *Schedule
	db       *database.BaseMapper[*Record]
}

func NewWrapper(task Task, schedule *Schedule) *WrappedTask {
	return &WrappedTask{
		Task:     task,
		schedule: schedule,
		db:       database.NewMapper(database.GetDB(), &Record{}),
	}
}

func (w *WrappedTask) Run() {
	// attempt to lock
	if err := eventbus.GetEventBus().TryLock(w.schedule.ID.String()); err != nil {
		logrus.Debugf("task %s try lock failed, skip, error: %v", w.schedule.ID, err)
		return
	}

	defer func(bus eventbus.EventBus, key string) {
		err := bus.UnLock(key)
		if err != nil {
			logrus.Errorf("task %s unlock failed, error: %v", w.schedule.ID, err)
		}
	}(eventbus.GetEventBus(), w.schedule.ID.String())

	// insert record
	record := &Record{
		ScheduleID: w.schedule.ID,
		Title:      w.schedule.Title,
		TaskName:   w.schedule.Executor,
		Params:     w.schedule.Params,
		Status:     TaskStatusRunning,
	}

	err := w.db.Insert(record)
	if err != nil {
		logrus.Errorf("task %s record insert record failed, error: %v", w.schedule.ID, err)
	}

	defer func() {
		if finalErr := recover(); finalErr != nil {
			logrus.Errorf("task %s recover failed, error: %v", w.schedule.ID, finalErr)
			record.Status = TaskStatusError
		} else {
			record.Status = TaskStatusSuccess
		}

		if updateErr := w.db.Update(&Record{ID: record.ID}, structutil.Struct2Map(record)); updateErr != nil {
			logrus.Errorf("task %s record update failed, error: %v", w.schedule.ID, updateErr)
		}
	}()

	// run task
	w.Task.Run()
}
