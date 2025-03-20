package http

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/eser/ajan/datafx"
	"github.com/eser/ajan/httpfx"
	"github.com/eser/ajan/httpfx/middlewares"
	"github.com/eser/ajan/httpfx/modules/healthcheck"
	"github.com/eser/ajan/httpfx/modules/openapi"
	"github.com/eser/ajan/httpfx/modules/profiling"
	"github.com/eser/ajan/lib"
	"github.com/eser/ajan/logfx"
	"github.com/eser/ajan/metricsfx"
	"github.com/eser/go-service/pkg/sample/adapters/storage"
	"github.com/eser/go-service/pkg/sample/business/channels"
)

func RegisterHttpRoutes(routes *httpfx.Router, logger *logfx.Logger, dataRegistry *datafx.Registry) {
	routes.
		Route("GET /channels", func(ctx *httpfx.Context) httpfx.Result {
			queries, err := storage.NewFromDefault(dataRegistry)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			service := channels.NewService(queries)

			channels, err := service.List(ctx.Request.Context())
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			return ctx.Results.Json(channels)
		}).
		HasSummary("List channels").
		HasDescription("List channels.").
		HasResponse(http.StatusOK)

	routes.
		Route("POST /send", func(ctx *httpfx.Context) httpfx.Result {
			body, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			var payload channels.Channel
			err = json.Unmarshal(body, &payload)
			if err != nil {
				return ctx.Results.Error(http.StatusBadRequest, []byte(err.Error()))
			}

			logger.Info(
				"Send",
				slog.String("id", payload.Id),
				slog.String("name", payload.Name.String),
			)

			return ctx.Results.Ok()
		}).
		HasSummary("Send a message to a channel").
		HasDescription("Send a message to a channel.").
		HasResponse(http.StatusOK)
}

func Run(ctx context.Context, config *httpfx.Config, metricsProvider *metricsfx.MetricsProvider, logger *logfx.Logger, dataRegistry *datafx.Registry) error { //nolint:lll
	routes := httpfx.NewRouter("/")
	httpService := httpfx.NewHttpService(config, routes, metricsProvider, logger)

	// http middlewares
	routes.Use(middlewares.ErrorHandlerMiddleware())
	routes.Use(middlewares.ResolveAddressMiddleware())
	routes.Use(middlewares.ResponseTimeMiddleware())
	routes.Use(middlewares.CorrelationIdMiddleware())
	routes.Use(middlewares.CorsMiddleware())
	routes.Use(middlewares.MetricsMiddleware(httpService.InnerMetrics))

	// http modules
	healthcheck.RegisterHttpRoutes(routes, config)
	openapi.RegisterHttpRoutes(routes, config)
	profiling.RegisterHttpRoutes(routes, config)

	// http routes
	RegisterHttpRoutes(routes, logger, dataRegistry) //nolint:contextcheck

	// run
	cleanup, err := httpService.Start(ctx)
	if err != nil {
		return err //nolint:wrapcheck
	}

	defer cleanup()

	lib.WaitForSignal()

	return nil
}
