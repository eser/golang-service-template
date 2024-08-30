package datafx

import (
	"context"
	"database/sql"
)

type DbExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	// Exec(query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	// Prepare(query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	// Query(query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	// QueryRow(query string, args ...any) *sql.Row
}

type Repository interface {
	DbExecutor() DbExecutor
}
