package tenants

import (
	"context"
	"errors"
	"fmt"
)

var ErrRecordNotFound = errors.New("record not found")

type Repository interface {
	GetById(ctx context.Context, id string) (*Tenant, error)
	GetByName(ctx context.Context, name string) (*Tenant, error)
	List(ctx context.Context) ([]*Tenant, error)
	Create(ctx context.Context, tenant *Tenant) (*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) error
	SoftDelete(ctx context.Context, id string) error
}

type Service struct {
	repo       Repository
	generateID TenantIDGenerator
}

func NewService(repo Repository, idGenerator TenantIDGenerator) *Service {
	return &Service{repo: repo, generateID: idGenerator}
}

// create

type CreateTenantArgs struct {
	Name        string
	Description string
}

type CreateTenantResult struct {
	Tenant Tenant
}

func (s *Service) GetById(ctx context.Context, id string) (*Tenant, error) {
	tenant, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant by id: %w", err)
	}

	return tenant, nil
}

func (s *Service) GetByName(ctx context.Context, name string) (*Tenant, error) {
	tenant, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant by name: %w", err)
	}

	return tenant, nil
}

func (s *Service) List(ctx context.Context) ([]*Tenant, error) {
	tenants, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	return tenants, nil
}

func (s *Service) Create(ctx context.Context, tenant *Tenant) (*Tenant, error) {
	tenant, err := s.repo.Create(ctx, tenant)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return tenant, nil
}
