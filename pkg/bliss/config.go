package bliss

import (
	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/eser/go-service/pkg/bliss/datafx"
	"github.com/eser/go-service/pkg/bliss/grpcfx"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/logfx"
)

type BaseConfig struct {
	Data datafx.Config `conf:"DATA"`
	// AppName           string `conf:"APP_NAME" default:"go-service"`
	// JwtSignature      string `conf:"JWT_SIGNATURE"`
	// CorsOrigin        string `conf:"CORS_ORIGIN"`
	// CorsStrictHeaders bool   `conf:"CORS_STRICT_HEADERS"`
	// DataConnstr       string `conf:"DATA_CONNSTR"`
	Env string `conf:"ENV" default:"development"`

	Log  logfx.Config  `conf:"LOG"`
	Grpc grpcfx.Config `conf:"GRPC"`
	Http httpfx.Config `conf:"HTTP"`
}

func LoadConfig[T any](
	fnc func(cl configfx.ConfigLoader) (*T, error),
) func(configfx.ConfigLoader) (*T, error) {
	return func(cl configfx.ConfigLoader) (*T, error) {
		config, err := fnc(cl)
		if err != nil {
			return new(T), err
		}

		return config, nil
	}
}
