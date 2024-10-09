package script

import (
	"encoding/json"
	"fmt"
	"github.com/MR5356/aurora/internal/domain/host"
	"github.com/MR5356/aurora/internal/domain/schedule"
	"github.com/MR5356/aurora/internal/infrastructure/cache"
	database2 "github.com/MR5356/aurora/internal/infrastructure/database"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/MR5356/aurora/pkg/util/validate"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	onceService sync.Once
	service     *Service
)

type Service struct {
	scriptDB database2.Mapper[*Script]
	hostDB   *database2.BaseMapper[*host.Host]
	recordDB *database2.BaseMapper[*Record]
}

func GetService() *Service {
	onceService.Do(func() {
		service = &Service{
			scriptDB: database2.NewCachedMapper(database2.GetDB(), &Script{}, cache.GetCache()),
			hostDB:   database2.NewMapper(database2.GetDB(), &host.Host{}),
			recordDB: database2.NewMapper(database2.GetDB(), &Record{}),
		}
	})
	return service
}

func (s *Service) AddScript(script *Script) error {
	if err := validate.Validate(script); err != nil {
		return err
	}

	return s.scriptDB.Insert(script)
}

func (s *Service) UpdateScript(script *Script) error {
	if err := validate.Validate(script); err != nil {
		return err
	}

	return s.scriptDB.Update(&Script{ID: script.ID}, structutil.Struct2Map(script))
}

func (s *Service) DeleteScript(script *Script) error {
	return s.scriptDB.Delete(script)
}

func (s *Service) BatchDeleteScript(ids []uuid.UUID) error {
	tx := s.scriptDB.GetDB().Begin()
	defer tx.Rollback()

	scripts := make([]*Script, 0)
	for _, id := range ids {
		if id == uuid.Nil {
			return fmt.Errorf("invalid id: %s", id)
		}
		scripts = append(scripts, &Script{ID: id})
	}

	err := s.scriptDB.GetDB().Delete(scripts).Error

	if err != nil {
		logrus.Errorf("batch delete script failed, error: %v", err)
		return err
	}
	tx.Commit()
	return nil
}

func (s *Service) PageScript(num, size int, script *Script) (*database2.Pager[*Script], error) {
	return s.scriptDB.Page(script, int64(num), int64(size))
}

func (s *Service) DetailScript(id uuid.UUID) (*Script, error) {
	return s.scriptDB.Detail(&Script{ID: id})
}

func (s *Service) GetScriptFile(id uuid.UUID) (string, error) {
	if script, err := s.scriptDB.Detail(&Script{ID: id}); err != nil {
		return "", err
	} else {
		return script.Content, nil
	}
}

func (s *Service) RunScriptOnHosts(rsp *RunScriptParams) error {
	task := NewTask()
	psStr, _ := json.Marshal(rsp)
	task.SetParams(string(psStr))
	go task.Run()
	return nil
}

func (s *Service) PageRecord(num, size int, record *Record) (*database2.Pager[*Record], error) {
	if res, err := s.recordDB.Page(record, int64(num), int64(size)); err != nil {
		return nil, err
	} else {
		for _, r := range res.Data {
			r.Hosts = ""
		}
		return res, nil
	}
}

func (s *Service) StopScript(id uuid.UUID) error {
	if job, ok := jobMap.Load(id.String()); ok {
		job.(*JobInfo).ctxCancel()
		jobMap.Delete(id.String())
		return nil
	} else {
		return fmt.Errorf("task %s not found", id)
	}
}

func (s *Service) GetJobLog(id uuid.UUID) (map[string][]string, error) {
	if record, err := s.recordDB.Detail(&Record{ID: id}); err != nil {
		return nil, err
	} else {
		log := make(map[string][]string)
		if err := json.Unmarshal([]byte(record.Result), &log); err != nil {
			return nil, err
		}
		return log, nil
	}
}

func (s *Service) Initialize() error {
	if err := database2.GetDB().AutoMigrate(&Script{}, &Record{}); err != nil {
		return err
	}

	if err := schedule.GetExecutorManager().Register(schedule.Executor{
		Name:        "script",
		DisplayName: "script executor",
	}, func() schedule.Task {
		return NewTask()
	}); err != nil {
		return err
	}
	return nil
}
