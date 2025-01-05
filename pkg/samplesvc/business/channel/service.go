package channel

import (
	"context"
	"fmt"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*Channel, error)
	GetByName(ctx context.Context, name string) (*Channel, error)
	List(ctx context.Context) ([]*Channel, error)
	Create(ctx context.Context, channel *Channel) (*Channel, error)
	Update(ctx context.Context, channel *Channel) error
	SoftDelete(ctx context.Context, id string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetById(ctx context.Context, id string) (*Channel, error) {
	channel, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel by id: %w", err)
	}

	return channel, nil
}

func (s *Service) GetByName(ctx context.Context, name string) (*Channel, error) {
	channel, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel by name: %w", err)
	}

	return channel, nil
}

func (s *Service) List(ctx context.Context) ([]*Channel, error) {
	channels, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}

	return channels, nil
}

func (s *Service) Create(ctx context.Context, channel *Channel) (*Channel, error) {
	channel, err := s.repo.Create(ctx, channel)
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	return channel, nil
}
