package app

import (
	"time"

	"github.com/eser/go-service/pkg/bliss/configfx"
)

type AppConfig struct {
	configfx.BaseConfig

	AppName  string `conf:"name" default:"go-service"`
	Postgres struct {
		dsn string `conf:"dsn" default:"postgres://localhost:5432"`
	} `conf:"postgres"`
	ReadTimeout time.Duration `conf:"read_timeout" default:"5s"`
}
