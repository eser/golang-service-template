# bliss/configfx

## Overview

The **configfx** package provides a flexible and powerful configuration loader for Go applications. It supports loading configuration from various sources, including environment files, JSON files, and system environment variables. The package is designed to work seamlessly with the `go.uber.org/fx` framework.

The documentation below provides an overview of the package, its types, functions, and usage examples. For more detailed information, refer to the source code and tests.


## Fx

The `configfx` package provides an `FxModule` that can be used to integrate with the `go.uber.org/fx` framework.

```go
import (
  ...
	"github.com/eser/go-service/pkg/bliss/configfx"
	"go.uber.org/fx"
  ...
)

app := fx.New(
	configfx.FxModule,                 // registers configfx.ConfigLoader
	...
)

app.Run()
```


## API

### ConfigLoader interface
Defines methods for loading configuration.

```go
type ConfigLoader interface {
	LoadMeta(i any) (ConfigItemMeta, error)
	LoadMap(resources ...ConfigResource) (*map[string]any, error)
	Load(i any, resources ...ConfigResource) error

	FromEnvFileDirect(filename string) ConfigResource
	FromEnvFile(filename string) ConfigResource
	FromSystemEnv() ConfigResource

	FromJsonFileDirect(filename string) ConfigResource
	FromJsonFile(filename string) ConfigResource
}
```


### NewConfigLoader function

Creates a new `ConfigLoader` object based on the provided configuration.

```go
// func NewConfigLoader() *ConfigLoaderImpl

cl := logfx.NewConfigLoader()
```


### Load function

The `Load` method loads configuration from multiple resources.

Example:
```go
type AppConfig struct {
	AppName  string `conf:"NAME" default:"go-service"`

  Postgres struct {
		Dsn string `conf:"DSN" default:"postgres://localhost:5432"`
	} `conf:"POSTGRES"`
}

func loadConfig() (*AppConfig, error) {
  conf := &AppConfig{}

  cl := logfx.NewConfigLoader()

  err := cl.Load(
		conf,
                                      // load order:
		cl.FromJsonFile("config.json"),   // - attempts to read from config.json,
                                      //                         config.local.json,
                                      //                         config.[env].json,
                                      //                         config.[env].local.json
		cl.FromEnvFile(".env"),           // - attempts to read from .env
                                      //                         .env.local
                                      //                         .env.[env]
                                      //                         .env.[env].local
		cl.FromSystemEnv(),               // - attempts to read from system environment variables
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

  return conf, nil
}

func main() {
	appConfig, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)

    return
	}

// Searches JSON files first, then checks the POSTGRES__DSN among environment variables.
	// If the config variable is not specified, it falls back to the default value "postgres://localhost:5432".
  fmt.Println(appConfig.Postgres.Dsn)
}
```
