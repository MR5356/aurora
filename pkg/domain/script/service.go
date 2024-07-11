package script

import (
	"fmt"
	"github.com/MR5356/aurora/pkg/middleware/database"
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
	scriptDB *database.BaseMapper[*Script]
}

func GetService() *Service {
	onceService.Do(func() {
		service = &Service{
			scriptDB: database.NewMapper(database.GetDB(), &Script{}),
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
	tx := s.scriptDB.DB.Begin()
	defer tx.Rollback()

	scripts := make([]*Script, 0)
	for _, id := range ids {
		if id == uuid.Nil {
			return fmt.Errorf("invalid id: %s", id)
		}
		scripts = append(scripts, &Script{ID: id})
	}

	err := s.scriptDB.DB.Delete(scripts).Error

	if err != nil {
		logrus.Errorf("batch delete script failed, error: %v", err)
		return err
	}
	tx.Commit()
	return nil
}

func (s *Service) PageScript(num, size int, script *Script) (*database.Pager[*Script], error) {
	return s.scriptDB.Page(script, int64(num), int64(size))
}

func (s *Service) DetailScript(id uuid.UUID) (*Script, error) {
	return s.scriptDB.Detail(&Script{ID: id})
}

func (s *Service) Initialize() error {
	if err := database.GetDB().AutoMigrate(&Script{}); err != nil {
		return err
	}
	return nil
}
