package domain

import (
	"context"
)

type StockTransferRepository interface {
	Create(ctx context.Context, t *StockTransfer) error
	FindByID(ctx context.Context, id string) (*StockTransfer, error)
	FindAllByProduct(ctx context.Context, productID string, limit, offset int32) ([]*StockTransfer, error)
	FindAllByFromLocation(ctx context.Context, fromLocationID string, limit, offset int32) ([]*StockTransfer, error)
	FindAllByToLocation(ctx context.Context, toLocationID string, limit, offset int32) ([]*StockTransfer, error)
	UpdateStatus(ctx context.Context, t *StockTransfer) error
}
