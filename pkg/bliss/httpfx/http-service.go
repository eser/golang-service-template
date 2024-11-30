package httpfx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/eser/go-service/pkg/bliss/metricsfx"
)

type HttpService interface {
	Server() *http.Server
	Router() Router

	Start(ctx context.Context) (func(), error)
}

type HttpServiceImpl struct {
	InnerServer  *http.Server
	InnerRouter  Router
	InnerMetrics *Metrics

	Config *Config
	logger *slog.Logger
}

var _ HttpService = (*HttpServiceImpl)(nil)

func NewHttpService(
	config *Config,
	router Router,
	metricsProvider metricsfx.MetricsProvider,
	logger *slog.Logger,
) *HttpServiceImpl {
	server := &http.Server{ //nolint:exhaustruct
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,

		Addr: config.Addr,

		Handler: router.GetMux(),
	}

	return &HttpServiceImpl{
		InnerServer:  server,
		InnerRouter:  router,
		InnerMetrics: NewMetrics(metricsProvider),
		Config:       config,
		logger:       logger,
	}
}

func (hs *HttpServiceImpl) Server() *http.Server {
	return hs.InnerServer
}

func (hs *HttpServiceImpl) Router() Router { //nolint:ireturn
	return hs.InnerRouter
}

func (hs *HttpServiceImpl) Start(ctx context.Context) (func(), error) {
	hs.logger.InfoContext(ctx, "HttpService is starting...", slog.String("addr", hs.Config.Addr))

	listener, lnErr := net.Listen("tcp", hs.InnerServer.Addr)
	if lnErr != nil {
		return nil, fmt.Errorf("HttpService Net Listen error: %w", lnErr)
	}

	go func() {
		if sErr := hs.InnerServer.Serve(listener); sErr != nil && !errors.Is(sErr, http.ErrServerClosed) {
			hs.logger.ErrorContext(ctx, "HttpService Serve error: %w", slog.Any("error", sErr))
		}
	}()

	cleanup := func() {
		hs.logger.InfoContext(ctx, "Shutting down server...")

		newCtx, cancel := context.WithTimeout(ctx, hs.Config.GracefulShutdownTimeout)
		defer cancel()

		if err := hs.InnerServer.Shutdown(newCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			hs.logger.ErrorContext(ctx, "HttpService forced to shutdown", slog.Any("error", err))

			return
		}

		hs.logger.InfoContext(ctx, "HttpService has gracefully stopped.")
	}

	return cleanup, nil
}
