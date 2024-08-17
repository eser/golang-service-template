package bliss

import (
	"github.com/eser/go-service/pkg/bliss/httpfx"
)

type BaseConfig struct {
	Env string `conf:"ENV" default:"development"`

	Http httpfx.Config `conf:"HTTP"`

	// AppName           string `conf:"APP_NAME" default:"go-service"`
	// LogTarget         string `conf:"LOG_TARGET" default:"stdout"`
	// Port              string `conf:"PORT" default:"8080"`
	// JwtSignature      string `conf:"JWT_SIGNATURE"`
	// CorsOrigin        string `conf:"CORS_ORIGIN"`
	// CorsStrictHeaders bool   `conf:"CORS_STRICT_HEADERS"`
	// DataConnstr       string `conf:"DATA_CONNSTR"`
}
