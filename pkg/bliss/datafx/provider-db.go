package datafx

import (
	"context"
	"database/sql"
)

type DataProviderDb struct {
	Connection *sql.DB
}

func (dataProviderDb *DataProviderDb) UseUnitOfWork(ctx context.Context) (UnitOfWork, error) { //nolint:ireturn
	uow, err := UseUnitOfWork(ctx, dataProviderDb.Connection)
	if err != nil {
		return nil, err
	}

	return uow, err
}
