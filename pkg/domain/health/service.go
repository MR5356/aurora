package health

import (
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/MR5356/aurora/pkg/util/validate"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"sync"
	"time"
)

var (
	onceService sync.Once
	service     *Service
)

type Service struct {
	healthDb       *database.BaseMapper[*Health]
	healthRecordDb *database.BaseMapper[*Record]
	cron           *cron.Cron
	cronJobMap     sync.Map
}

func GetService() *Service {
	onceService.Do(func() {
		c := cron.New(cron.WithSeconds())
		c.Start()
		service = &Service{
			healthDb:       database.NewMapper(database.GetDB(), &Health{}),
			healthRecordDb: database.NewMapper(database.GetDB(), &Record{}),
			cron:           c,
			cronJobMap:     sync.Map{},
		}
	})
	return service
}

func (s *Service) ListHealth(health *Health) ([]*Health, error) {
	return s.healthDb.List(health)
}

func (s *Service) AddHealth(health *Health) error {
	health.ID = uuid.Nil
	if err := validate.Validate(health); err != nil {
		return err
	}
	if err := s.healthDb.Insert(health); err != nil {
		return err
	}

	return s.startChecker(health)
}

func (s *Service) UpdateHealth(health *Health) error {
	if err := validate.Validate(health); err != nil {
		return err
	}
	if err := s.healthDb.Update(&Health{ID: health.ID}, structutil.Struct2Map(health)); err != nil {
		return err
	}

	if err := s.stopChecker(health); err != nil {
		return err
	}

	return s.startChecker(health)
}

func (s *Service) DeleteHealth(health *Health) error {
	if err := s.stopChecker(health); err != nil {
		return err
	}
	return s.healthDb.Delete(health)
}

func (s *Service) DetailHealth(id uuid.UUID) (*Health, error) {
	return s.healthDb.Detail(&Health{ID: id})
}

func (s *Service) GetTimeRangeRecord(healthId uuid.UUID, startTime time.Time, endTime time.Time) ([]*Record, error) {
	res := make([]*Record, 0)
	if err := s.healthRecordDb.DB.Where("parent_id = ?", healthId).Where("created_at BETWEEN ? AND ?", startTime, endTime).Find(&res).Error; err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (s *Service) GetHealthCheckTypes() []CheckType {
	return []CheckType{
		{
			Title: "http",
			Type:  "http",
			Desc:  "http check, support http and https GET method",
			Params: []Param{
				{
					Key:      "url",
					Value:    "",
					Title:    "URL",
					Type:     "string",
					Required: true,
					Desc:     "url for check",
				},
			},
		},
		{
			Title: "ssh",
			Type:  "ssh",
			Desc:  "ssh check",
			Params: []Param{
				{
					Key:      "host",
					Value:    "",
					Title:    "Host",
					Type:     "string",
					Required: true,
					Desc:     "host or ip",
				},
				{
					Key:      "port",
					Value:    "22",
					Title:    "Port",
					Type:     "number",
					Required: true,
					Desc:     "port",
				},
				{
					Key:      "username",
					Value:    "",
					Title:    "Username",
					Type:     "string",
					Required: true,
					Desc:     "username",
				},
				{
					Key:      "password",
					Value:    "",
					Title:    "Password",
					Type:     "string",
					Required: false,
					Desc:     "password",
				},
				{
					Key:      "privateKey",
					Value:    "",
					Title:    "PrivateKey",
					Type:     "string",
					Required: false,
					Desc:     "PrivateKey",
				},
				{
					Key:      "passphrase",
					Value:    "",
					Title:    "PrivateKey passphrase",
					Type:     "string",
					Required: false,
					Desc:     "privateKey passphrase",
				},
			},
		},
		{
			Title: "ping",
			Type:  "ping",
			Desc:  "ping check",
			Params: []Param{
				{
					Key:      "host",
					Value:    "",
					Title:    "Host",
					Type:     "string",
					Required: true,
					Desc:     "host or ip",
				},
			},
		},
		{
			Title: "database",
			Type:  "database",
			Desc:  "database check, support mysql sqlite3 and postgresql",
			Params: []Param{
				{
					Key:      "dbDriverType",
					Value:    "",
					Title:    "Database Driver Type",
					Type:     "string",
					Required: true,
					Desc:     "Database Driver Type, support mysql sqlite3 and postgresql",
				},
				{
					Key:      "dsn",
					Value:    "",
					Title:    "Dsn",
					Type:     "string",
					Required: true,
					Desc:     "Connection information",
				},
			},
		},
	}
}

func (s *Service) Initialize() error {
	if err := database.GetDB().AutoMigrate(&Health{}, &Record{}); err != nil {
		return err
	}
	if err := s.initChecker(); err != nil {
		return err
	}
	return nil
}
