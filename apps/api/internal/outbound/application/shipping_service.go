package application

import (
	"context"

	dderr "github.com/aleodoni/go-ddd/errors"

	inventorydomain "github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/domain"
)

type ShippingService struct {
	shippingRepo domain.ShippingRepository
	packingRepo  domain.PackingRepository
	orderRepo    domain.OrderRepository
	balanceRepo  inventorydomain.StockBalanceRepository
	movementRepo inventorydomain.StockMovementRepository
}

func NewShippingService(
	shippingRepo domain.ShippingRepository,
	packingRepo domain.PackingRepository,
	orderRepo domain.OrderRepository,
	balanceRepo inventorydomain.StockBalanceRepository,
	movementRepo inventorydomain.StockMovementRepository,
) *ShippingService {
	return &ShippingService{
		shippingRepo: shippingRepo,
		packingRepo:  packingRepo,
		orderRepo:    orderRepo,
		balanceRepo:  balanceRepo,
		movementRepo: movementRepo,
	}
}

func (s *ShippingService) Create(ctx context.Context, orderID string) (*domain.Shipping, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status != domain.OrderStatusPacking {
		return nil, dderr.New("INVALID_STATUS", "only packing orders can start shipping", nil)
	}

	packing, err := s.packingRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if packing.Status != domain.PackingStatusCompleted {
		return nil, dderr.New("INVALID_STATUS", "packing must be completed before shipping", nil)
	}

	shipping, err := domain.NewShipping(orderID, packing.ID, "")
	if err != nil {
		return nil, err
	}

	if err := s.shippingRepo.Create(ctx, shipping); err != nil {
		return nil, err
	}

	return shipping, nil
}

func (s *ShippingService) GetByID(ctx context.Context, id string) (*domain.Shipping, error) {
	return s.shippingRepo.FindByID(ctx, id)
}

func (s *ShippingService) Ship(ctx context.Context, id, trackingCode string) (*domain.Shipping, error) {
	shipping, err := s.shippingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := shipping.Ship(trackingCode); err != nil {
		return nil, err
	}

	if err := s.shippingRepo.UpdateStatus(ctx, shipping); err != nil {
		return nil, err
	}

	order, err := s.orderRepo.FindByID(ctx, shipping.OrderID)
	if err != nil {
		return nil, err
	}

	items, err := s.orderRepo.FindAllItems(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		movement, err := inventorydomain.NewStockMovement(
			item.ProductID,
			item.LocationID,
			inventorydomain.MovementTypeOut,
			item.Quantity,
			"shipped order: "+order.Reference,
		)
		if err != nil {
			return nil, err
		}

		if err := s.movementRepo.Create(ctx, movement); err != nil {
			return nil, err
		}

		balance, err := s.balanceRepo.FindByProductAndLocation(ctx, item.ProductID, item.LocationID)
		if err != nil && err != dderr.ErrNotFound {
			return nil, err
		}

		if err == dderr.ErrNotFound || balance == nil {
			balance = inventorydomain.NewStockBalance(item.ProductID, item.LocationID, 0)
		}

		balance.Apply(movement)

		if err := s.balanceRepo.Upsert(ctx, balance); err != nil {
			return nil, err
		}
	}

	order.Ship()
	if err := s.orderRepo.UpdateStatus(ctx, order); err != nil {
		return nil, err
	}

	return shipping, nil
}

func (s *ShippingService) Cancel(ctx context.Context, id string) (*domain.Shipping, error) {
	shipping, err := s.shippingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := shipping.Cancel(); err != nil {
		return nil, err
	}

	if err := s.shippingRepo.UpdateStatus(ctx, shipping); err != nil {
		return nil, err
	}

	order, err := s.orderRepo.FindByID(ctx, shipping.OrderID)
	if err != nil {
		return nil, err
	}

	if err := order.Cancel(); err != nil {
		return nil, err
	}

	if err := s.orderRepo.UpdateStatus(ctx, order); err != nil {
		return nil, err
	}

	return shipping, nil
}
