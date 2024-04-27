package config

import (
	"github.com/mcuadros/go-defaults"
	"time"
)

const (
	DefaultRetryCount = 3
	DefaultRetryDelay = time.Second * 0
)

var config *Config

type Config struct {
	Server   Server   `json:"server" yaml:"server"`
	Database Database `json:"database" yaml:"database"`
}

type Server struct {
	Port        int    `json:"port" yaml:"port" default:"8888"`
	Prefix      string `json:"prefix" yaml:"prefix" default:"/api/v1"`
	Debug       bool   `json:"debug" yaml:"debug" default:"false"`
	GracePeriod int    `json:"gracePeriod" yaml:"gracePeriod" default:"30"`
}

func Current(cfgs ...Cfg) *Config {
	if config == nil {
		config = New(cfgs...)
	}
	return config
}

func New(cfgs ...Cfg) *Config {
	config = new(Config)
	defaults.SetDefaults(config)

	for _, cfg := range cfgs {
		cfg(config)
	}

	return config
}

type Database struct {
	Driver string `json:"driver" yaml:"driver" default:"sqlite"`
	DSN    string `json:"dsn" yaml:"dsn" default:":memory:"`

	MaxIdleConn int           `json:"maxIdleConn" yaml:"maxIdleConn" default:"10"`
	MaxOpenConn int           `json:"maxOpenConn" yaml:"maxOpenConn" default:"40"`
	ConnMaxLift time.Duration `json:"connMaxLift" yaml:"connMaxLift" default:"0s"`
	ConnMaxIdle time.Duration `json:"connMaxIdle" yaml:"connMaxIdle" default:"0s"`
}

type Cfg func(c *Config)

func WithPort(port int) Cfg {
	return func(c *Config) {
		c.Server.Port = port
	}
}

func WithDatabase(driver, dsn string) Cfg {
	return func(c *Config) {
		c.Database.Driver = driver
		c.Database.DSN = dsn
	}
}

func WithDebug(debug bool) Cfg {
	return func(c *Config) {
		c.Server.Debug = debug
	}
}

func WithGracePeriod(gracePeriod int) Cfg {
	return func(c *Config) {
		c.Server.GracePeriod = gracePeriod
	}
}
