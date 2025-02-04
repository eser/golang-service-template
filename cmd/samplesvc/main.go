package main

import (
	"context"

	"github.com/eser/go-service/pkg/samplesvc/adapters/appcontext"
	"github.com/eser/go-service/pkg/samplesvc/adapters/http"
)

func main() {
	ctx := context.Background()

	appContext, err := appcontext.NewAppContext(ctx)
	if err != nil {
		panic(err)
	}

	err = http.Run(ctx, &appContext.Config.Http, appContext.Metrics, appContext.Logger, appContext.Data)
	if err != nil {
		panic(err)
	}
}
