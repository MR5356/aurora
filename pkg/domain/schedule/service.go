package schedule

import (
	"encoding/json"
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/MR5356/aurora/pkg/middleware/eventbus"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/MR5356/aurora/pkg/util/validate"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	topicAddCronTask = "topic.schedule.add_cron_task"
	topicDelCronTask = "topic.schedule.del_cron_task"
)

var (
	onceService sync.Once
	service     *Service
)

type Service struct {
	scheduleDB *database.BaseMapper[*Schedule]
	recordDB   *database.BaseMapper[*Record]
	cron       *cron.Cron
	cronJobMap sync.Map
}

func GetService() *Service {
	onceService.Do(func() {
		c := cron.New(cron.WithSeconds())
		c.Start()

		service = &Service{
			scheduleDB: database.NewMapper(database.GetDB(), &Schedule{}),
			recordDB:   database.NewMapper(database.GetDB(), &Record{}),
			cron:       c,
			cronJobMap: sync.Map{},
		}
	})
	return service
}

// AddSchedule add schedule
func (s *Service) AddSchedule(schedule *Schedule) error {
	if err := s.verifyTaskParams(schedule); err != nil {
		logrus.Errorf("verify task params failed, error: %v", err)
		return err
	}

	// use transaction
	tx := s.scheduleDB.DB.Begin()
	defer tx.Rollback()
	if err := s.scheduleDB.Insert(schedule, tx); err != nil {
		logrus.Errorf("insert schedule failed, error: %v", err)
		return err
	}

	if schedule.Enabled {
		ps, err := json.Marshal(schedule)
		if err != nil {
			logrus.Errorf("marshal schedule params failed, error: %v", err)
			return err
		}
		if err := eventbus.GetEventBus().Publish(topicAddCronTask, string(ps)); err != nil {
			logrus.Errorf("publish add cron task failed, error: %v", err)
			return err
		}
	}

	tx.Commit()
	return nil
}

func (s *Service) addCronTask(params string) {
	logrus.Infof("add cron task: %s", params)

	schedule := new(Schedule)
	err := json.Unmarshal([]byte(params), schedule)
	if err != nil {
		logrus.Errorf("unmarshal schedule params failed, error: %v", err)
		return
	}

	if _, ok := s.cronJobMap.Load(schedule.ID); ok {
		logrus.Errorf("task %s already exists", schedule.ID)
		return
	}

	taskFunc, err := GetExecutorManager().GetExecutor(schedule.Executor)
	if err != nil {
		logrus.Errorf("get task executor failed, error: %v", err)
		return
	}

	f := taskFunc()
	f.SetParams(schedule.Params)

	jobId, err := s.cron.AddJob(schedule.CronString, NewWrapper(f, schedule))
	if err != nil {
		logrus.Errorf("add cron job failed, error: %v", err)
		return
	}
	s.cronJobMap.Store(schedule.ID, jobId)
}

func (s *Service) delCronTask(id uuid.UUID) {
	logrus.Infof("del cron task: %s", id)

	jobId, ok := s.cronJobMap.Load(id)
	if !ok {
		logrus.Errorf("task %s not found", id)
		return
	}

	s.cron.Remove(jobId.(cron.EntryID))
	s.cronJobMap.Delete(id)
}

// DeleteSchedule delete schedule
func (s *Service) DeleteSchedule(id uuid.UUID) error {
	tx := s.scheduleDB.DB.Begin()
	defer tx.Rollback()
	if err := s.scheduleDB.Delete(&Schedule{ID: id}, tx); err != nil {
		logrus.Errorf("delete schedule failed, error: %v", err)
		return err
	}

	if err := eventbus.GetEventBus().Publish(topicDelCronTask, id); err != nil {
		logrus.Errorf("publish del cron task failed, error: %v", err)
		return err
	}

	tx.Commit()
	return nil
}

// UpdateSchedule update schedule
func (s *Service) UpdateSchedule(schedule *Schedule) error {
	if err := s.verifyTaskParams(schedule); err != nil {
		logrus.Errorf("verify task params failed, error: %v", err)
		return err
	}

	tx := s.scheduleDB.DB.Begin()
	defer tx.Rollback()
	if err := s.scheduleDB.Update(&Schedule{ID: schedule.ID}, structutil.Struct2Map(schedule), tx); err != nil {
		logrus.Errorf("update schedule failed, error: %v", err)
		return err
	}

	// delete cron task
	if err := eventbus.GetEventBus().Publish(topicDelCronTask, schedule.ID); err != nil {
		logrus.Errorf("publish del cron task failed, error: %v", err)
		return err
	}

	// add cron task
	if schedule.Enabled {
		ps, err := json.Marshal(schedule)
		if err != nil {
			logrus.Errorf("marshal schedule params failed, error: %v", err)
			return err
		}
		if err := eventbus.GetEventBus().Publish(topicAddCronTask, string(ps)); err != nil {
			logrus.Errorf("publish add cron task failed, error: %v", err)
			return err
		}
	}

	tx.Commit()
	return nil
}

// DetailSchedule detail schedule
func (s *Service) DetailSchedule(id uuid.UUID) (*Schedule, error) {
	return s.scheduleDB.Detail(&Schedule{ID: id})
}

func (s *Service) verifyTaskParams(schedule *Schedule) error {
	if err := validate.Validate(schedule); err != nil {
		return err
	}

	// verify task executor
	if _, err := GetExecutorManager().GetExecutor(schedule.Executor); err != nil {
		return err
	}

	// verify cron string
	parser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	if _, err := parser.Parse("*/5 * * * * *"); err != nil {
		return err
	}
	return nil
}

// GetTaskExecutors get task executors
func (s *Service) GetTaskExecutors() []Executor {
	return GetExecutorManager().GetExecutors()
}

// PageScheduleRecord page task records
func (s *Service) PageScheduleRecord(num, size int, record *Record) (*database.Pager[*Record], error) {
	return s.recordDB.Page(record, int64(num), int64(size))
}

// PageSchedule page schedules
func (s *Service) PageSchedule(num, size int, schedule *Schedule) (*database.Pager[*Schedule], error) {
	return s.scheduleDB.Page(schedule, int64(num), int64(size))
}

func (s *Service) Initialize() (err error) {
	if err = database.GetDB().AutoMigrate(&Record{}, &Schedule{}); err != nil {
		return err
	}

	if err := GetExecutorManager().Register("test", func() Task {
		return &TestTask{}
	}); err != nil {
		logrus.Errorf("register test task failed, error: %v", err)
	}

	if err = eventbus.GetEventBus().Subscribe(topicAddCronTask, s.addCronTask); err != nil {
		return err
	}

	if err = eventbus.GetEventBus().Subscribe(topicDelCronTask, s.delCronTask); err != nil {
		return err
	}

	if jobs, err := s.scheduleDB.List(&Schedule{Enabled: true}); err != nil {
		return err
	} else {
		for _, job := range jobs {
			if js, err := json.Marshal(job); err != nil {
				return err
			} else {
				if err = eventbus.GetEventBus().Publish(topicAddCronTask, string(js)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
