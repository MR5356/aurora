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
	Server      Server                 `json:"server" yaml:"server"`
	Database    Database               `json:"database" yaml:"database"`
	JWT         JWT                    `json:"jwt" yaml:"jwt"`
	OAuthConfig map[string]OAuthConfig `json:"oauth" yaml:"oauth"`
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

type Server struct {
	BaseURL     string `json:"baseURL" yaml:"baseURL" default:"http://localhost"`
	Port        int    `json:"port" yaml:"port" default:"80"`
	Prefix      string `json:"prefix" yaml:"prefix" default:"/api/v1"`
	Debug       bool   `json:"debug" yaml:"debug" default:"false"`
	GracePeriod int    `json:"gracePeriod" yaml:"gracePeriod" default:"30"`
}

type JWT struct {
	Secret string        `json:"secret" yaml:"secret" default:"aurora"`
	Issuer string        `json:"issuer" yaml:"issuer" default:"fun.toodo.aurora"`
	Expire time.Duration `json:"expire" yaml:"expire" default:"720h"`
}

type OAuthConfig struct {
	AuthType     string `json:"authType" yaml:"authType"`
	AuthURL      string `json:"authURL" yaml:"authURL"`
	TokenURL     string `json:"tokenURL" yaml:"tokenURL"`
	ClientId     string `json:"clientId" yaml:"clientId"`
	ClientSecret string `json:"clientSecret" yaml:"clientSecret"`
}

type Database struct {
	Driver string `json:"driver" yaml:"driver" default:"sqlite"`
	DSN    string `json:"dsn" yaml:"dsn" default:"db.sqlite"`

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
