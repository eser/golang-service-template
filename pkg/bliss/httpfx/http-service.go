package httpfx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
)

type HttpService interface {
	Server() *http.Server
	Router() Router

	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type HttpServiceImpl struct {
	InnerServer  *http.Server
	InnerRouter  Router
	InnerMetrics *Metrics

	Config *Config
}

var _ HttpService = (*HttpServiceImpl)(nil)

func NewHttpService(config *Config, router Router) *HttpServiceImpl {
	server := &http.Server{ //nolint:exhaustruct
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,

		Addr: config.Addr,

		Handler: router.GetMux(),
	}

	return &HttpServiceImpl{
		InnerServer: server,
		InnerRouter: router,
		Config:      config,
	}
}

func (hs *HttpServiceImpl) Server() *http.Server {
	return hs.InnerServer
}

func (hs *HttpServiceImpl) Router() Router { //nolint:ireturn
	return hs.InnerRouter
}

func (hs *HttpServiceImpl) Start(ctx context.Context) error {
	slog.InfoContext(ctx, "HttpService is starting...", slog.String("addr", hs.Config.Addr))

	// serverErr := make(chan error, 1)

	go func() {
		ln, lnErr := net.Listen("tcp", hs.InnerServer.Addr)

		if lnErr != nil {
			// serverErr <- fmt.Errorf("HttpService Net Listen error: %w", lnErr)
			os.Exit(1)

			return
		}

		if sErr := hs.Server().Serve(ln); sErr != nil && !errors.Is(sErr, http.ErrServerClosed) {
			// serverErr <- fmt.Errorf("HttpService Serve error: %w", sErr)
			os.Exit(1)

			return
		}

		// serverErr <- nil
	}() //nolint:wsl

	// if err := <-serverErr; err != nil {
	// 	return err
	// }

	return nil
}

func (hs *HttpServiceImpl) Stop(ctx context.Context) error {
	slog.InfoContext(ctx, "HttpService is stopping...")

	shutdownCtx, cancel := context.WithTimeout(ctx, hs.Config.GracefulShutdownTimeout)
	defer cancel()

	err := hs.InnerServer.Shutdown(shutdownCtx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HttpService forced to shutdown: %w", err)
	}

	<-shutdownCtx.Done()
	slog.InfoContext(ctx, "HttpService has stopped...")

	return nil
}
