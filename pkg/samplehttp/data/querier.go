// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package data

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateChannel(ctx context.Context, arg CreateChannelParams) (Channel, error)
	DeleteChannel(ctx context.Context, id string) (sql.Result, error)
	GetChannelById(ctx context.Context, id string) (Channel, error)
	GetChannelByName(ctx context.Context, name sql.NullString) (Channel, error)
	ListChannels(ctx context.Context) ([]Channel, error)
	UpdateChannel(ctx context.Context, arg UpdateChannelParams) (sql.Result, error)
}

var _ Querier = (*Queries)(nil)
