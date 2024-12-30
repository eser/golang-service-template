package datafx

import (
	"context"
	"database/sql"
)

type DataManagerDb struct {
	connection *sql.DB
}

var _ SqlDataStorer = (*DataManagerDb)(nil)

func NewDataManagerDb(connection *sql.DB) *DataManagerDb {
	return &DataManagerDb{
		connection: connection,
	}
}

func (dataManagerDb *DataManagerDb) GetConnection() *sql.DB {
	return dataManagerDb.connection
}

func (dataManagerDb *DataManagerDb) UseUnitOfWork(ctx context.Context) (UnitOfWork, error) { //nolint:ireturn
	uow, err := UseUnitOfWork(ctx, dataManagerDb.connection)
	if err != nil {
		return nil, err
	}

	return uow, err
}
