package storage

import (
	"errors"
	"fmt"

	"github.com/eser/ajan/datafx"
)

var ErrDataSourceNotFound = errors.New("data source not found")

func NewFromDefault(dataRegistry *datafx.Registry) (*Queries, error) {
	dataSource := dataRegistry.GetDefault()

	if dataSource == nil {
		return nil, fmt.Errorf("%w - default", ErrDataSourceNotFound)
	}

	db := dataSource.GetConnection()

	return &Queries{db: db}, nil
}

func NewFromNamed(dataRegistry *datafx.Registry, name string) (*Queries, error) {
	dataSource := dataRegistry.GetNamed(name)

	if dataSource == nil {
		return nil, fmt.Errorf("%w - %s", ErrDataSourceNotFound, name)
	}

	db := dataSource.GetConnection()

	return &Queries{db: db}, nil
}
