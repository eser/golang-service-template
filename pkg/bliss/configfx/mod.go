package configfx

import (
	"errors"

	"go.uber.org/fx"
)

var ErrConfigDecoding = errors.New("config decoding error")

//nolint:gochecknoglobals
var Module = fx.Module(
	"config",
	fx.Provide(
		New,
	),
)

type Result struct {
	fx.Out

	ConfigLoader *ConfigLoader
}

type BaseConfig struct {
	Env string `conf:"env"`

	// AppName           string `conf:"APP_NAME"`
	// LogTarget         string `conf:"LOG_TARGET"`
	// JwtSignature      string `conf:"JWT_SIGNATURE"`
	// Port              string `conf:"PORT"`
	// CorsOrigin        string `conf:"CORS_ORIGIN"`
	// CorsStrictHeaders bool   `conf:"CORS_STRICT_HEADERS"`
	// DataConnstr       string `conf:"DATA_CONNSTR"`
}

func New() (Result, error) {
	// envMap := map[string]string{
	// 	"APP_NAME":   "go-service",
	// 	"LOG_TARGET": "stdout",
	// 	"PORT":       "8080",
	// }

	// err := TryLoadEnv(&envMap)
	// if err != nil {
	// 	return nil, err
	// }

	// config := Config{}
	// // err = mapstructure.Decode(envMap, &config)
	// // if err != nil {
	// // 	return nil, fmt.Errorf("config decoding error: %w %w", ErrConfigDecoding, err)
	// // }

	return Result{
		ConfigLoader: &ConfigLoader{},
	}, nil
}
