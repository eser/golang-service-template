# bliss/httpfx

## Overview

The **httpfx** package provides a framework for building HTTP services with support for routing, middleware, and OpenAPI
documentation generation. The package is designed to work seamlessly with the `bliss/di` package.

The documentation below provides an overview of the package, its types, functions, and usage examples. For more detailed
information, refer to the source code and tests.

## Configuration

Configuration struct for the logger:

```go
type Config struct {
	ReadHeaderTimeout time.Duration `conf:"READ_HEADER_TIMEOUT" default:"5s"`
	ReadTimeout       time.Duration `conf:"READ_TIMEOUT"        default:"10s"`
	WriteTimeout      time.Duration `conf:"WRITE_TIMEOUT"       default:"10s"`
	IdleTimeout       time.Duration `conf:"IDLE_TIMEOUT"        default:"120s"`

	InitializationTimeout   time.Duration `conf:"INIT_TIMEOUT"     default:"25s"`
	GracefulShutdownTimeout time.Duration `conf:"SHUTDOWN_TIMEOUT" default:"5s"`

	Addr string `conf:"ADDR" default:":8080"`
}
```

## Bliss DI

The `httpfx` package provides a `RegisterDependencies` function that can be used to integrate with the `bliss/di`
package.

```go
import (
  ...
  "github.com/eser/go-service/pkg/bliss/di"
	"github.com/eser/go-service/pkg/bliss/httpfx"
  ...
)

err := di.RegisterFn(
	httpfx.RegisterDependencies,                    // registers httpfx.HttpService and httpfx.Router
	...
)
```

## API

### NewRouter function

Create a new `Router` object.

```go
// func NewRouter(path string) *RouterImpl

router := httpfx.NewRouter("/")
```

### NewHttpService function

Creates a new `HttpService` object based on the provided configuration.

```go
// func NewHttpService(config *Config, router Router) *HttpServiceImpl

router := httpfx.NewRouter("/")
hs := httpfx.HttpService(config, router)
```

TODO(@eser) rest of the documentation
