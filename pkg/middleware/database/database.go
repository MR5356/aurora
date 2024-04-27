package database

import (
	"encoding/json"
	"errors"
	"github.com/MR5356/aurora/pkg/config"
	"github.com/avast/retry-go"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

var (
	once sync.Once
	db   *Database
)

type Database struct {
	*gorm.DB
}

var (
	DBDriverNotSupport = errors.New("database driver not support")
)

func GetDB() *Database {
	once.Do(func() {
		err := retry.Do(
			func() (err error) {
				db, err = initDB()
				return err
			},
			retry.Attempts(config.DefaultRetryCount),
			retry.Delay(config.DefaultRetryDelay),
			retry.LastErrorOnly(true),
			retry.DelayType(retry.DefaultDelayType),
			retry.OnRetry(func(n uint, err error) {
				logrus.Warnf("[%d/%d]: retry to initialize database: %v", n+1, config.DefaultRetryCount, err)
			}),
		)
		if err != nil {
			logrus.Fatalf("Failed to initialize database: %v", err)
		}
	})
	return db
}

func initDB() (database *Database, err error) {
	cfg := config.Current()
	var driver gorm.Dialector
	logrus.Debugf("database driver: %s", cfg.Database.Driver)
	switch cfg.Database.Driver {
	case "sqlite":
		driver = sqlite.Open(cfg.Database.DSN)
	case "mysql":
		driver = mysql.Open(cfg.Database.DSN)
	case "postgres":
		driver = postgres.Open(cfg.Database.DSN)
	default:
		return nil, DBDriverNotSupport
	}

	var dbLogLevel = logger.Error
	if cfg.Server.Debug {
		dbLogLevel = logger.Info
	}
	logrus.Debugf("database log level: %+v", dbLogLevel)

	client, err := gorm.Open(driver, &gorm.Config{
		Logger: logger.Default.LogMode(dbLogLevel),
	})
	if err != nil {
		return nil, err
	}

	db, err := client.DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.Database.MaxIdleConn)
	db.SetMaxOpenConns(cfg.Database.MaxOpenConn)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLift)
	db.SetConnMaxIdleTime(cfg.Database.ConnMaxIdle)

	dbStat, _ := json.Marshal(db.Stats())
	logrus.Debugf("database stats: %s", dbStat)
	return &Database{client}, nil
}
