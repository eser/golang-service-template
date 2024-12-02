package samplesvc

import (
	"context"
	"fmt"

	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/eser/go-service/pkg/bliss/di"
	"github.com/eser/go-service/pkg/bliss/grpcfx"
	"github.com/eser/go-service/pkg/bliss/lib"
	"github.com/eser/go-service/pkg/bliss/logfx"
	"github.com/eser/go-service/pkg/bliss/metricsfx"
)

func LoadConfig(loader configfx.ConfigLoader) (*AppConfig, *logfx.Config, *grpcfx.Config, error) {
	appConfig := &AppConfig{} //nolint:exhaustruct

	err := loader.LoadDefaults(appConfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	return appConfig, &appConfig.Log, &appConfig.Grpc, nil
}

func Run() error {
	err := di.RegisterFn(
		di.Default,
		configfx.RegisterDependencies,
		LoadConfig,

		logfx.RegisterDependencies,
		metricsfx.RegisterDependencies,
		grpcfx.RegisterDependencies,

		RegisterGrpcService,
	)
	if err != nil {
		panic(err)
	}

	run := di.CreateInvoker(
		di.Default,
		func(
			grpcService grpcfx.GrpcService,
		) error {
			ctx := context.Background()

			cleanup, err := grpcService.Start(ctx)
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
