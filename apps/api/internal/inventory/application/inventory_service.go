package application

import (
	"context"

	dderr "github.com/aleodoni/go-ddd/errors"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
)

type InventoryService struct {
	movementRepo domain.StockMovementRepository
	balanceRepo  domain.StockBalanceRepository
}

func NewInventoryService(movementRepo domain.StockMovementRepository, balanceRepo domain.StockBalanceRepository) *InventoryService {
	return &InventoryService{
		movementRepo: movementRepo,
		balanceRepo:  balanceRepo,
	}
}

type RegisterMovementInput struct {
	ProductID  string
	LocationID string
	Type       string
	Quantity   float64
	Note       string
}

type ListMovementsInput struct {
	ProductID  string
	LocationID string
	Page       int32
	PerPage    int32
}

type ListMovementsOutput struct {
	Data    []*domain.StockMovement
	Page    int32
	PerPage int32
}

func (s *InventoryService) RegisterMovement(ctx context.Context, input RegisterMovementInput) (*domain.StockMovement, error) {
	movementType := domain.MovementType(input.Type)

	movement, err := domain.NewStockMovement(
		input.ProductID,
		input.LocationID,
		movementType,
		input.Quantity,
		input.Note,
	)
	if err != nil {
		return nil, err
	}

	if movementType == domain.MovementTypeOut {
		balance, err := s.balanceRepo.FindByProductAndLocation(ctx, input.ProductID, input.LocationID)
		if err != nil && err != dderr.ErrNotFound {
			return nil, err
		}

		currentQty := 0.0
		if balance != nil {
			currentQty = balance.Quantity
		}

		if currentQty < input.Quantity {
			return nil, dderr.New("INSUFFICIENT_STOCK", "insufficient stock for this operation", nil)
		}
	}

	if err := s.movementRepo.Create(ctx, movement); err != nil {
		return nil, err
	}

	balance, err := s.balanceRepo.FindByProductAndLocation(ctx, input.ProductID, input.LocationID)
	if err != nil && err != dderr.ErrNotFound {
		return nil, err
	}

	if err == dderr.ErrNotFound || balance == nil {
		balance = domain.NewStockBalance(input.ProductID, input.LocationID, 0)
	}

	balance.Apply(movement)

	if err := s.balanceRepo.Upsert(ctx, balance); err != nil {
		return nil, err
	}

	return movement, nil
}

func (s *InventoryService) GetBalance(ctx context.Context, productID, locationID string) (*domain.StockBalance, error) {
	return s.balanceRepo.FindByProductAndLocation(ctx, productID, locationID)
}

func (s *InventoryService) ListBalancesByProduct(ctx context.Context, productID string) ([]*domain.StockBalance, error) {
	return s.balanceRepo.FindAllByProduct(ctx, productID)
}

func (s *InventoryService) ListBalancesByLocation(ctx context.Context, locationID string) ([]*domain.StockBalance, error) {
	return s.balanceRepo.FindAllByLocation(ctx, locationID)
}

func (s *InventoryService) ListMovements(ctx context.Context, input ListMovementsInput) (*ListMovementsOutput, error) {
	if input.PerPage == 0 {
		input.PerPage = 20
	}
	if input.Page == 0 {
		input.Page = 1
	}

	offset := (input.Page - 1) * input.PerPage

	var movements []*domain.StockMovement
	var err error

	switch {
	case input.ProductID != "":
		movements, err = s.movementRepo.FindAllByProduct(ctx, input.ProductID, input.PerPage, offset)
	case input.LocationID != "":
		movements, err = s.movementRepo.FindAllByLocation(ctx, input.LocationID, input.PerPage, offset)
	default:
		return nil, dderr.New("INVALID_FILTER", "product_id or location_id is required", nil)
	}

	if err != nil {
		return nil, err
	}

	return &ListMovementsOutput{
		Data:    movements,
		Page:    input.Page,
		PerPage: input.PerPage,
	}, nil
}
