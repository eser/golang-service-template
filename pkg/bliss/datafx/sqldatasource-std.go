package datafx

import (
	"context"
	"database/sql"
	"fmt"
)

type SqlDataSourceStd struct {
	connection *sql.DB
}

var _ SqlDataSource = (*SqlDataSourceStd)(nil)

func NewSqlDataSourceStd(ctx context.Context, dialect Dialect, dsn string) (*SqlDataSourceStd, error) {
	connection, err := sql.Open(string(dialect), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := connection.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &SqlDataSourceStd{
		connection: connection,
	}, nil
}

func (dataSource *SqlDataSourceStd) GetConnection() SqlExecutor { //nolint:ireturn
	return dataSource.connection
}

func (dataSource *SqlDataSourceStd) UseUnitOfWork(ctx context.Context) (UnitOfWork, error) { //nolint:ireturn
	uow, err := UseUnitOfWork(ctx, dataSource.connection)
	if err != nil {
		return nil, err
	}

	return uow, err
}
