package application

import (
	"context"

	dderr "github.com/aleodoni/go-ddd/errors"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/domain"
)

type ZoneService struct {
	repo          domain.ZoneRepository
	warehouseRepo domain.WarehouseRepository
}

func NewZoneService(repo domain.ZoneRepository, warehouseRepo domain.WarehouseRepository) *ZoneService {
	return &ZoneService{repo: repo, warehouseRepo: warehouseRepo}
}

type CreateZoneInput struct {
	WarehouseID string
	Code        string
	Name        string
}

type UpdateZoneInput struct {
	Code string
	Name string
}

type ListZonesInput struct {
	WarehouseID string
	Page        int32
	PerPage     int32
}

type ListZonesOutput struct {
	Data    []*domain.Zone
	Page    int32
	PerPage int32
}

func (s *ZoneService) Create(ctx context.Context, input CreateZoneInput) (*domain.Zone, error) {
	_, err := s.warehouseRepo.FindByID(ctx, input.WarehouseID)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.FindByCode(ctx, input.WarehouseID, input.Code)
	if err == nil {
		return nil, dderr.New("CODE_ALREADY_EXISTS", "a zone with this code already exists in this warehouse", nil)
	}
	if err != dderr.ErrNotFound {
		return nil, err
	}

	zone, err := domain.NewZone(input.WarehouseID, input.Code, input.Name)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, zone); err != nil {
		return nil, err
	}

	return zone, nil
}

func (s *ZoneService) GetByID(ctx context.Context, id string) (*domain.Zone, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ZoneService) List(ctx context.Context, input ListZonesInput) (*ListZonesOutput, error) {
	if input.PerPage == 0 {
		input.PerPage = 20
	}
	if input.Page == 0 {
		input.Page = 1
	}

	offset := (input.Page - 1) * input.PerPage

	zones, err := s.repo.FindAllByWarehouse(ctx, input.WarehouseID, input.PerPage, offset)
	if err != nil {
		return nil, err
	}

	return &ListZonesOutput{
		Data:    zones,
		Page:    input.Page,
		PerPage: input.PerPage,
	}, nil
}

func (s *ZoneService) Update(ctx context.Context, id string, input UpdateZoneInput) (*domain.Zone, error) {
	zone, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if zone.Code != input.Code {
		_, err := s.repo.FindByCode(ctx, zone.WarehouseID, input.Code)
		if err == nil {
			return nil, dderr.New("CODE_ALREADY_EXISTS", "a zone with this code already exists in this warehouse", nil)
		}
		if err != dderr.ErrNotFound {
			return nil, err
		}
	}

	if err := zone.Update(input.Code, input.Name); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, zone); err != nil {
		return nil, err
	}

	return zone, nil
}

func (s *ZoneService) Deactivate(ctx context.Context, id string) error {
	zone, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	zone.Deactivate()

	return s.repo.Delete(ctx, zone)
}
