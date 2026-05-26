package domain

import (
	"context"
)

type ReceiptRepository interface {
	Create(ctx context.Context, r *Receipt) error
	FindByID(ctx context.Context, id string) (*Receipt, error)
	FindByPurchaseOrderID(ctx context.Context, purchaseOrderID string) (*Receipt, error)
	UpdateStatus(ctx context.Context, r *Receipt) error
	AddItem(ctx context.Context, item *ReceiptItem) error
	FindAllItems(ctx context.Context, receiptID string) ([]*ReceiptItem, error)
	UpdateItemReceivedQuantity(ctx context.Context, item *ReceiptItem) error
}
