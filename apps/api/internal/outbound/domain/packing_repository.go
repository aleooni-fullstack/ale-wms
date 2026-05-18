package domain

import (
	"context"
)

type PackingRepository interface {
	Create(ctx context.Context, p *Packing) error
	FindByID(ctx context.Context, id string) (*Packing, error)
	FindByOrderID(ctx context.Context, orderID string) (*Packing, error)
	UpdateStatus(ctx context.Context, p *Packing) error
}
