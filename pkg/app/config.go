package app

import (
	"time"

	"github.com/eser/go-service/pkg/bliss"
)

type AppConfig struct {
	bliss.BaseConfig

	AppName  string `conf:"name" default:"go-service"`
	Postgres struct {
		dsn string `conf:"dsn" default:"postgres://localhost:5432"`
	} `conf:"postgres"`
	ReadTimeout time.Duration `conf:"read_timeout" default:"5s"`
}
