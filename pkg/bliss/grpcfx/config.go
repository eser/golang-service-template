package grpcfx

import (
	"time"
)

type Config struct {
	Addr                    string        `conf:"ADDR"             default:":9090"`
	Reflection              bool          `conf:"REFLECTION"       default:"true"`
	InitializationTimeout   time.Duration `conf:"INIT_TIMEOUT"     default:"25s"`
	GracefulShutdownTimeout time.Duration `conf:"SHUTDOWN_TIMEOUT" default:"5s"`
}
