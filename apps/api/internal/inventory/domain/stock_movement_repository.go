package domain

import (
	"context"
)

type StockMovementRepository interface {
	Create(ctx context.Context, m *StockMovement) error
	FindByID(ctx context.Context, id string) (*StockMovement, error)
	FindAllByProduct(ctx context.Context, productID string, limit, offset int32) ([]*StockMovement, error)
	FindAllByLocation(ctx context.Context, locationID string, limit, offset int32) ([]*StockMovement, error)
}
