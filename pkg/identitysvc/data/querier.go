// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package data

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteUser(ctx context.Context, id string) (sql.Result, error)
	GetUserByGithubRemoteId(ctx context.Context, githubRemoteID sql.NullString) (User, error)
	GetUserById(ctx context.Context, id string) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error)
}

var _ Querier = (*Queries)(nil)