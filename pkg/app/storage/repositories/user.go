package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eser/go-service/pkg/app/storage/db"
	"github.com/eser/go-service/pkg/bliss/datafx"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	tx      datafx.DbExecutor
	queries *db.Queries
}

var _ datafx.Repository = (*UserRepository)(nil)

func NewUserRepository(tx datafx.DbExecutor) UserRepository {
	return UserRepository{
		tx:      tx,
		queries: db.New(tx),
	}
}

func (r UserRepository) DbExecutor() datafx.DbExecutor {
	return r.tx
}

func (r UserRepository) GetById(ctx context.Context, id string) (*db.User, error) {
	row, err := r.queries.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &row, nil
}

func (r UserRepository) GetByGithubRemoteId(ctx context.Context, githubRemoteId string) (*db.User, error) {
	row, err := r.queries.GetUserByGithubRemoteId(ctx, sql.NullString{String: githubRemoteId, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &row, nil
}

func (r UserRepository) List(ctx context.Context) ([]db.User, error) {
	rows, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r UserRepository) Create(ctx context.Context, user *db.User) (*db.User, error) {
	row, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		GithubRemoteId: user.GithubRemoteId,
		Name:           user.Name,
		Email:          user.Email,
	})
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r UserRepository) Update(ctx context.Context, user *db.User) error {
	result, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r UserRepository) SoftDelete(ctx context.Context, id string) error {
	result, err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
