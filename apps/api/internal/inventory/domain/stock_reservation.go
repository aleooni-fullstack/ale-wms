package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/nrednav/cuid2"
)

type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "PENDING"
	ReservationStatusConfirmed ReservationStatus = "CONFIRMED"
	ReservationStatusFulfilled ReservationStatus = "FULFILLED"
	ReservationStatusReleased  ReservationStatus = "RELEASED"
	ReservationStatusCancelled ReservationStatus = "CANCELLED"
)

type StockReservation struct {
	domain.AggregateRoot[string]
	ProductID  string
	LocationID string
	Quantity   float64
	Status     ReservationStatus
	Reference  string
	Note       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewStockReservation(productID, locationID string, quantity float64, reference, note string) (*StockReservation, error) {
	if productID == "" {
		return nil, dderr.New("INVALID_PRODUCT_ID", "product_id is required", nil)
	}
	if locationID == "" {
		return nil, dderr.New("INVALID_LOCATION_ID", "location_id is required", nil)
	}
	if quantity <= 0 {
		return nil, dderr.New("INVALID_QUANTITY", "quantity must be greater than zero", nil)
	}
	if reference == "" {
		return nil, dderr.New("INVALID_REFERENCE", "reference is required", nil)
	}

	return &StockReservation{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		ProductID:     productID,
		LocationID:    locationID,
		Quantity:      quantity,
		Status:        ReservationStatusPending,
		Reference:     reference,
		Note:          note,
	}, nil
}

func RestoreStockReservation(id, productID, locationID string, quantity float64, status ReservationStatus, reference, note string, createdAt, updatedAt time.Time) *StockReservation {
	return &StockReservation{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		ProductID:     productID,
		LocationID:    locationID,
		Quantity:      quantity,
		Status:        status,
		Reference:     reference,
		Note:          note,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (r *StockReservation) Confirm() error {
	if r.Status != ReservationStatusPending {
		return dderr.New("INVALID_STATUS", "only pending reservations can be confirmed", nil)
	}

	r.Status = ReservationStatusConfirmed

	return nil
}

func (r *StockReservation) Fulfill() error {
	if r.Status != ReservationStatusConfirmed {
		return dderr.New("INVALID_STATUS", "only confirmed reservations can be fulfilled", nil)
	}

	r.Status = ReservationStatusFulfilled

	return nil
}

func (r *StockReservation) Release() error {
	if r.Status != ReservationStatusPending && r.Status != ReservationStatusConfirmed {
		return dderr.New("INVALID_STATUS", "only pending or confirmed reservations can be released", nil)
	}

	r.Status = ReservationStatusReleased

	return nil
}

func (r *StockReservation) Cancel() error {
	if r.Status != ReservationStatusPending {
		return dderr.New("INVALID_STATUS", "only pending reservations can be cancelled", nil)
	}

	r.Status = ReservationStatusCancelled

	return nil
}
