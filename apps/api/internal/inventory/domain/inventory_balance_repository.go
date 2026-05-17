package domain

import (
	"context"
)

type InventoryBalanceRepository interface {
	Create(ctx context.Context, b *InventoryBalance) error
	FindByID(ctx context.Context, id string) (*InventoryBalance, error)
	FindAllByLocation(ctx context.Context, locationID string, limit, offset int32) ([]*InventoryBalance, error)
	UpdateStatus(ctx context.Context, b *InventoryBalance) error
	AddItem(ctx context.Context, item *InventoryBalanceItem) error
	FindItem(ctx context.Context, inventoryBalanceID, productID string) (*InventoryBalanceItem, error)
	FindAllItems(ctx context.Context, inventoryBalanceID string) ([]*InventoryBalanceItem, error)
	UpdateItemCountedQuantity(ctx context.Context, item *InventoryBalanceItem) error
}
