package samplehttp

import (
	"github.com/eser/go-service/pkg/bliss"
)

type AppConfig struct {
	AppName string `conf:"NAME" default:"samplehttp"`
	bliss.BaseConfig
}