package config

import (
	"github.com/eser/ajan"
)

type AppConfig struct {
	AppName string `conf:"NAME" default:"samplesvc"`
	ajan.BaseConfig
}
