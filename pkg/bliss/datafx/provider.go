package datafx

type DataProvider interface {
	// CreateUnitOfWork() UnitOfWork
}

type DataProviderImpl struct{}

var _ DataProvider = (*DataProviderImpl)(nil)

func NewDataProvider() *DataProviderImpl {
	return &DataProviderImpl{}
}
