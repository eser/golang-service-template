# bliss/logfx

## Overview

The **logfx** package is a configurable logging solution leverages the `log/slog` of the standard library for structured
logging. It includes pretty-printing options and a fx module for the `bliss/di` package. The package also has extensive
tests to ensure reliability and correctness, covering configuration parsing, handler behavior and the custom error
formatting logic.

The documentation below provides an overview of the package, its types, functions, and usage examples. For more detailed
information, refer to the source code and tests.

## Configuration

Configuration struct for the logger:

```go
type Config struct {
	Level      string `conf:"LEVEL"      default:"INFO"`
	PrettyMode bool   `conf:"PRETTY"     default:"true"`
	AddSource  bool   `conf:"ADD_SOURCE" default:"false"`
}
```

## Bliss DI

The `logfx` package provides a `RegisterDependencies` function that can be used to integrate with the `bliss/di`
package.

```go
import (
  ...
  "github.com/eser/go-service/pkg/bliss/di"
	"github.com/eser/go-service/pkg/bliss/logfx"
  ...
)

err := di.RegisterFn(
	logfx.RegisterDependencies,
	...
)
```

## API

### NewLogger function

Creates a new `slog.Logger` object based on the provided configuration.

```go
// func NewLogger(config *Config) (*slog.Logger, error)

logger, err := logfx.NewLogger(config)
```

### NewLoggerAsDefault function

Creates a new `slog.Logger` object based on the provided configuration and makes it default slog instance.

```go
// func NewLoggerAsDefault(config *Config) (*slog.Logger, error)

logger, err := logfx.NewLoggerAsDefault(config)
```

### Colored function

Returns a ANSI-colored string for terminal output.

```go
// func Colored(color Color, message string) string

// available colors:
//	ColorReset        ColorDimGray
//	ColorRed          ColorLightRed
//	ColorGreen        ColorLightGreen
//	ColorYellow       ColorLightYellow
//	ColorBlue         ColorLightBlue
//	ColorMagenta      ColorLightMagenta
//	ColorCyan         ColorLightCyan
//	ColorGray         ColorLightGray

fmt.Println(logfx.Colored(logfx.ColorRed, "test"))
```
