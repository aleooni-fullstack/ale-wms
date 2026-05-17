package application

import (
	"context"

	dderr "github.com/aleodoni/go-ddd/errors"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
)

type InventoryBalanceService struct {
	balanceRepo      domain.InventoryBalanceRepository
	stockBalanceRepo domain.StockBalanceRepository
	movementRepo     domain.StockMovementRepository
}

func NewInventoryBalanceService(
	balanceRepo domain.InventoryBalanceRepository,
	stockBalanceRepo domain.StockBalanceRepository,
	movementRepo domain.StockMovementRepository,
) *InventoryBalanceService {
	return &InventoryBalanceService{
		balanceRepo:      balanceRepo,
		stockBalanceRepo: stockBalanceRepo,
		movementRepo:     movementRepo,
	}
}

type CreateInventoryBalanceInput struct {
	LocationID string
	Note       string
}

type AddItemInput struct {
	InventoryBalanceID string
	ProductID          string
}

type CountItemInput struct {
	InventoryBalanceID string
	ProductID          string
	CountedQuantity    float64
}

type ListInventoryBalancesInput struct {
	LocationID string
	Page       int32
	PerPage    int32
}

type ListInventoryBalancesOutput struct {
	Data    []*domain.InventoryBalance
	Page    int32
	PerPage int32
}

func (s *InventoryBalanceService) Create(ctx context.Context, input CreateInventoryBalanceInput) (*domain.InventoryBalance, error) {
	balance, err := domain.NewInventoryBalance(input.LocationID, input.Note)
	if err != nil {
		return nil, err
	}

	if err := s.balanceRepo.Create(ctx, balance); err != nil {
		return nil, err
	}

	return balance, nil
}

func (s *InventoryBalanceService) GetByID(ctx context.Context, id string) (*domain.InventoryBalance, error) {
	balance, err := s.balanceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.balanceRepo.FindAllItems(ctx, id)
	if err != nil {
		return nil, err
	}

	balance.Items = items

	return balance, nil
}

func (s *InventoryBalanceService) List(ctx context.Context, input ListInventoryBalancesInput) (*ListInventoryBalancesOutput, error) {
	if input.PerPage == 0 {
		input.PerPage = 20
	}
	if input.Page == 0 {
		input.Page = 1
	}

	if input.LocationID == "" {
		return nil, dderr.New("INVALID_FILTER", "location_id is required", nil)
	}

	offset := (input.Page - 1) * input.PerPage

	balances, err := s.balanceRepo.FindAllByLocation(ctx, input.LocationID, input.PerPage, offset)
	if err != nil {
		return nil, err
	}

	return &ListInventoryBalancesOutput{
		Data:    balances,
		Page:    input.Page,
		PerPage: input.PerPage,
	}, nil
}

func (s *InventoryBalanceService) Start(ctx context.Context, id string) (*domain.InventoryBalance, error) {
	balance, err := s.balanceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := balance.Start(); err != nil {
		return nil, err
	}

	if err := s.balanceRepo.UpdateStatus(ctx, balance); err != nil {
		return nil, err
	}

	return balance, nil
}

func (s *InventoryBalanceService) AddItem(ctx context.Context, input AddItemInput) (*domain.InventoryBalanceItem, error) {
	balance, err := s.balanceRepo.FindByID(ctx, input.InventoryBalanceID)
	if err != nil {
		return nil, err
	}

	if balance.Status != domain.InventoryBalanceStatusDraft && balance.Status != domain.InventoryBalanceStatusInProgress {
		return nil, dderr.New("INVALID_STATUS", "items can only be added to draft or in_progress balances", nil)
	}

	_, err = s.balanceRepo.FindItem(ctx, input.InventoryBalanceID, input.ProductID)
	if err == nil {
		return nil, dderr.New("ITEM_ALREADY_EXISTS", "product already added to this balance", nil)
	}
	if err != dderr.ErrNotFound {
		return nil, err
	}

	stockBalance, err := s.stockBalanceRepo.FindByProductAndLocation(ctx, input.ProductID, balance.LocationID)
	systemQty := 0.0
	if err == nil {
		systemQty = stockBalance.Quantity
	} else if err != dderr.ErrNotFound {
		return nil, err
	}

	item, err := domain.NewInventoryBalanceItem(input.InventoryBalanceID, input.ProductID, systemQty)
	if err != nil {
		return nil, err
	}

	if err := s.balanceRepo.AddItem(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *InventoryBalanceService) CountItem(ctx context.Context, input CountItemInput) (*domain.InventoryBalanceItem, error) {
	balance, err := s.balanceRepo.FindByID(ctx, input.InventoryBalanceID)
	if err != nil {
		return nil, err
	}

	if balance.Status != domain.InventoryBalanceStatusInProgress {
		return nil, dderr.New("INVALID_STATUS", "items can only be counted in in_progress balances", nil)
	}

	item, err := s.balanceRepo.FindItem(ctx, input.InventoryBalanceID, input.ProductID)
	if err != nil {
		return nil, err
	}

	if err := item.Count(input.CountedQuantity); err != nil {
		return nil, err
	}

	if err := s.balanceRepo.UpdateItemCountedQuantity(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *InventoryBalanceService) Complete(ctx context.Context, id string) (*domain.InventoryBalance, error) {
	balance, err := s.balanceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := balance.Complete(); err != nil {
		return nil, err
	}

	items, err := s.balanceRepo.FindAllItems(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.CountedQuantity == nil {
			continue
		}

		movement, err := domain.NewStockMovement(
			item.ProductID,
			balance.LocationID,
			domain.MovementTypeAdjustment,
			*item.CountedQuantity,
			"inventory balance adjustment: "+balance.ID,
		)
		if err != nil {
			return nil, err
		}

		if err := s.movementRepo.Create(ctx, movement); err != nil {
			return nil, err
		}

		stockBalance, err := s.stockBalanceRepo.FindByProductAndLocation(ctx, item.ProductID, balance.LocationID)
		if err != nil && err != dderr.ErrNotFound {
			return nil, err
		}

		if err == dderr.ErrNotFound || stockBalance == nil {
			stockBalance = domain.NewStockBalance(item.ProductID, balance.LocationID, 0)
		}

		stockBalance.Apply(movement)

		if err := s.stockBalanceRepo.Upsert(ctx, stockBalance); err != nil {
			return nil, err
		}
	}

	if err := s.balanceRepo.UpdateStatus(ctx, balance); err != nil {
		return nil, err
	}

	balance.Items = items

	return balance, nil
}

func (s *InventoryBalanceService) Cancel(ctx context.Context, id string) (*domain.InventoryBalance, error) {
	balance, err := s.balanceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := balance.Cancel(); err != nil {
		return nil, err
	}

	if err := s.balanceRepo.UpdateStatus(ctx, balance); err != nil {
		return nil, err
	}

	return balance, nil
}
