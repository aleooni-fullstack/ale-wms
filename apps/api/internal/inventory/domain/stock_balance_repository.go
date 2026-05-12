package domain

import (
	"context"
)

type StockBalanceRepository interface {
	FindByProductAndLocation(ctx context.Context, productID, locationID string) (*StockBalance, error)
	FindAllByProduct(ctx context.Context, productID string) ([]*StockBalance, error)
	FindAllByLocation(ctx context.Context, locationID string) ([]*StockBalance, error)
	Upsert(ctx context.Context, b *StockBalance) error
}
