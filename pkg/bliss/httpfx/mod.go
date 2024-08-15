package httpfx

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"go.uber.org/fx"
)

//nolint:gochecknoglobals
var Module = fx.Module(
	"httpservice",
	fx.Provide(
		New,
	),
	fx.Invoke(
		RegisterHooks,
	),
)

type Result struct {
	fx.Out

	HttpService *HttpService
	Routes      *Router
}

type HttpService struct {
	Server *http.Server

	Config *Config

	Routes *Router
}

func New() (Result, error) {
	routes := NewRouter("/")

	config, err := NewConfig()
	if err != nil {
		return Result{}, fmt.Errorf("error creating new config: %w", err)
	}

	server := &http.Server{
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,

		Addr: config.Addr,

		Handler: routes.Mux,
	}

	return Result{
		HttpService: &HttpService{server, config, routes},
		Routes:      routes,
	}, nil
}

// , conf *config.Config, logger *log.Logger.
func RegisterHooks(lc fx.Lifecycle, hs *HttpService) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// logger.Info("HttpService is starting...", log.String("env", conf.Env), log.String("port", conf.Port))

			// serverErr := make(chan error, 1)

			go func() {
				ln, lnErr := net.Listen("tcp", hs.Server.Addr) //nolint:varnamelen

				if lnErr != nil {
					// serverErr <- fmt.Errorf("HttpService Net Listen error: %w", lnErr)
					os.Exit(1)

					return
				}

				if sErr := hs.Server.Serve(ln); sErr != nil && !errors.Is(sErr, http.ErrServerClosed) {
					// serverErr <- fmt.Errorf("HttpService Serve error: %w", sErr)
					os.Exit(1)

					return
				}

				// serverErr <- nil
			}()

			// if err := <-serverErr; err != nil {
			// 	return err
			// }

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// logger.Info("HttpService is stopping...")

			shutdownCtx, cancel := context.WithTimeout(ctx, hs.Config.GracefulShutdownTimeout)
			defer cancel()

			err := hs.Server.Shutdown(shutdownCtx)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				return fmt.Errorf("HttpService forced to shutdown: %w", err)
			}

			<-shutdownCtx.Done()
			// logger.Info("HttpService has stopped...")

			return nil
		},
	})
}
