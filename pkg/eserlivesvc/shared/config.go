package shared

import (
	"github.com/eser/go-service/pkg/bliss"
)

type AppConfig struct {
	AppName string `conf:"NAME" default:"eserlivesvc"`
	bliss.BaseConfig
}
