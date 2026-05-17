package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/nrednav/cuid2"
)

type InventoryBalanceStatus string

const (
	InventoryBalanceStatusDraft      InventoryBalanceStatus = "DRAFT"
	InventoryBalanceStatusInProgress InventoryBalanceStatus = "IN_PROGRESS"
	InventoryBalanceStatusCompleted  InventoryBalanceStatus = "COMPLETED"
	InventoryBalanceStatusCancelled  InventoryBalanceStatus = "CANCELLED"
)

type InventoryBalanceItem struct {
	domain.AggregateRoot[string]
	InventoryBalanceID string
	ProductID          string
	SystemQuantity     float64
	CountedQuantity    *float64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type InventoryBalance struct {
	domain.AggregateRoot[string]
	LocationID string
	Status     InventoryBalanceStatus
	Note       string
	Items      []*InventoryBalanceItem
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewInventoryBalance(locationID, note string) (*InventoryBalance, error) {
	if locationID == "" {
		return nil, dderr.New("INVALID_LOCATION_ID", "location_id is required", nil)
	}

	return &InventoryBalance{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		LocationID:    locationID,
		Status:        InventoryBalanceStatusDraft,
		Note:          note,
		Items:         []*InventoryBalanceItem{},
	}, nil
}

func RestoreInventoryBalance(id, locationID string, status InventoryBalanceStatus, note string, createdAt, updatedAt time.Time) *InventoryBalance {
	return &InventoryBalance{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		LocationID:    locationID,
		Status:        status,
		Note:          note,
		Items:         []*InventoryBalanceItem{},
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func NewInventoryBalanceItem(inventoryBalanceID, productID string, systemQuantity float64) (*InventoryBalanceItem, error) {
	if inventoryBalanceID == "" {
		return nil, dderr.New("INVALID_INVENTORY_BALANCE_ID", "inventory_balance_id is required", nil)
	}
	if productID == "" {
		return nil, dderr.New("INVALID_PRODUCT_ID", "product_id is required", nil)
	}
	if systemQuantity < 0 {
		return nil, dderr.New("INVALID_QUANTITY", "system_quantity must be greater than or equal to zero", nil)
	}

	return &InventoryBalanceItem{
		AggregateRoot:      domain.NewAggregateRoot[string](cuid2.Generate()),
		InventoryBalanceID: inventoryBalanceID,
		ProductID:          productID,
		SystemQuantity:     systemQuantity,
	}, nil
}

func RestoreInventoryBalanceItem(id, inventoryBalanceID, productID string, systemQuantity float64, countedQuantity *float64, createdAt, updatedAt time.Time) *InventoryBalanceItem {
	return &InventoryBalanceItem{
		AggregateRoot:      domain.NewAggregateRoot[string](id),
		InventoryBalanceID: inventoryBalanceID,
		ProductID:          productID,
		SystemQuantity:     systemQuantity,
		CountedQuantity:    countedQuantity,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}
}

func (b *InventoryBalance) Start() error {
	if b.Status != InventoryBalanceStatusDraft {
		return dderr.New("INVALID_STATUS", "only draft balances can be started", nil)
	}

	b.Status = InventoryBalanceStatusInProgress

	return nil
}

func (b *InventoryBalance) Complete() error {
	if b.Status != InventoryBalanceStatusInProgress {
		return dderr.New("INVALID_STATUS", "only in_progress balances can be completed", nil)
	}

	b.Status = InventoryBalanceStatusCompleted

	return nil
}

func (b *InventoryBalance) Cancel() error {
	if b.Status != InventoryBalanceStatusDraft && b.Status != InventoryBalanceStatusInProgress {
		return dderr.New("INVALID_STATUS", "only draft or in_progress balances can be cancelled", nil)
	}

	b.Status = InventoryBalanceStatusCancelled

	return nil
}

func (item *InventoryBalanceItem) Count(quantity float64) error {
	if quantity < 0 {
		return dderr.New("INVALID_QUANTITY", "counted_quantity must be greater than or equal to zero", nil)
	}

	item.CountedQuantity = &quantity

	return nil
}

func (item *InventoryBalanceItem) Difference() *float64 {
	if item.CountedQuantity == nil {
		return nil
	}

	diff := *item.CountedQuantity - item.SystemQuantity

	return &diff
}
