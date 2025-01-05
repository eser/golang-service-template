package samplesvc

import (
	"context"

	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/eser/go-service/pkg/bliss/datafx"
	"github.com/eser/go-service/pkg/bliss/di"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/httpfx/modules/healthcheck"
	"github.com/eser/go-service/pkg/bliss/httpfx/modules/openapi"
	"github.com/eser/go-service/pkg/bliss/httpfx/modules/profiling"
	"github.com/eser/go-service/pkg/bliss/lib"
	"github.com/eser/go-service/pkg/bliss/logfx"
	"github.com/eser/go-service/pkg/bliss/metricsfx"
	"github.com/eser/go-service/pkg/samplesvc/adapters/http"
)

func RunWithDi() error {
	err := di.RegisterFn(
		di.Default,
		configfx.RegisterDependencies,
		LoadConfig,

		logfx.RegisterDependencies,
		metricsfx.RegisterDependencies,
		httpfx.RegisterDependencies,
		datafx.RegisterDependencies,

		RegisterHttpMiddlewares,

		healthcheck.RegisterHttpRoutes,
		openapi.RegisterHttpRoutes,
		profiling.RegisterHttpRoutes,

		http.RegisterHttpRoutes,
	)
	if err != nil {
		panic(err)
	}

	run := di.CreateInvoker(
		di.Default,
		func(
			httpService httpfx.HttpService,
		) error {
			ctx := context.Background()

			cleanup, err := httpService.Start(ctx)
			if err != nil {
				return err //nolint:wrapcheck
			}

			lib.WaitForSignal()

			cleanup()

			return nil
		},
	)

	di.Seal(di.Default)

	err = run()

	return err
}
