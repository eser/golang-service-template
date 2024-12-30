package datafx

import (
	"context"
	"database/sql"
)

const DefaultDB = "DEFAULT"

type SqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type TransactionFinalizer interface {
	Rollback() error
	Commit() error
}

type DataStorer interface{}

type SqlDataStorer interface {
	GetConnection() *sql.DB
	UseUnitOfWork(ctx context.Context) (UnitOfWork, error)
}

type DataProvider interface {
	GetDefaultSql() SqlDataStorer
	GetNamedSql(name string) SqlDataStorer
}
