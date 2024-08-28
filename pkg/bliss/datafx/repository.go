package datafx

import (
	"context"
	"database/sql"
)

type DbExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryRow(query string, args ...any) *sql.Row
}

type Repository interface {
	DbExecutor() DbExecutor
}

type RepositoryBase struct {
	dbExecutor DbExecutor
}

var _ Repository = (*RepositoryBase)(nil)

func NewRepositoryBase(db DbExecutor) RepositoryBase {
	return RepositoryBase{
		dbExecutor: db,
	}
}

func (r *RepositoryBase) DbExecutor() DbExecutor {
	return r.dbExecutor
}
