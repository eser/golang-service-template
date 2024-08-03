package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/fx"
)

//nolint:gochecknoglobals
var Module = fx.Module(
	"config",
	fx.Provide(
		NewConfig,
	),
)

type Config struct {
	Env string `mapstructure:"ENV"`

	AppName           string `mapstructure:"APP_NAME"`
	LogTarget         string `mapstructure:"LOG_TARGET"`
	JwtSignature      string `mapstructure:"JWT_SIGNATURE"`
	Port              string `mapstructure:"PORT"`
	CorsOrigin        string `mapstructure:"CORS_ORIGIN"`
	CorsStrictHeaders bool   `mapstructure:"CORS_STRICT_HEADERS"`
	DataConnstr       string `mapstructure:"DATA_CONNSTR"`
}

func OverrideWithSystemEnv(envMap *map[string]string) {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2) //nolint:gomnd
		(*envMap)[pair[0]] = pair[1]
	}
}

func TryLoadEnv(envMap *map[string]string) error {
	env := strings.ToLower(os.Getenv("ENV"))

	if env == "" {
		env = "development"
	}

	if env != "test" {
		TryParseFiles(envMap, ".env", ".env."+env, ".env.local", ".env."+env+".local")
	} else {
		TryParseFiles(envMap, ".env", ".env."+env, ".env."+env+".local")
	}

	OverrideWithSystemEnv(envMap)

	(*envMap)["ENV"] = env

	return nil
}

func NewConfig() (*Config, error) {
	envMap := map[string]string{
		"APP_NAME":   "go-service",
		"LOG_TARGET": "stdout",
		"PORT":       "8080",
	}

	err := TryLoadEnv(&envMap)
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = mapstructure.Decode(envMap, &config)
	if err != nil {
		return nil, fmt.Errorf("config decoding error: %w", err)
	}

	return &config, nil
}
