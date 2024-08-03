package httpservice

import (
	"fmt"
	"time"
)

const (
	DefaultReadHeaderTimeout = "5s"
	DefaultReadTimeout       = "10s"
	DefaultWriteTimeout      = "10s"
	DefaultIdleTimeout       = "120s"

	DefaultInitializationTimeout   = "25s"
	DefaultGracefulShutdownTimeout = "5s"

	DefaultAddr = ":8000"
)

type Config struct {
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration

	InitializationTimeout   time.Duration
	GracefulShutdownTimeout time.Duration

	Addr string
}

func NewConfig() (Config, error) {
	readHeaderTimeout, err := time.ParseDuration(DefaultReadHeaderTimeout)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing read header timeout: %w", err)
	}

	readTimeout, err := time.ParseDuration(DefaultReadTimeout)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing read timeout: %w", err)
	}

	writeTimeout, err := time.ParseDuration(DefaultWriteTimeout)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing write timeout: %w", err)
	}

	idleTimeout, err := time.ParseDuration(DefaultIdleTimeout)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing idle timeout: %w", err)
	}

	initializationTimeout, err := time.ParseDuration(DefaultInitializationTimeout)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing initialization timeout: %w", err)
	}

	gracefulShutdownTimeout, err := time.ParseDuration(DefaultGracefulShutdownTimeout)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing graceful shutdown timeout: %w", err)
	}

	return Config{
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,

		InitializationTimeout:   initializationTimeout,
		GracefulShutdownTimeout: gracefulShutdownTimeout,

		Addr: DefaultAddr,
	}, nil
}
