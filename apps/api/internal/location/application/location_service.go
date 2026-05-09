package application

import (
	"context"

	dderr "github.com/aleodoni/go-ddd/errors"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/domain"
)

type LocationService struct {
	repo     domain.LocationRepository
	zoneRepo domain.ZoneRepository
}

func NewLocationService(repo domain.LocationRepository, zoneRepo domain.ZoneRepository) *LocationService {
	return &LocationService{repo: repo, zoneRepo: zoneRepo}
}

type CreateLocationInput struct {
	ZoneID string
	Code   string
	Name   string
}

type UpdateLocationInput struct {
	Code string
	Name string
}

type ListLocationsInput struct {
	ZoneID  string
	Page    int32
	PerPage int32
}

type ListLocationsOutput struct {
	Data    []*domain.Location
	Page    int32
	PerPage int32
}

func (s *LocationService) Create(ctx context.Context, input CreateLocationInput) (*domain.Location, error) {
	_, err := s.zoneRepo.FindByID(ctx, input.ZoneID)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.FindByCode(ctx, input.ZoneID, input.Code)
	if err == nil {
		return nil, dderr.New("CODE_ALREADY_EXISTS", "a location with this code already exists in this zone", nil)
	}
	if err != dderr.ErrNotFound {
		return nil, err
	}

	location, err := domain.NewLocation(input.ZoneID, input.Code, input.Name)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, location); err != nil {
		return nil, err
	}

	return location, nil
}

func (s *LocationService) GetByID(ctx context.Context, id string) (*domain.Location, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *LocationService) List(ctx context.Context, input ListLocationsInput) (*ListLocationsOutput, error) {
	if input.PerPage == 0 {
		input.PerPage = 20
	}
	if input.Page == 0 {
		input.Page = 1
	}

	offset := (input.Page - 1) * input.PerPage

	locations, err := s.repo.FindAllByZone(ctx, input.ZoneID, input.PerPage, offset)
	if err != nil {
		return nil, err
	}

	return &ListLocationsOutput{
		Data:    locations,
		Page:    input.Page,
		PerPage: input.PerPage,
	}, nil
}

func (s *LocationService) Update(ctx context.Context, id string, input UpdateLocationInput) (*domain.Location, error) {
	location, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if location.Code != input.Code {
		_, err := s.repo.FindByCode(ctx, location.ZoneID, input.Code)
		if err == nil {
			return nil, dderr.New("CODE_ALREADY_EXISTS", "a location with this code already exists in this zone", nil)
		}
		if err != dderr.ErrNotFound {
			return nil, err
		}
	}

	if err := location.Update(input.Code, input.Name); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, location); err != nil {
		return nil, err
	}

	return location, nil
}

func (s *LocationService) Deactivate(ctx context.Context, id string) error {
	location, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	location.Deactivate()

	return s.repo.Delete(ctx, location)
}
