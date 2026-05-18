package domain

import (
	"context"
)

type ShippingRepository interface {
	Create(ctx context.Context, s *Shipping) error
	FindByID(ctx context.Context, id string) (*Shipping, error)
	FindByOrderID(ctx context.Context, orderID string) (*Shipping, error)
	UpdateStatus(ctx context.Context, s *Shipping) error
}
