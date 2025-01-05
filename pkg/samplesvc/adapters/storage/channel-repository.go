package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eser/go-service/pkg/bliss/datafx"
	"github.com/eser/go-service/pkg/samplesvc/adapters/storage/data"
	"github.com/eser/go-service/pkg/samplesvc/business/channel"
)

var ErrChannelNotFound = errors.New("channel not found")

type ChannelRepository struct {
	scope   datafx.SqlDataSource
	queries *data.Queries
}

func NewChannelRepository(dataSource datafx.SqlDataSource) ChannelRepository {
	return ChannelRepository{
		scope:   dataSource,
		queries: data.New(dataSource.GetConnection()),
	}
}

func (r ChannelRepository) GetById(ctx context.Context, id string) (*channel.Channel, error) {
	row, err := r.queries.GetChannelById(ctx, id)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	if row == nil {
		return nil, ErrChannelNotFound
	}

	return row, nil
}

func (r ChannelRepository) GetByName(ctx context.Context, name string) (*channel.Channel, error) {
	row, err := r.queries.GetChannelByName(ctx, sql.NullString{String: name, Valid: true})
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	if row == nil {
		return nil, ErrChannelNotFound
	}

	return row, nil
}

func (r ChannelRepository) List(ctx context.Context) ([]*channel.Channel, error) {
	rows, err := r.queries.ListChannels(ctx)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return rows, nil
}

func (r ChannelRepository) Create(ctx context.Context, channel *channel.Channel) (*channel.Channel, error) {
	row, err := r.queries.CreateChannel(ctx, data.CreateChannelParams{
		Id:   channel.Id,
		Name: channel.Name,
	})
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return row, nil
}

func (r ChannelRepository) Update(ctx context.Context, channel *channel.Channel) error {
	rowsAffected, err := r.queries.UpdateChannel(ctx, data.UpdateChannelParams{
		Id:   channel.Id,
		Name: channel.Name,
	})
	if err != nil {
		return err //nolint:wrapcheck
	}

	if rowsAffected == 0 {
		return ErrChannelNotFound
	}

	return nil
}

func (r ChannelRepository) SoftDelete(ctx context.Context, id string) error {
	rowsAffected, err := r.queries.DeleteChannel(ctx, id)
	if err != nil {
		return err //nolint:wrapcheck
	}

	if rowsAffected == 0 {
		return ErrChannelNotFound
	}

	return nil
}
