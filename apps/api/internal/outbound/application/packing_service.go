package application

import (
	"context"

	dderr "github.com/aleodoni/go-ddd/errors"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/domain"
)

type PackingService struct {
	packingRepo domain.PackingRepository
	pickingRepo domain.PickingRepository
	orderRepo   domain.OrderRepository
}

func NewPackingService(
	packingRepo domain.PackingRepository,
	pickingRepo domain.PickingRepository,
	orderRepo domain.OrderRepository,
) *PackingService {
	return &PackingService{
		packingRepo: packingRepo,
		pickingRepo: pickingRepo,
		orderRepo:   orderRepo,
	}
}

func (s *PackingService) Create(ctx context.Context, orderID string) (*domain.Packing, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status != domain.OrderStatusPicking {
		return nil, dderr.New("INVALID_STATUS", "only picking orders can start packing", nil)
	}

	picking, err := s.pickingRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if picking.Status != domain.PickingStatusCompleted {
		return nil, dderr.New("INVALID_STATUS", "picking must be completed before packing", nil)
	}

	packing, err := domain.NewPacking(orderID, picking.ID, "")
	if err != nil {
		return nil, err
	}

	if err := s.packingRepo.Create(ctx, packing); err != nil {
		return nil, err
	}

	order.StartPacking()
	if err := s.orderRepo.UpdateStatus(ctx, order); err != nil {
		return nil, err
	}

	return packing, nil
}

func (s *PackingService) GetByID(ctx context.Context, id string) (*domain.Packing, error) {
	return s.packingRepo.FindByID(ctx, id)
}

func (s *PackingService) Start(ctx context.Context, id string) (*domain.Packing, error) {
	packing, err := s.packingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := packing.Start(); err != nil {
		return nil, err
	}

	if err := s.packingRepo.UpdateStatus(ctx, packing); err != nil {
		return nil, err
	}

	return packing, nil
}

func (s *PackingService) Complete(ctx context.Context, id string) (*domain.Packing, error) {
	packing, err := s.packingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := packing.Complete(); err != nil {
		return nil, err
	}

	if err := s.packingRepo.UpdateStatus(ctx, packing); err != nil {
		return nil, err
	}

	return packing, nil
}

func (s *PackingService) Cancel(ctx context.Context, id string) (*domain.Packing, error) {
	packing, err := s.packingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := packing.Cancel(); err != nil {
		return nil, err
	}

	if err := s.packingRepo.UpdateStatus(ctx, packing); err != nil {
		return nil, err
	}

	order, err := s.orderRepo.FindByID(ctx, packing.OrderID)
	if err != nil {
		return nil, err
	}

	if err := order.Cancel(); err != nil {
		return nil, err
	}

	if err := s.orderRepo.UpdateStatus(ctx, order); err != nil {
		return nil, err
	}

	return packing, nil
}
