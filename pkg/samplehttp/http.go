package samplehttp

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/eser/go-service/pkg/bliss/datafx"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/proto-go/broadcast"
)

func RegisterHttpRoutes(routes httpfx.Router, appConfig *AppConfig, logger *slog.Logger, dataProvider datafx.DataProvider) { //nolint:lll
	routes.
		Route("GET /channels", func(ctx *httpfx.Context) httpfx.Result {
			dataStorage := dataProvider.GetDefaultSql()

			channels, err := NewChannelService(dataStorage).List(ctx.Request.Context())
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

			var payload broadcast.SendRequest
			err = json.Unmarshal(body, &payload)
			if err != nil {
				return ctx.Results.Error(http.StatusBadRequest, []byte(err.Error()))
			}

			logger.Info(
				"Send",
				slog.String("channelId", payload.GetChannelId()),
				slog.Any("message", payload.GetMessage()),
			)

			return ctx.Results.Ok()
		}).
		HasSummary("Send a message to a channel").
		HasDescription("Send a message to a channel.").
		HasResponse(http.StatusOK)
}
