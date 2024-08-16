package app

import (
	"time"

	"github.com/eser/go-service/pkg/bliss"
)

type AppConfig struct {
	bliss.BaseConfig

	AppName  string `conf:"NAME" default:"go-service"`
	Postgres struct {
		dsn string `conf:"DSN" default:"postgres://localhost:5432"`
	} `conf:"POSTGRES"`
	ReadTimeout time.Duration `conf:"READ_TIMEOUT" default:"5s"`
}
