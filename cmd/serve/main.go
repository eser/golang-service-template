package main

import (
	"context"
	"log/slog"

	"github.com/eser/go-service/pkg/sample/adapters/appcontext"
	"github.com/eser/go-service/pkg/sample/adapters/http"
)

func main() {
	ctx := context.Background()

	appContext, err := appcontext.NewAppContext(ctx)
	if err != nil {
		panic(err)
	}

	appContext.Logger.InfoContext(
		ctx,
		"Starting service",
		slog.String("name", appContext.Config.AppName),
		slog.String("environment", appContext.Config.AppEnv),
		slog.Any("features", appContext.Config.Features),
	)

	err = http.Run(ctx, &appContext.Config.Http, appContext.Metrics, appContext.Logger, appContext.Data)
	if err != nil {
		panic(err)
	}
}
