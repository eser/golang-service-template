package datafx

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

const DefaultDB = "DEFAULT"

type DataProvider interface {
	GetDefault() *DataProviderDb
	GetNamed(name string) *DataProviderDb
}

type DataProviderImpl struct {
	dbs    map[string]*DataProviderDb
	logger *slog.Logger
}

var _ DataProvider = (*DataProviderImpl)(nil)

func NewDataProvider(logger *slog.Logger) *DataProviderImpl {
	dbs := make(map[string]*DataProviderDb)

	return &DataProviderImpl{
		dbs:    dbs,
		logger: logger,
	}
}

func (dataProvider *DataProviderImpl) GetDefault() *DataProviderDb {
	return dataProvider.dbs[DefaultDB]
}

func (dataProvider *DataProviderImpl) GetNamed(name string) *DataProviderDb {
	if db, exists := dataProvider.dbs[name]; exists {
		return db
	}

	return nil
}

func (dataProvider *DataProviderImpl) AddConnection(name string, dsn string) error {
	dataProvider.logger.Info(
		"adding database connection",
		slog.String("name", name),
		slog.String("dialect", string(DetermineDialect(dsn))),
	)

	dialect := DetermineDialect(dsn)

	database, err := sql.Open(string(dialect), dsn)
	if err != nil {
		dataProvider.logger.Error(
			"failed to open database connection",
			slog.String("error", err.Error()),
			slog.String("name", name),
		)

		return fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := database.PingContext(context.TODO()); err != nil {
		dataProvider.logger.Error("failed to ping database", slog.String("error", err.Error()), slog.String("name", name))

		return fmt.Errorf("failed to ping database: %w", err)
	}

	dataProvider.dbs[name] = &DataProviderDb{Connection: database}
	dataProvider.logger.Info("successfully added database connection", slog.String("name", name))

	return nil
}

func (dataProvider *DataProviderImpl) LoadFromConfig(config *Config) error {
	for name, source := range config.Sources {
		err := dataProvider.AddConnection(name, source.DSN)
		if err != nil {
			return fmt.Errorf("failed to add connection for %s: %w", name, err)
		}
	}

	return nil
}
