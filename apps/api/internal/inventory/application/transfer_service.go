package application

import (
	"context"

	dderr "github.com/aleodoni/go-ddd/errors"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
)

type TransferService struct {
	transferRepo domain.StockTransferRepository
	balanceRepo  domain.StockBalanceRepository
	movementRepo domain.StockMovementRepository
}

func NewTransferService(
	transferRepo domain.StockTransferRepository,
	balanceRepo domain.StockBalanceRepository,
	movementRepo domain.StockMovementRepository,
) *TransferService {
	return &TransferService{
		transferRepo: transferRepo,
		balanceRepo:  balanceRepo,
		movementRepo: movementRepo,
	}
}

type CreateTransferInput struct {
	ProductID      string
	FromLocationID string
	ToLocationID   string
	Quantity       float64
	Note           string
}

type ListTransfersInput struct {
	ProductID      string
	FromLocationID string
	ToLocationID   string
	Page           int32
	PerPage        int32
}

type ListTransfersOutput struct {
	Data    []*domain.StockTransfer
	Page    int32
	PerPage int32
}

func (s *TransferService) Create(ctx context.Context, input CreateTransferInput) (*domain.StockTransfer, error) {
	balance, err := s.balanceRepo.FindByProductAndLocation(ctx, input.ProductID, input.FromLocationID)
	if err != nil && err != dderr.ErrNotFound {
		return nil, err
	}

	currentQty := 0.0
	if balance != nil {
		currentQty = balance.Quantity
	}

	if currentQty < input.Quantity {
		return nil, dderr.New("INSUFFICIENT_STOCK", "insufficient stock for this transfer", nil)
	}

	transfer, err := domain.NewStockTransfer(
		input.ProductID,
		input.FromLocationID,
		input.ToLocationID,
		input.Quantity,
		input.Note,
	)
	if err != nil {
		return nil, err
	}

	if err := s.transferRepo.Create(ctx, transfer); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *TransferService) Complete(ctx context.Context, id string) (*domain.StockTransfer, error) {
	transfer, err := s.transferRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := transfer.Complete(); err != nil {
		return nil, err
	}

	outMovement, err := domain.NewStockMovement(
		transfer.ProductID,
		transfer.FromLocationID,
		domain.MovementTypeOut,
		transfer.Quantity,
		"transfer out: "+transfer.ID,
	)
	if err != nil {
		return nil, err
	}

	if err := s.movementRepo.Create(ctx, outMovement); err != nil {
		return nil, err
	}

	fromBalance, err := s.balanceRepo.FindByProductAndLocation(ctx, transfer.ProductID, transfer.FromLocationID)
	if err != nil && err != dderr.ErrNotFound {
		return nil, err
	}
	if err == dderr.ErrNotFound || fromBalance == nil {
		fromBalance = domain.NewStockBalance(transfer.ProductID, transfer.FromLocationID, 0)
	}
	fromBalance.Apply(outMovement)
	if err := s.balanceRepo.Upsert(ctx, fromBalance); err != nil {
		return nil, err
	}

	inMovement, err := domain.NewStockMovement(
		transfer.ProductID,
		transfer.ToLocationID,
		domain.MovementTypeIn,
		transfer.Quantity,
		"transfer in: "+transfer.ID,
	)
	if err != nil {
		return nil, err
	}

	if err := s.movementRepo.Create(ctx, inMovement); err != nil {
		return nil, err
	}

	toBalance, err := s.balanceRepo.FindByProductAndLocation(ctx, transfer.ProductID, transfer.ToLocationID)
	if err != nil && err != dderr.ErrNotFound {
		return nil, err
	}
	if err == dderr.ErrNotFound || toBalance == nil {
		toBalance = domain.NewStockBalance(transfer.ProductID, transfer.ToLocationID, 0)
	}
	toBalance.Apply(inMovement)
	if err := s.balanceRepo.Upsert(ctx, toBalance); err != nil {
		return nil, err
	}

	if err := s.transferRepo.UpdateStatus(ctx, transfer); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *TransferService) Cancel(ctx context.Context, id string) (*domain.StockTransfer, error) {
	transfer, err := s.transferRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := transfer.Cancel(); err != nil {
		return nil, err
	}

	if err := s.transferRepo.UpdateStatus(ctx, transfer); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *TransferService) GetByID(ctx context.Context, id string) (*domain.StockTransfer, error) {
	return s.transferRepo.FindByID(ctx, id)
}

func (s *TransferService) List(ctx context.Context, input ListTransfersInput) (*ListTransfersOutput, error) {
	if input.PerPage == 0 {
		input.PerPage = 20
	}
	if input.Page == 0 {
		input.Page = 1
	}

	offset := (input.Page - 1) * input.PerPage

	var transfers []*domain.StockTransfer
	var err error

	switch {
	case input.ProductID != "":
		transfers, err = s.transferRepo.FindAllByProduct(ctx, input.ProductID, input.PerPage, offset)
	case input.FromLocationID != "":
		transfers, err = s.transferRepo.FindAllByFromLocation(ctx, input.FromLocationID, input.PerPage, offset)
	case input.ToLocationID != "":
		transfers, err = s.transferRepo.FindAllByToLocation(ctx, input.ToLocationID, input.PerPage, offset)
	default:
		return nil, dderr.New("INVALID_FILTER", "product_id, from_location_id or to_location_id is required", nil)
	}

	if err != nil {
		return nil, err
	}

	return &ListTransfersOutput{
		Data:    transfers,
		Page:    input.Page,
		PerPage: input.PerPage,
	}, nil
}
