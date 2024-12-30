package users

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/eser/go-service/pkg/bliss/datafx"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/eserlivesvc/shared"
)

func RegisterHttpRoutes(routes httpfx.Router, appConfig *shared.AppConfig, logger *slog.Logger, dataProvider datafx.DataProvider) { //nolint:lll
	routes.
		Route("GET /users", func(ctx *httpfx.Context) httpfx.Result {
			dataStorage := dataProvider.GetDefaultSql()

			cursor := shared.Cursor{
				Offset: ctx.Request.URL.Query().Get("cursor"),
				Limit:  10, //nolint:mnd
			}

			users, err := NewUserService(dataStorage).List(ctx.Request.Context(), cursor)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			return ctx.Results.Json(users)
		}).
		HasSummary("List users").
		HasDescription("List users.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /users/{id}", func(ctx *httpfx.Context) httpfx.Result {
			id := ctx.Request.PathValue("id")
			dataStorage := dataProvider.GetDefaultSql()

			user, err := NewUserService(dataStorage).GetById(ctx.Request.Context(), id)
			if err != nil {
				if errors.Is(err, ErrUserNotFound) {
					return ctx.Results.Error(http.StatusNotFound, []byte("user not found"))
				}

				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			return ctx.Results.Json(user)
		}).
		HasSummary("Get user by ID").
		HasDescription("Get a single user by their ID.").
		HasResponse(http.StatusOK).
		HasResponse(http.StatusNotFound)
}
