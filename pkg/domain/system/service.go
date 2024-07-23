package system

import (
	"context"
	"fmt"
	"github.com/MR5356/aurora/pkg/domain/health"
	"github.com/MR5356/aurora/pkg/domain/host"
	"github.com/MR5356/aurora/pkg/domain/schedule"
	"github.com/MR5356/aurora/pkg/domain/user"
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/MR5356/aurora/pkg/version"
	"github.com/google/go-github/v61/github"
	"github.com/spf13/cast"
	"sync"
	"time"
)

var (
	once    sync.Once
	service *Service
)

type Service struct {
	recordDB   *database.BaseMapper[*Record]
	userDB     *database.BaseMapper[*user.User]
	scheduleDB *database.BaseMapper[*schedule.Schedule]
	hostDB     *database.BaseMapper[*host.Host]
	healthDB   *database.BaseMapper[*health.Health]
}

func GetService() *Service {
	once.Do(func() {
		service = &Service{
			recordDB:   database.NewMapper(database.GetDB(), &Record{}),
			userDB:     database.NewMapper(database.GetDB(), &user.User{}),
			scheduleDB: database.NewMapper(database.GetDB(), &schedule.Schedule{}),
			hostDB:     database.NewMapper(database.GetDB(), &host.Host{}),
			healthDB:   database.NewMapper(database.GetDB(), &health.Health{}),
		}
	})
	return service
}

func (s *Service) InsertRecord(record *Record) error {
	return s.recordDB.Insert(record)
}

func (s *Service) GetVersionInfo() *Version {
	result := &Version{
		Version: version.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	client := github.NewClient(nil)
	res, _, err := client.Repositories.GetLatestRelease(ctx, "MR5356", "aurora")
	if err != nil {
		return result
	}
	result.LatestInfo = cast.ToString(*res.Body)
	result.LatestVersion = cast.ToString(*res.TagName)
	result.LatestUrl = cast.ToString(*res.HTMLURL)

	return result
}

func (s *Service) GetStatistic() ([]*Statistic, error) {
	res := make([]*Statistic, 0)
	// user
	uc, err := s.userDB.Count(&user.User{})
	if err != nil {
		return res, err
	}

	// schedule
	scAll, err := s.scheduleDB.Count(&schedule.Schedule{})
	if err != nil {
		return res, err
	}

	sc, err := s.scheduleDB.Count(&schedule.Schedule{Enabled: true})
	if err != nil {
		return res, err
	}

	// record
	var rc int64
	err = s.recordDB.DB.Model(&Record{}).Where("is_api = ?", true).Count(&rc).Error
	if err != nil {
		return res, err
	}

	// host
	rh, err := s.hostDB.Count(&host.Host{})
	if err != nil {
		return res, err
	}

	// health
	rHealth, err := s.healthDB.Count(&health.Health{})
	if err != nil {
		return res, err
	}

	res = append(res, &Statistic{
		Name:  "statistic.user",
		Count: fmt.Sprintf("%d", uc),
		Icon:  "user",
	})

	res = append(res, &Statistic{
		Name:  "statistic.host",
		Count: fmt.Sprintf("%d", rh),
		Path:  "/host",
		Icon:  "host",
	})

	res = append(res, &Statistic{
		Name:  "statistic.terminal",
		Count: "i18n://terminal.info",
		Path:  "/terminal",
		Icon:  "terminal",
	})

	res = append(res, &Statistic{
		Name:  "statistic.health",
		Count: fmt.Sprintf("%d", rHealth),
		Path:  "/health",
		Icon:  "health",
	})

	res = append(res, &Statistic{
		Name:  "statistic.schedule",
		Count: fmt.Sprintf("%d/%d", sc, scAll),
		Path:  "/schedule/list",
		Icon:  "schedule",
	})

	res = append(res, &Statistic{
		Name:  "statistic.record",
		Count: fmt.Sprintf("%d", rc),
		Icon:  "record",
	})
	return res, nil
}

func (s *Service) Initialize() error {
	if err := database.GetDB().AutoMigrate(&Record{}); err != nil {
		return err
	}
	return nil
}
