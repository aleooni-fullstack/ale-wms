package application

import (
	"context"

	dderr "github.com/aleodoni/go-ddd/errors"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
)

type ReservationService struct {
	reservationRepo domain.StockReservationRepository
	balanceRepo     domain.StockBalanceRepository
	movementRepo    domain.StockMovementRepository
}

func NewReservationService(
	reservationRepo domain.StockReservationRepository,
	balanceRepo domain.StockBalanceRepository,
	movementRepo domain.StockMovementRepository,
) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
		balanceRepo:     balanceRepo,
		movementRepo:    movementRepo,
	}
}

type CreateReservationInput struct {
	ProductID  string
	LocationID string
	Quantity   float64
	Reference  string
	Note       string
}

type ListReservationsInput struct {
	ProductID  string
	LocationID string
	Reference  string
	Page       int32
	PerPage    int32
}

type ListReservationsOutput struct {
	Data    []*domain.StockReservation
	Page    int32
	PerPage int32
}

func (s *ReservationService) Create(ctx context.Context, input CreateReservationInput) (*domain.StockReservation, error) {
	balance, err := s.balanceRepo.FindByProductAndLocation(ctx, input.ProductID, input.LocationID)
	if err != nil && err != dderr.ErrNotFound {
		return nil, err
	}

	if err == dderr.ErrNotFound || balance == nil {
		return nil, dderr.New("INSUFFICIENT_STOCK", "no stock available for this product in this location", nil)
	}

	if balance.AvailableQuantity() < input.Quantity {
		return nil, dderr.New("INSUFFICIENT_STOCK", "insufficient available stock for this reservation", nil)
	}

	reservation, err := domain.NewStockReservation(
		input.ProductID,
		input.LocationID,
		input.Quantity,
		input.Reference,
		input.Note,
	)
	if err != nil {
		return nil, err
	}

	if err := s.reservationRepo.Create(ctx, reservation); err != nil {
		return nil, err
	}

	balance.Reserve(input.Quantity)

	if err := s.balanceRepo.Upsert(ctx, balance); err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *ReservationService) Confirm(ctx context.Context, id string) (*domain.StockReservation, error) {
	reservation, err := s.reservationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := reservation.Confirm(); err != nil {
		return nil, err
	}

	if err := s.reservationRepo.UpdateStatus(ctx, reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *ReservationService) Fulfill(ctx context.Context, id string) (*domain.StockReservation, error) {
	reservation, err := s.reservationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := reservation.Fulfill(); err != nil {
		return nil, err
	}

	outMovement, err := domain.NewStockMovement(
		reservation.ProductID,
		reservation.LocationID,
		domain.MovementTypeOut,
		reservation.Quantity,
		"fulfillment of reservation: "+reservation.ID,
	)
	if err != nil {
		return nil, err
	}

	if err := s.movementRepo.Create(ctx, outMovement); err != nil {
		return nil, err
	}

	balance, err := s.balanceRepo.FindByProductAndLocation(ctx, reservation.ProductID, reservation.LocationID)
	if err != nil {
		return nil, err
	}

	balance.Apply(outMovement)
	balance.Release(reservation.Quantity)

	if err := s.balanceRepo.Upsert(ctx, balance); err != nil {
		return nil, err
	}

	if err := s.reservationRepo.UpdateStatus(ctx, reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *ReservationService) Release(ctx context.Context, id string) (*domain.StockReservation, error) {
	reservation, err := s.reservationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := reservation.Release(); err != nil {
		return nil, err
	}

	balance, err := s.balanceRepo.FindByProductAndLocation(ctx, reservation.ProductID, reservation.LocationID)
	if err != nil {
		return nil, err
	}

	balance.Release(reservation.Quantity)

	if err := s.balanceRepo.Upsert(ctx, balance); err != nil {
		return nil, err
	}

	if err := s.reservationRepo.UpdateStatus(ctx, reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *ReservationService) Cancel(ctx context.Context, id string) (*domain.StockReservation, error) {
	reservation, err := s.reservationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := reservation.Cancel(); err != nil {
		return nil, err
	}

	balance, err := s.balanceRepo.FindByProductAndLocation(ctx, reservation.ProductID, reservation.LocationID)
	if err != nil {
		return nil, err
	}

	balance.Release(reservation.Quantity)

	if err := s.balanceRepo.Upsert(ctx, balance); err != nil {
		return nil, err
	}

	if err := s.reservationRepo.UpdateStatus(ctx, reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *ReservationService) GetByID(ctx context.Context, id string) (*domain.StockReservation, error) {
	return s.reservationRepo.FindByID(ctx, id)
}

func (s *ReservationService) List(ctx context.Context, input ListReservationsInput) (*ListReservationsOutput, error) {
	if input.PerPage == 0 {
		input.PerPage = 20
	}
	if input.Page == 0 {
		input.Page = 1
	}

	offset := (input.Page - 1) * input.PerPage

	var reservations []*domain.StockReservation
	var err error

	switch {
	case input.Reference != "":
		reservations, err = s.reservationRepo.FindAllByReference(ctx, input.Reference)
	case input.ProductID != "":
		reservations, err = s.reservationRepo.FindAllByProduct(ctx, input.ProductID, input.PerPage, offset)
	case input.LocationID != "":
		reservations, err = s.reservationRepo.FindAllByLocation(ctx, input.LocationID, input.PerPage, offset)
	default:
		return nil, dderr.New("INVALID_FILTER", "product_id, location_id or reference is required", nil)
	}

	if err != nil {
		return nil, err
	}

	return &ListReservationsOutput{
		Data:    reservations,
		Page:    input.Page,
		PerPage: input.PerPage,
	}, nil
}
