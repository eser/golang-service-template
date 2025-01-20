package httpfx

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type HttpService struct {
	InnerServer  *http.Server
	InnerRouter  *Router
	InnerMetrics *Metrics

	Config *Config
	logger *slog.Logger
}

type MetricsProvider interface {
	GetRegistry() *prometheus.Registry
}

func NewHttpService(
	config *Config,
	router *Router,
	metricsProvider MetricsProvider,
	logger *slog.Logger,
) *HttpService {
	server := &http.Server{ //nolint:exhaustruct
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,

		Addr: config.Addr,

		Handler: router.GetMux(),
	}

	if config.CertString != "" && config.KeyString != "" {
		cert, err := tls.X509KeyPair([]byte(config.CertString), []byte(config.KeyString))
		if err != nil {
			panic(fmt.Errorf("failed to load certificate: %w", err))
		}

		server.TLSConfig = &tls.Config{ //nolint:exhaustruct
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
	}

	return &HttpService{
		InnerServer:  server,
		InnerRouter:  router,
		InnerMetrics: NewMetrics(metricsProvider),
		Config:       config,
		logger:       logger,
	}
}

func (hs *HttpService) Server() *http.Server {
	return hs.InnerServer
}

func (hs *HttpService) Router() *Router {
	return hs.InnerRouter
}

func (hs *HttpService) Start(ctx context.Context) (func(), error) {
	hs.logger.InfoContext(ctx, "HttpService is starting...", slog.String("addr", hs.Config.Addr))

	listener, lnErr := net.Listen("tcp", hs.InnerServer.Addr)
	if lnErr != nil {
		return nil, fmt.Errorf("HttpService Net Listen error: %w", lnErr)
	}

	go func() {
		var sErr error

		if hs.Server().TLSConfig != nil {
			sErr = hs.InnerServer.ServeTLS(listener, "", "")
		} else {
			sErr = hs.InnerServer.Serve(listener)
		}

		if sErr != nil && !errors.Is(sErr, http.ErrServerClosed) {
			hs.logger.ErrorContext(ctx, "HttpService ServeTLS error: %w", slog.Any("error", sErr))
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
