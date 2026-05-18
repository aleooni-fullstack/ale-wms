package domain

import (
	"context"
)

type PickingRepository interface {
	Create(ctx context.Context, p *Picking) error
	FindByID(ctx context.Context, id string) (*Picking, error)
	FindByOrderID(ctx context.Context, orderID string) (*Picking, error)
	UpdateStatus(ctx context.Context, p *Picking) error
	AddItem(ctx context.Context, item *PickingItem) error
	FindAllItems(ctx context.Context, pickingID string) ([]*PickingItem, error)
	UpdateItemPicked(ctx context.Context, item *PickingItem) error
}
