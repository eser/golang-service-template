package appcontext

import (
	"github.com/eser/ajan"
)

type FeatureFlags struct {
	Dummy bool `conf:"DUMMY" default:"false"` // dummy feature flag
}

type AppConfig struct {
	ajan.BaseConfig

	Features FeatureFlags `conf:"FEATURES"`
}
