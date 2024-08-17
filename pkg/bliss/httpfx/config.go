package httpfx

import (
	"time"
)

type Config struct {
	ReadHeaderTimeout time.Duration `conf:"READ_HEADER_TIMEOUT" default:"5s"`
	ReadTimeout       time.Duration `conf:"READ_TIMEOUT"        default:"10s"`
	WriteTimeout      time.Duration `conf:"WRITE_TIMEOUT"       default:"10s"`
	IdleTimeout       time.Duration `conf:"IDLE_TIMEOUT"        default:"120s"`

	InitializationTimeout   time.Duration `conf:"INIT_TIMEOUT"     default:"25s"`
	GracefulShutdownTimeout time.Duration `conf:"SHUTDOWN_TIMEOUT" default:"5s"`

	Addr string `conf:"ADDR" default:":8080"`
}
