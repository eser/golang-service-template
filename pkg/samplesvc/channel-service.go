package samplesvc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eser/go-service/pkg/bliss/datafx"
	"github.com/eser/go-service/pkg/samplesvc/data"
)

var ErrChannelNotFound = errors.New("channel not found")

type ChannelService struct {
	scope   datafx.DbExecutor
	queries *data.Queries
}

func NewChannelService(dataStorage datafx.DataStorer) ChannelService {
	return ChannelService{
		scope:   scope,
		queries: data.New(scope),
	}
}

// func (r ChannelService) DbScope() datafx.DbExecutor { //nolint:ireturn
// 	return r.scope
// }

func (r ChannelService) GetById(ctx context.Context, id string) (*data.Channel, error) {
	row, err := r.queries.GetChannelById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrChannelNotFound
		}

		return nil, err //nolint:wrapcheck
	}

	return &row, nil
}

func (r ChannelService) GetByName(ctx context.Context, name string) (*data.Channel, error) {
	row, err := r.queries.GetChannelByName(ctx, sql.NullString{String: name, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrChannelNotFound
		}

		return nil, err //nolint:wrapcheck
	}

	return &row, nil
}

func (r ChannelService) List(ctx context.Context) ([]data.Channel, error) {
	rows, err := r.queries.ListChannels(ctx)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return rows, nil
}

func (r ChannelService) Create(ctx context.Context, channel *data.Channel) (*data.Channel, error) {
	row, err := r.queries.CreateChannel(ctx, data.CreateChannelParams{
		Id:   channel.Id,
		Name: channel.Name,
	})
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return &row, nil
}

func (r ChannelService) Update(ctx context.Context, channel *data.Channel) error {
	result, err := r.queries.UpdateChannel(ctx, data.UpdateChannelParams{
		Id:   channel.Id,
		Name: channel.Name,
	})
	if err != nil {
		return err //nolint:wrapcheck
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err //nolint:wrapcheck
	}

	if rowsAffected == 0 {
		return ErrChannelNotFound
	}

	return nil
}

func (r ChannelService) SoftDelete(ctx context.Context, id string) error {
	result, err := r.queries.DeleteChannel(ctx, id)
	if err != nil {
		return err //nolint:wrapcheck
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err //nolint:wrapcheck
	}

	if rowsAffected == 0 {
		return ErrChannelNotFound
	}

	return nil
}
