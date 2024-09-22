package bliss

import (
	"reflect"

	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/logfx"
	"go.uber.org/fx"
)

type BaseConfig struct {
	Env string `conf:"ENV" default:"development"`

	Log  logfx.Config  `conf:"LOG"`
	Http httpfx.Config `conf:"HTTP"`

	// AppName           string `conf:"APP_NAME" default:"go-service"`
	// JwtSignature      string `conf:"JWT_SIGNATURE"`
	// CorsOrigin        string `conf:"CORS_ORIGIN"`
	// CorsStrictHeaders bool   `conf:"CORS_STRICT_HEADERS"`
	// DataConnstr       string `conf:"DATA_CONNSTR"`
}

type RegisterConfigResult[T any] struct {
	fx.Out

	Config     *T
	BaseConfig *BaseConfig
	LogConfig  *logfx.Config
	HttpConfig *httpfx.Config
}

func LoadConfig[T any](
	fnc func(cl configfx.ConfigLoader) (*T, error),
) func(configfx.ConfigLoader) (RegisterConfigResult[T], error) {
	return func(cl configfx.ConfigLoader) (RegisterConfigResult[T], error) {
		config, err := fnc(cl)
		if err != nil {
			return RegisterConfigResult[T]{}, err
		}

		return RegisterConfig(config), nil
	}
}

func RegisterConfig[T any](config *T) RegisterConfigResult[T] {
	rConfig := reflect.ValueOf(config).Elem()
	rBaseConfig := rConfig.FieldByName("BaseConfig")

	if !rBaseConfig.IsValid() {
		panic("config must have BaseConfig field")
	}

	baseConfig, ok := rBaseConfig.Interface().(BaseConfig)

	if !ok {
		panic("BaseConfig field must be of type BaseConfig")
	}

	return RegisterConfigResult[T]{ //nolint:exhaustruct
		Config:     config,
		BaseConfig: &baseConfig,
		LogConfig:  &baseConfig.Log,
		HttpConfig: &baseConfig.Http,
	}
}
