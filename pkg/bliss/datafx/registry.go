package datafx

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

const DefaultDB = "DEFAULT"

type SqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type SqlDataSource interface {
	GetConnection() SqlExecutor
	UseUnitOfWork(ctx context.Context) (*UnitOfWork, error)
}

type Registry struct {
	dbs    map[string]SqlDataSource
	logger *slog.Logger
}

func NewRegistry(logger *slog.Logger) *Registry {
	dbs := make(map[string]SqlDataSource)

	return &Registry{
		dbs:    dbs,
		logger: logger,
	}
}

func (registry *Registry) GetDefaultSql() SqlDataSource { //nolint:ireturn
	return registry.dbs[DefaultDB]
}

func (registry *Registry) GetNamedSql(name string) SqlDataSource { //nolint:ireturn
	if db, exists := registry.dbs[name]; exists {
		return db
	}

	return nil
}

func (registry *Registry) AddConnection(ctx context.Context, name string, dsn string) error {
	registry.logger.Info(
		"adding database connection",
		slog.String("name", name),
		slog.String("dialect", string(DetermineDialect(dsn))),
	)

	dialect := DetermineDialect(dsn)

	// var db SqlDataSource

	// var err error

	// if dialect == DialectPostgresPgx {
	// 	db, err = NewSqlDataSourcePgx(ctx, dialect, dsn)
	// } else {
	db, err := NewSqlDataSourceStd(ctx, dialect, dsn) //nolint:varnamelen
	// }
	if err != nil {
		registry.logger.Error(
			"failed to open database connection",
			slog.String("error", err.Error()),
			slog.String("name", name),
		)

		return fmt.Errorf("failed to add connection for %s: %w", name, err)
	}

	registry.dbs[name] = db

	registry.logger.Info("successfully added database connection", slog.String("name", name))

	return nil
}

func (registry *Registry) LoadFromConfig(ctx context.Context, config *Config) error {
	for name, source := range config.Sources {
		err := registry.AddConnection(ctx, name, source.DSN)
		if err != nil {
			return fmt.Errorf("failed to add connection for %s: %w", name, err)
		}
	}

	return nil
}
