package host

import (
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/MR5356/aurora/pkg/util/validate"
	"github.com/google/uuid"
	"sync"
)

var (
	onceService sync.Once
	service     *Service
)

type Service struct {
	hostDb  *database.BaseMapper[*Host]
	groupDb *database.BaseMapper[*Group]
}

func GetService() *Service {
	onceService.Do(func() {
		service = &Service{
			hostDb:  database.NewMapper(database.GetDB(), &Host{}),
			groupDb: database.NewMapper(database.GetDB(), &Group{}),
		}
	})
	return service
}

// ListGroup list host group
func (s *Service) ListGroup(group *Group) ([]*Group, error) {
	res := make([]*Group, 0)
	if err := s.groupDb.DB.Preload("Hosts").Find(&res, group).Error; err != nil {
		return res, err
	}

	for _, g := range res {
		for _, h := range g.Hosts {
			h.HostInfo.Password = ""
		}
	}

	return res, nil
}

// AddGroup add host group
func (s *Service) AddGroup(group *Group) error {
	group.ID = uuid.Nil
	if err := validate.Validate(group); err != nil {
		return err
	}
	return s.groupDb.Insert(group)
}

// DeleteGroup delete host group
func (s *Service) DeleteGroup(id uuid.UUID) error {
	if err := s.groupDb.DB.Where("group_id = ?", id).Delete(&Host{GroupId: id}).Error; err != nil {
		return err
	}

	return s.groupDb.Delete(&Group{ID: id})
}

func (s *Service) UpdateGroup(group *Group) error {
	if err := validate.Validate(group); err != nil {
		return err
	}
	return s.groupDb.DB.Updates(group).Error
}

func (s *Service) Initialize() error {
	if err := database.GetDB().AutoMigrate(&Host{}, &Group{}); err != nil {
		return err
	}

	if err := s.groupDb.DB.FirstOrCreate(&Group{}, &Group{Title: "default"}).Error; err != nil {
		return err
	}
	return nil
}
