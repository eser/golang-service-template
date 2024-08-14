package configfx

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/eser/go-service/pkg/bliss/configfx/envparser"
	"github.com/eser/go-service/pkg/bliss/configfx/jsonparser"
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

type ConfigLoader struct{}

type Config struct {
	Env string `conf:"env"`

	// AppName           string `conf:"APP_NAME"`
	// LogTarget         string `conf:"LOG_TARGET"`
	// JwtSignature      string `conf:"JWT_SIGNATURE"`
	// Port              string `conf:"PORT"`
	// CorsOrigin        string `conf:"CORS_ORIGIN"`
	// CorsStrictHeaders bool   `conf:"CORS_STRICT_HEADERS"`
	// DataConnstr       string `conf:"DATA_CONNSTR"`
}

func GetCurrentEnv() string {
	// FIXME(@eser) no need to use os.Lookupenv here
	env := strings.ToLower(os.Getenv("ENV"))

	if env == "" {
		env = "development"
	}

	return env
}

func splitFilename(filename string) (string, string, string) {
	dir, file := filepath.Split(filename)
	ext := filepath.Ext(file)
	rest := len(file) - len(ext)

	if rest == 0 {
		return dir, file, ""
	}

	return dir, file[:rest], ext
}

func GetFilenamesForEnv(env string, filename string) []string {
	dirname, basename, ext := splitFilename(filename)

	filenames := []string{
		filename,
		dirname + basename + "." + env + ext,
	}

	if env != "test" {
		filenames = append(filenames, dirname+basename+".local"+ext)
	}

	filenames = append(filenames, dirname+basename+"."+env+".local"+ext)

	return filenames
}

func OverrideWithSystemEnv(m *map[string]string) {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2) //nolint:gomnd,mnd
		(*m)[pair[0]] = pair[1]
	}
}

func (dcl *ConfigLoader) TryLoadEnv(m *map[string]string) error {
	env := GetCurrentEnv()
	filenames := GetFilenamesForEnv(env, ".env")

	err := envparser.TryParseFiles(m, filenames...)
	if err != nil {
		return err //nolint:wrapcheck
	}

	OverrideWithSystemEnv(m)

	return nil
}

func (dcl *ConfigLoader) TryLoadJson(m *map[string]any) error {
	env := GetCurrentEnv()
	filenames := GetFilenamesForEnv(env, "config.json")

	err := jsonparser.TryParseFiles(m, filenames...)
	if err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}

func (dcl *ConfigLoader) Load(out *any) error {
	return nil
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
