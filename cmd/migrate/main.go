package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/eser/go-service/pkg/samplesvc/adapters/appcontext"
	"github.com/pressly/goose/v3"
)

var (
	ErrCommandRequired          = errors.New("command is required")
	ErrAppContextNotInitialized = errors.New("app context is not initialized")
	ErrDatabaseNotInitialized   = errors.New("database is not initialized")
	ErrDatabaseNotSqlDb         = errors.New("database is not an instance of *sql.DB")
	ErrFailedToRunGoose         = errors.New("failed to run goose")
)

func run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return ErrCommandRequired
	}

	appContext, err := appcontext.NewAppContext(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAppContextNotInitialized, err)
	}

	defaultSql := appContext.Data.GetDefaultSql()
	if defaultSql == nil {
		return ErrDatabaseNotInitialized
	}

	db, dbOk := defaultSql.GetConnection().(*sql.DB) //nolint:varnamelen
	if !dbOk {
		return ErrDatabaseNotSqlDb
	}

	command := args[0]
	rest := args[1:]

	dialect := string(defaultSql.GetDialect())

	err = goose.SetDialect(dialect)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToRunGoose, err)
	}

	err = goose.RunContext(ctx, command, db, "./etc/data/migrations", rest...)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToRunGoose, err)
	}

	return nil
}

func main() {
	ctx := context.Background()

	err := run(ctx, os.Args[1:])
	if err != nil {
		panic(err)
	}
}
