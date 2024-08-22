package datafx

type DataProvider interface {
	// CreateUnitOfWork() UnitOfWork
}

type DataProviderImpl struct{}

func NewDataProvider() *DataProviderImpl {
	return &DataProviderImpl{}
}
