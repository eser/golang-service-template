package httpfx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"go.uber.org/fx"
)

var FxModule = fx.Module( //nolint:gochecknoglobals
	"httpservice",
	fx.Provide(
		New,
	),
	fx.Invoke(
		RegisterHooks,
	),
)

type FxResult struct {
	fx.Out

	HttpService *HttpService
	Routes      Router
}

type HttpService struct {
	Server *http.Server

	Config *Config

	Routes Router
}

func New(config *Config) (FxResult, error) {
	routes := NewRouter("/")

	server := &http.Server{ //nolint:exhaustruct
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,

		Addr: config.Addr,

		Handler: routes.GetMux(),
	}

	return FxResult{ //nolint:exhaustruct
		HttpService: &HttpService{Server: server, Config: config, Routes: routes},
		Routes:      routes,
	}, nil
}

// , conf *config.Config, logger *log.Logger.
func RegisterHooks(lc fx.Lifecycle, hs *HttpService, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("HttpService is starting...", slog.String("addr", hs.Config.Addr))

			// serverErr := make(chan error, 1)

			go func() {
				ln, lnErr := net.Listen("tcp", hs.Server.Addr)

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
			logger.Info("HttpService is stopping...")

			shutdownCtx, cancel := context.WithTimeout(ctx, hs.Config.GracefulShutdownTimeout)
			defer cancel()

			err := hs.Server.Shutdown(shutdownCtx)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				return fmt.Errorf("HttpService forced to shutdown: %w", err)
			}

			<-shutdownCtx.Done()
			logger.Info("HttpService has stopped...")

			return nil
		},
	})
}
