package datafx

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type DataManager struct {
	dbs    map[string]*DataManagerDb
	logger *slog.Logger
}

var _ DataProvider = (*DataManager)(nil)

func NewDataManager(logger *slog.Logger) *DataManager {
	dbs := make(map[string]*DataManagerDb)

	return &DataManager{
		dbs:    dbs,
		logger: logger,
	}
}

func (dataProvider *DataManager) GetDefaultSql() SqlDataStorer { //nolint:ireturn
	return dataProvider.dbs[DefaultDB]
}

func (dataProvider *DataManager) GetNamedSql(name string) SqlDataStorer { //nolint:ireturn
	if db, exists := dataProvider.dbs[name]; exists {
		return db
	}

	return nil
}

func (dataProvider *DataManager) AddConnection(name string, dsn string) error {
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

	dataProvider.dbs[name] = NewDataManagerDb(database)

	dataProvider.logger.Info("successfully added database connection", slog.String("name", name))

	return nil
}

func (dataProvider *DataManager) LoadFromConfig(config *Config) error {
	for name, source := range config.Sources {
		err := dataProvider.AddConnection(name, source.DSN)
		if err != nil {
			return fmt.Errorf("failed to add connection for %s: %w", name, err)
		}
	}

	return nil
}
