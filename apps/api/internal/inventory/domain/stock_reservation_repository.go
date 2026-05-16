package domain

import (
	"context"
)

type StockReservationRepository interface {
	Create(ctx context.Context, r *StockReservation) error
	FindByID(ctx context.Context, id string) (*StockReservation, error)
	FindAllByProduct(ctx context.Context, productID string, limit, offset int32) ([]*StockReservation, error)
	FindAllByLocation(ctx context.Context, locationID string, limit, offset int32) ([]*StockReservation, error)
	FindAllByReference(ctx context.Context, reference string) ([]*StockReservation, error)
	UpdateStatus(ctx context.Context, r *StockReservation) error
}
