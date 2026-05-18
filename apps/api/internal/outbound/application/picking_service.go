package application

import (
	"context"

	dderr "github.com/aleodoni/go-ddd/errors"

	inventorydomain "github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/domain"
)

type PickingService struct {
	pickingRepo  domain.PickingRepository
	orderRepo    domain.OrderRepository
	balanceRepo  inventorydomain.StockBalanceRepository
	movementRepo inventorydomain.StockMovementRepository
}

func NewPickingService(
	pickingRepo domain.PickingRepository,
	orderRepo domain.OrderRepository,
	balanceRepo inventorydomain.StockBalanceRepository,
	movementRepo inventorydomain.StockMovementRepository,
) *PickingService {
	return &PickingService{
		pickingRepo:  pickingRepo,
		orderRepo:    orderRepo,
		balanceRepo:  balanceRepo,
		movementRepo: movementRepo,
	}
}

func (s *PickingService) Create(ctx context.Context, orderID string) (*domain.Picking, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status != domain.OrderStatusConfirmed {
		return nil, dderr.New("INVALID_STATUS", "only confirmed orders can start picking", nil)
	}

	picking, err := domain.NewPicking(orderID, "")
	if err != nil {
		return nil, err
	}

	if err := s.pickingRepo.Create(ctx, picking); err != nil {
		return nil, err
	}

	items, err := s.orderRepo.FindAllItems(ctx, orderID)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		pickingItem, err := domain.NewPickingItem(picking.ID, item.ProductID, item.LocationID, item.Quantity)
		if err != nil {
			return nil, err
		}

		if err := s.pickingRepo.AddItem(ctx, pickingItem); err != nil {
			return nil, err
		}

		picking.Items = append(picking.Items, pickingItem)
	}

	order.StartPicking()
	if err := s.orderRepo.UpdateStatus(ctx, order); err != nil {
		return nil, err
	}

	return picking, nil
}

func (s *PickingService) GetByID(ctx context.Context, id string) (*domain.Picking, error) {
	picking, err := s.pickingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.pickingRepo.FindAllItems(ctx, id)
	if err != nil {
		return nil, err
	}

	picking.Items = items

	return picking, nil
}

func (s *PickingService) Start(ctx context.Context, id string) (*domain.Picking, error) {
	picking, err := s.pickingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := picking.Start(); err != nil {
		return nil, err
	}

	if err := s.pickingRepo.UpdateStatus(ctx, picking); err != nil {
		return nil, err
	}

	return picking, nil
}

func (s *PickingService) PickItem(ctx context.Context, pickingID, productID string) (*domain.PickingItem, error) {
	picking, err := s.pickingRepo.FindByID(ctx, pickingID)
	if err != nil {
		return nil, err
	}

	if picking.Status != domain.PickingStatusInProgress {
		return nil, dderr.New("INVALID_STATUS", "only in_progress pickings can have items picked", nil)
	}

	items, err := s.pickingRepo.FindAllItems(ctx, pickingID)
	if err != nil {
		return nil, err
	}

	var target *domain.PickingItem
	for _, item := range items {
		if item.ProductID == productID {
			target = item
			break
		}
	}

	if target == nil {
		return nil, dderr.ErrNotFound
	}

	target.Pick()

	if err := s.pickingRepo.UpdateItemPicked(ctx, target); err != nil {
		return nil, err
	}

	return target, nil
}

func (s *PickingService) Complete(ctx context.Context, id string) (*domain.Picking, error) {
	picking, err := s.pickingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.pickingRepo.FindAllItems(ctx, id)
	if err != nil {
		return nil, err
	}

	picking.Items = items

	if err := picking.Complete(); err != nil {
		return nil, err
	}

	if err := s.pickingRepo.UpdateStatus(ctx, picking); err != nil {
		return nil, err
	}

	return picking, nil
}

func (s *PickingService) Cancel(ctx context.Context, id string) (*domain.Picking, error) {
	picking, err := s.pickingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := picking.Cancel(); err != nil {
		return nil, err
	}

	if err := s.pickingRepo.UpdateStatus(ctx, picking); err != nil {
		return nil, err
	}

	order, err := s.orderRepo.FindByID(ctx, picking.OrderID)
	if err != nil {
		return nil, err
	}

	if err := order.Cancel(); err != nil {
		return nil, err
	}

	if err := s.orderRepo.UpdateStatus(ctx, order); err != nil {
		return nil, err
	}

	return picking, nil
}

func (s *PickingService) findBalance(ctx context.Context, productID, locationID string) (*inventorydomain.StockBalance, error) {
	balance, err := s.balanceRepo.FindByProductAndLocation(ctx, productID, locationID)
	if err != nil {
		return nil, err
	}
	return balance, nil
}
