package datafx

import "database/sql"

type DataProvider interface {
	CreateUnitOfWork() (UnitOfWork, error)
	GetDB() *sql.DB
}

type DataProviderImpl struct {
	db *sql.DB
}

var _ DataProvider = (*DataProviderImpl)(nil)

func NewDataProvider(db *sql.DB) *DataProviderImpl {
	return &DataProviderImpl{db: db}
}

func (dp *DataProviderImpl) GetDB() *sql.DB {
	return dp.db
}

func (dp *DataProviderImpl) CreateUnitOfWork() (UnitOfWork, error) {
	tx, err := dp.db.Begin()
	if err != nil {
		return nil, err
	}

	return NewUnitOfWork(tx), nil
}
