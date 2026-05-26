package application

import (
	"context"

	dderr "github.com/aleodoni/go-ddd/errors"

	inventorydomain "github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/domain"
)

type PutAwayService struct {
	putAwayRepo  domain.PutAwayRepository
	receiptRepo  domain.ReceiptRepository
	poRepo       domain.PurchaseOrderRepository
	balanceRepo  inventorydomain.StockBalanceRepository
	movementRepo inventorydomain.StockMovementRepository
}

func NewPutAwayService(
	putAwayRepo domain.PutAwayRepository,
	receiptRepo domain.ReceiptRepository,
	poRepo domain.PurchaseOrderRepository,
	balanceRepo inventorydomain.StockBalanceRepository,
	movementRepo inventorydomain.StockMovementRepository,
) *PutAwayService {
	return &PutAwayService{
		putAwayRepo:  putAwayRepo,
		receiptRepo:  receiptRepo,
		poRepo:       poRepo,
		balanceRepo:  balanceRepo,
		movementRepo: movementRepo,
	}
}

type AddPutAwayItemInput struct {
	PutAwayID  string
	ProductID  string
	LocationID string
	Quantity   float64
}

type StoreItemInput struct {
	PutAwayID string
	ProductID string
}

func (s *PutAwayService) Create(ctx context.Context, receiptID string) (*domain.PutAway, error) {
	receipt, err := s.receiptRepo.FindByID(ctx, receiptID)
	if err != nil {
		return nil, err
	}

	if receipt.Status != domain.ReceiptStatusCompleted {
		return nil, dderr.New("INVALID_STATUS", "only completed receipts can have a put away", nil)
	}

	putAway, err := domain.NewPutAway(receiptID, "")
	if err != nil {
		return nil, err
	}

	if err := s.putAwayRepo.Create(ctx, putAway); err != nil {
		return nil, err
	}

	return putAway, nil
}

func (s *PutAwayService) GetByID(ctx context.Context, id string) (*domain.PutAway, error) {
	putAway, err := s.putAwayRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.putAwayRepo.FindAllItems(ctx, id)
	if err != nil {
		return nil, err
	}

	putAway.Items = items

	return putAway, nil
}

func (s *PutAwayService) Start(ctx context.Context, id string) (*domain.PutAway, error) {
	putAway, err := s.putAwayRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := putAway.Start(); err != nil {
		return nil, err
	}

	if err := s.putAwayRepo.UpdateStatus(ctx, putAway); err != nil {
		return nil, err
	}

	return putAway, nil
}

func (s *PutAwayService) AddItem(ctx context.Context, input AddPutAwayItemInput) (*domain.PutAwayItem, error) {
	putAway, err := s.putAwayRepo.FindByID(ctx, input.PutAwayID)
	if err != nil {
		return nil, err
	}

	if putAway.Status != domain.PutAwayStatusPending && putAway.Status != domain.PutAwayStatusInProgress {
		return nil, dderr.New("INVALID_STATUS", "items can only be added to pending or in_progress put aways", nil)
	}

	item, err := domain.NewPutAwayItem(input.PutAwayID, input.ProductID, input.LocationID, input.Quantity)
	if err != nil {
		return nil, err
	}

	if err := s.putAwayRepo.AddItem(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *PutAwayService) StoreItem(ctx context.Context, input StoreItemInput) (*domain.PutAwayItem, error) {
	putAway, err := s.putAwayRepo.FindByID(ctx, input.PutAwayID)
	if err != nil {
		return nil, err
	}

	if putAway.Status != domain.PutAwayStatusInProgress {
		return nil, dderr.New("INVALID_STATUS", "only in_progress put aways can store items", nil)
	}

	items, err := s.putAwayRepo.FindAllItems(ctx, input.PutAwayID)
	if err != nil {
		return nil, err
	}

	var target *domain.PutAwayItem
	for _, item := range items {
		if item.ProductID == input.ProductID {
			target = item
			break
		}
	}

	if target == nil {
		return nil, dderr.ErrNotFound
	}

	target.Store()

	if err := s.putAwayRepo.UpdateItemPutAway(ctx, target); err != nil {
		return nil, err
	}

	return target, nil
}

func (s *PutAwayService) Complete(ctx context.Context, id string) (*domain.PutAway, error) {
	putAway, err := s.putAwayRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.putAwayRepo.FindAllItems(ctx, id)
	if err != nil {
		return nil, err
	}

	putAway.Items = items

	if err := putAway.Complete(); err != nil {
		return nil, err
	}

	for _, item := range items {
		movement, err := inventorydomain.NewStockMovement(
			item.ProductID,
			item.LocationID,
			inventorydomain.MovementTypeIn,
			item.Quantity,
			"put away receipt: "+putAway.ReceiptID,
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

	if err := s.putAwayRepo.UpdateStatus(ctx, putAway); err != nil {
		return nil, err
	}

	receipt, err := s.receiptRepo.FindByID(ctx, putAway.ReceiptID)
	if err != nil {
		return nil, err
	}

	po, err := s.poRepo.FindByID(ctx, receipt.PurchaseOrderID)
	if err != nil {
		return nil, err
	}

	if err := po.Complete(); err != nil {
		return nil, err
	}

	if err := s.poRepo.UpdateStatus(ctx, po); err != nil {
		return nil, err
	}

	return putAway, nil
}

func (s *PutAwayService) Cancel(ctx context.Context, id string) (*domain.PutAway, error) {
	putAway, err := s.putAwayRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := putAway.Cancel(); err != nil {
		return nil, err
	}

	if err := s.putAwayRepo.UpdateStatus(ctx, putAway); err != nil {
		return nil, err
	}

	return putAway, nil
}
