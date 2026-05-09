package application

import (
	"context"

	dderr "github.com/aleodoni/go-ddd/errors"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/domain"
)

type WarehouseService struct {
	repo domain.WarehouseRepository
}

func NewWarehouseService(repo domain.WarehouseRepository) *WarehouseService {
	return &WarehouseService{repo: repo}
}

type CreateWarehouseInput struct {
	Code    string
	Name    string
	Address string
}

type UpdateWarehouseInput struct {
	Code    string
	Name    string
	Address string
}

type ListWarehousesInput struct {
	Page    int32
	PerPage int32
}

type ListWarehousesOutput struct {
	Data    []*domain.Warehouse
	Page    int32
	PerPage int32
}

func (s *WarehouseService) Create(ctx context.Context, input CreateWarehouseInput) (*domain.Warehouse, error) {
	_, err := s.repo.FindByCode(ctx, input.Code)
	if err == nil {
		return nil, dderr.New("CODE_ALREADY_EXISTS", "a warehouse with this code already exists", nil)
	}
	if err != dderr.ErrNotFound {
		return nil, err
	}

	warehouse, err := domain.NewWarehouse(input.Code, input.Name, input.Address)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, warehouse); err != nil {
		return nil, err
	}

	return warehouse, nil
}

func (s *WarehouseService) GetByID(ctx context.Context, id string) (*domain.Warehouse, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *WarehouseService) List(ctx context.Context, input ListWarehousesInput) (*ListWarehousesOutput, error) {
	if input.PerPage == 0 {
		input.PerPage = 20
	}
	if input.Page == 0 {
		input.Page = 1
	}

	offset := (input.Page - 1) * input.PerPage

	warehouses, err := s.repo.FindAll(ctx, input.PerPage, offset)
	if err != nil {
		return nil, err
	}

	return &ListWarehousesOutput{
		Data:    warehouses,
		Page:    input.Page,
		PerPage: input.PerPage,
	}, nil
}

func (s *WarehouseService) Update(ctx context.Context, id string, input UpdateWarehouseInput) (*domain.Warehouse, error) {
	warehouse, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if warehouse.Code != input.Code {
		_, err := s.repo.FindByCode(ctx, input.Code)
		if err == nil {
			return nil, dderr.New("CODE_ALREADY_EXISTS", "a warehouse with this code already exists", nil)
		}
		if err != dderr.ErrNotFound {
			return nil, err
		}
	}

	if err := warehouse.Update(input.Code, input.Name, input.Address); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, warehouse); err != nil {
		return nil, err
	}

	return warehouse, nil
}

func (s *WarehouseService) Deactivate(ctx context.Context, id string) error {
	warehouse, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	warehouse.Deactivate()

	return s.repo.Delete(ctx, warehouse)
}
