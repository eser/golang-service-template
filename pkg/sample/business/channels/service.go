package channels

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	GetChannelById(ctx context.Context, id string) (*Channel, error)
	GetChannelByName(ctx context.Context, name sql.NullString) (*Channel, error)
	ListChannels(ctx context.Context) ([]*Channel, error)
	// CreateChannel(ctx context.Context, arg CreateChannelParams) (*Channel, error)
	// UpdateChannel(ctx context.Context, arg UpdateChannelParams) (int64, error)
	// DeleteChannel(ctx context.Context, id string) (int64, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetById(ctx context.Context, id string) (*Channel, error) {
	channel, err := s.repo.GetChannelById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel by id: %w", err)
	}

	return channel, nil
}

func (s *Service) GetByName(ctx context.Context, name string) (*Channel, error) {
	channel, err := s.repo.GetChannelByName(ctx, sql.NullString{String: name, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get channel by name: %w", err)
	}

	return channel, nil
}

func (s *Service) List(ctx context.Context) ([]*Channel, error) {
	channels, err := s.repo.ListChannels(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}

	return channels, nil
}

// func (s *Service) Create(ctx context.Context, channel *Channel) (*Channel, error) {
// 	channel, err := s.repo.CreateChannel(ctx, channel)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create channel: %w", err)
// 	}

// 	return channel, nil
// }
