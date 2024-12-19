package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eser/go-service/pkg/bliss/datafx"
	"github.com/eser/go-service/pkg/eserlivesvc/data"
	"github.com/eser/go-service/pkg/eserlivesvc/shared"
)

var ErrUserNotFound = errors.New("user not found")

type UserService struct {
	scope   datafx.DbExecutor
	queries *data.Queries
}

func NewUserService(scope datafx.DbExecutor) UserService {
	return UserService{
		scope:   scope,
		queries: data.New(scope),
	}
}

// func (r UserService) DbScope() datafx.DbExecutor { //nolint:ireturn
// 	return r.scope
// }

func (r UserService) GetById(ctx context.Context, id string) (*data.User, error) {
	row, err := r.queries.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err //nolint:wrapcheck
	}

	return &row, nil
}

func (r UserService) GetByGithubRemoteId(ctx context.Context, githubRemoteId string) (*data.User, error) {
	row, err := r.queries.GetUserByGithubRemoteId(ctx, sql.NullString{String: githubRemoteId, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err //nolint:wrapcheck
	}

	return &row, nil
}

func (r UserService) List(ctx context.Context, cursor shared.Cursor) (*shared.RecordsWithCursor[data.User], error) {
	rows, err := r.queries.ListUsers(ctx, data.ListUsersParams{
		Offset: cursor.Offset,
		Limit:  cursor.Limit,
	})
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	result := shared.NewRecordsWithCursor(rows, cursor.Limit, func(user *data.User) string {
		return user.Id
	})

	return result, nil
}

func (r UserService) Create(ctx context.Context, user *data.User) (*data.User, error) {
	row, err := r.queries.CreateUser(ctx, data.CreateUserParams{
		Id:             user.Id,
		Name:           user.Name,
		Email:          user.Email,
		Phone:          user.Phone,
		GithubRemoteId: user.GithubRemoteId,
		GithubHandle:   user.GithubHandle,
		XRemoteId:      user.XRemoteId,
		XHandle:        user.XHandle,
	})
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return &row, nil
}

func (r UserService) Update(ctx context.Context, user *data.User) error {
	result, err := r.queries.UpdateUser(ctx, data.UpdateUserParams{
		Id:   user.Id,
		Name: user.Name,
	})
	if err != nil {
		return err //nolint:wrapcheck
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err //nolint:wrapcheck
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r UserService) SoftDelete(ctx context.Context, id string) error {
	result, err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return err //nolint:wrapcheck
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err //nolint:wrapcheck
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
