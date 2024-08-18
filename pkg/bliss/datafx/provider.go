package datafx

type DataProvider interface {
	// CreateUnitOfWork() UnitOfWork
}

type DataProviderImpl struct{}

func NewDataProvider() DataProvider { //nolint:ireturn
	return &DataProviderImpl{}
}
