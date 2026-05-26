package domain

import (
	"context"
)

type PurchaseOrderRepository interface {
	Create(ctx context.Context, po *PurchaseOrder) error
	FindByID(ctx context.Context, id string) (*PurchaseOrder, error)
	FindByReference(ctx context.Context, reference string) (*PurchaseOrder, error)
	FindAll(ctx context.Context, limit, offset int32) ([]*PurchaseOrder, error)
	UpdateStatus(ctx context.Context, po *PurchaseOrder) error
	AddItem(ctx context.Context, item *PurchaseOrderItem) error
	FindAllItems(ctx context.Context, purchaseOrderID string) ([]*PurchaseOrderItem, error)
}
