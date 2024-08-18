package repositories

import (
	"context"

	"github.com/eser/go-service/pkg/app/data/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
}

type UserRepositoryImpl struct{}

func NewUserRepository() UserRepository { //nolint:ireturn
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, nil //nolint:nilnil
}

func (r *UserRepositoryImpl) GetUserById(ctx context.Context, id string) (*models.User, error) {
	return nil, nil //nolint:nilnil
}
