package health

import (
	"encoding/json"
	"errors"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/MR5356/health"
	"github.com/MR5356/health/database"
	"github.com/MR5356/health/host"
	"github.com/MR5356/health/url"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

const (
	defaultCron     = "*/10 * * * * *"
	defaultPingCron = "*/1 * * * * *"
	defaultHttpCron = "*/2 * * * * *"
	defaultSSHCron  = "*/10 * * * * *"
	defaultDBCron   = "*/10 * * * * *"
)

func getCron(t string) string {
	switch t {
	case "ping":
		return defaultPingCron
	case "http":
		return defaultHttpCron
	case "ssh":
		return defaultSSHCron
	case "database":
		return defaultDBCron
	default:
		return defaultCron
	}
}

func (s *Service) initChecker() error {
	if healths, err := s.healthDb.List(&Health{Enabled: true}); err != nil {
		return err
	} else {
		for _, health := range healths {
			if err := s.startChecker(health); err != nil {
				logrus.Errorf("start checker failed, error: %v", err)
			}
		}
	}
	return nil
}

type Checker struct {
	health  *Health
	service *Service
}

func (c *Checker) Run() {
	logrus.Debugf("health check: %s", c.health.Title)
	var checker health.Checker
	var params Params
	if err := json.Unmarshal([]byte(c.health.Params), &params); err != nil {
		logrus.Errorf("unmarshal params failed, error: %v", err)
		return
	}
	switch c.health.Type {
	case "http":
		if str, ok := params.GetKey("url").(string); ok {
			checker = url.NewChecker(str)
		} else {
			logrus.Errorf("http url is empty")
			return
		}
	case "ssh":
		privateKey, ok := params.GetKey("privateKey").(string)
		if !ok {
			logrus.Errorf("ssh private key is empty")
			return
		}
		passphrase, ok := params.GetKey("passphrase").(string)
		if !ok {
			logrus.Errorf("ssh passphrase is empty")
			return
		}
		hostStr, ok := params.GetKey("host").(string)
		if !ok {
			logrus.Errorf("ssh host is empty")
			return
		}
		port, ok := params.GetKey("port").(float64)
		if !ok {
			logrus.Errorf("ssh port is empty")
			return
		}
		username, ok := params.GetKey("username").(string)
		if !ok {
			logrus.Errorf("ssh username is empty")
			return
		}
		password, ok := params.GetKey("password").(string)
		if !ok {
			logrus.Errorf("ssh password is empty")
			return
		}
		checker = host.NewSSHChecker(&host.HostInfo{
			PrivateKey: privateKey,
			Passphrase: passphrase,
			Host:       hostStr,
			Port:       uint16(port),
			Username:   username,
			Password:   password,
		})
	case "ping":
		if str, ok := params.GetKey("host").(string); ok {
			checker = host.NewPingChecker(str)
		} else {
			logrus.Errorf("ping host is empty")
			return
		}
	case "database":
		dbType, ok := params.GetKey("dbDriverType").(string)
		if !ok {
			logrus.Errorf("database dbDriverType is empty")
			return
		}
		dsn, ok := params.GetKey("dsn").(string)
		if !ok {
			logrus.Errorf("database dsn is empty")
			return
		}
		checker = database.NewChecker(dbType, dsn)
	default:
		logrus.Errorf("unknown health check type: %s", c.health.Type)
		return
	}
	res := checker.Check()
	c.health.Status = string(res.Status)
	c.health.RTT = res.RTT
	if err := c.service.healthDb.Update(&Health{ID: c.health.ID}, structutil.Struct2Map(c.health)); err != nil {
		logrus.Errorf("update health failed, error: %v", err)
	}

	result, _ := json.Marshal(res.Result)
	healthRecord := &Record{
		ParentId: c.health.ID,
		Status:   string(res.Status),
		Rtt:      res.RTT,
		Result:   string(result),
	}
	if err := c.service.healthRecordDb.Insert(healthRecord); err != nil {
		logrus.Errorf("insert health record failed, error: %v", err)
	}
}

func (s *Service) startChecker(health *Health) error {
	if _, ok := s.cronJobMap.Load(health.ID); ok {
		return errors.New("cron job already exists")
	}
	if jobId, err := s.cron.AddJob(getCron(health.Type), &Checker{health: health, service: s}); err != nil {
		return err
	} else {
		s.cronJobMap.Store(health.ID, jobId)
		return nil
	}
}

func (s *Service) stopChecker(health *Health) error {
	if jobId, ok := s.cronJobMap.Load(health.ID); !ok {
		return errors.New("cron job not found")
	} else {
		s.cron.Remove(jobId.(cron.EntryID))
		s.cronJobMap.Delete(health.ID)
		return nil
	}
}
