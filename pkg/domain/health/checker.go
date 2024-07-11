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
	"github.com/spf13/cast"
)

const (
	defaultCron     = "*/10 * * * * *"
	defaultPingCron = "*/1 * * * * *"
	defaultHttpCron = "*/2 * * * * *"
	defaultSSHCron  = "*/10 * * * * *"
	defaultDBCron   = "*/10 * * * * *"
)

var (
	ErrParam    = errors.New("invalid param")
	StatusError = health.Status("error")
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
	c.health.RTT = 0
	c.health.Status = string(health.StatusUnknown)
	defer func() {
		if err := c.service.healthDb.Update(&Health{ID: c.health.ID}, structutil.Struct2Map(c.health)); err != nil {
			logrus.Errorf("update health failed, error: %v", err)
		}
	}()
	var checker health.Checker
	var params Params
	if err := json.Unmarshal([]byte(c.health.Params), &params); err != nil {
		logrus.Errorf("unmarshal params failed, error: %v", err)
		c.health.Status = string(StatusError)
		return
	}
	switch c.health.Type {
	case "http":
		if str, err := cast.ToStringE(params.GetKey("url")); err != nil {
			logrus.Errorf("http url is empty")
			c.health.Status = string(StatusError)
			return
		} else {
			checker = url.NewChecker(str)
		}
	case "ssh":
		privateKey, err := cast.ToStringE(params.GetKey("privateKey"))
		if err != nil {
			logrus.Errorf("ssh private key is empty")
			c.health.Status = string(StatusError)
			return
		}
		passphrase, err := cast.ToStringE(params.GetKey("passphrase"))
		if err != nil {
			logrus.Errorf("ssh passphrase is empty")
			c.health.Status = string(StatusError)
			return
		}
		hostStr, err := cast.ToStringE(params.GetKey("host"))
		if err != nil {
			logrus.Errorf("ssh host is empty")
			c.health.Status = string(StatusError)
			return
		}
		port, err := cast.ToUint16E(params.GetKey("port"))
		if err != nil {
			logrus.Errorf("ssh port is empty")
			c.health.Status = string(StatusError)
			return
		}
		username, err := cast.ToStringE(params.GetKey("username"))
		if err != nil {
			logrus.Errorf("ssh username is empty")
			c.health.Status = string(StatusError)
			return
		}
		password, err := cast.ToStringE(params.GetKey("password"))
		if err != nil {
			logrus.Errorf("ssh password is empty")
			c.health.Status = string(StatusError)
			return
		}
		checker = host.NewSSHChecker(&host.HostInfo{
			PrivateKey: privateKey,
			Passphrase: passphrase,
			Host:       hostStr,
			Port:       port,
			Username:   username,
			Password:   password,
		})
	case "ping":
		if str, err := cast.ToStringE(params.GetKey("host")); err == nil {
			checker = host.NewPingChecker(str)
		} else {
			logrus.Errorf("ping host is empty")
			c.health.Status = string(StatusError)
			return
		}
	case "database":
		dbType, err := cast.ToStringE(params.GetKey("dbDriverType"))
		if err != nil {
			logrus.Errorf("database dbDriverType is empty")
			c.health.Status = string(StatusError)
			return
		}
		dsn, err := cast.ToStringE(params.GetKey("dsn"))
		if err != nil {
			logrus.Errorf("database dsn is empty")
			c.health.Status = string(StatusError)
			return
		}
		checker = database.NewChecker(dbType, dsn)
	default:
		logrus.Errorf("unknown health check type: %s", c.health.Type)
		c.health.Status = string(StatusError)
		return
	}
	res := checker.Check()
	c.health.Status = string(res.Status)
	c.health.RTT = res.RTT

	//result, _ := json.Marshal(res.Result)
	//healthRecord := &Record{
	//	ParentId: c.health.ID,
	//	Status:   string(res.Status),
	//	Rtt:      res.RTT,
	//	Result:   string(result),
	//}
	//if err := c.service.healthRecordDb.Insert(healthRecord); err != nil {
	//	logrus.Errorf("insert health record failed, error: %v", err)
	//}
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
