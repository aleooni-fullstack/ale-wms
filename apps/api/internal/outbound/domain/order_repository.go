package domain

import (
	"context"
)

type OrderRepository interface {
	Create(ctx context.Context, o *Order) error
	FindByID(ctx context.Context, id string) (*Order, error)
	FindByReference(ctx context.Context, reference string) (*Order, error)
	FindAll(ctx context.Context, limit, offset int32) ([]*Order, error)
	UpdateStatus(ctx context.Context, o *Order) error
	AddItem(ctx context.Context, item *OrderItem) error
	FindAllItems(ctx context.Context, orderID string) ([]*OrderItem, error)
}
