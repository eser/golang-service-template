package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/eser/go-service/pkg/bliss"
	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/eser/go-service/pkg/bliss/datafx"
	"github.com/eser/go-service/pkg/bliss/di"
	"github.com/eser/go-service/pkg/bliss/logfx"
	"github.com/eser/go-service/pkg/bliss/metricsfx"
	"github.com/pressly/goose/v3"
)

type AppConfig bliss.BaseConfig

func LoadConfig(loader configfx.ConfigLoader) (*AppConfig, *logfx.Config, *datafx.Config, error) {
	appConfig := &AppConfig{} //nolint:exhaustruct

	err := loader.LoadDefaults(appConfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	return appConfig, &appConfig.Log, &appConfig.Data, nil
}

func main() {
	err := di.RegisterFn(
		di.Default,
		configfx.RegisterDependencies,
		LoadConfig,

		logfx.RegisterDependencies,
		metricsfx.RegisterDependencies,
		datafx.RegisterDependencies,
	)
	if err != nil {
		panic(err)
	}

	run := di.CreateInvoker(
		di.Default,
		func(
			dataProvider datafx.DataProvider,
		) error {
			allArgs := os.Args[1:]
			if len(allArgs) == 0 {
				return errors.New("command argument is required") //nolint:err113
			}

			command := allArgs[0]
			args := allArgs[1:]

			database := dataProvider.GetDefaultSql()
			if database == nil {
				return errors.New("database is not initialized") //nolint:err113
			}

			err := goose.RunWithOptionsContext(
				context.Background(),
				command,
				database.GetConnection(),
				"./ops/migrations/",
				args,
			)
			if err != nil {
				return fmt.Errorf("failed to run migrations: %w", err)
			}

			fmt.Println("migrations applied") //nolint:forbidigo

			return nil
		},
	)

	di.Seal(di.Default)

	err = run()
	if err != nil {
		panic(err)
	}
}
