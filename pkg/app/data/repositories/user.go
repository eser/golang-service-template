package repositories

import (
	"context"

	"github.com/eser/go-service/pkg/app/data/models"
	"github.com/eser/go-service/pkg/bliss/datafx"
)

type UserRepository interface {
	datafx.Repository

	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
}

type UserRepositoryImpl struct {
	datafx.RepositoryBase
}

var _ UserRepository = (*UserRepositoryImpl)(nil)

func NewUserRepository(db datafx.DbExecutor) UserRepositoryImpl {
	return UserRepositoryImpl{
		RepositoryBase: datafx.NewRepositoryBase(db),
	}
}

func (r UserRepositoryImpl) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, nil //nolint:nilnil
}

func (r UserRepositoryImpl) GetUserById(ctx context.Context, id string) (*models.User, error) {
	return nil, nil //nolint:nilnil
}
