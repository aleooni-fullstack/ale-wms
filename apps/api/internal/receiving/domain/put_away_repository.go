package domain

import (
	"context"
)

type PutAwayRepository interface {
	Create(ctx context.Context, p *PutAway) error
	FindByID(ctx context.Context, id string) (*PutAway, error)
	FindByReceiptID(ctx context.Context, receiptID string) (*PutAway, error)
	UpdateStatus(ctx context.Context, p *PutAway) error
	AddItem(ctx context.Context, item *PutAwayItem) error
	FindAllItems(ctx context.Context, putAwayID string) ([]*PutAwayItem, error)
	UpdateItemPutAway(ctx context.Context, item *PutAwayItem) error
}
