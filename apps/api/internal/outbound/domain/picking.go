package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/nrednav/cuid2"
)

type PickingStatus string

const (
	PickingStatusPending    PickingStatus = "PENDING"
	PickingStatusInProgress PickingStatus = "IN_PROGRESS"
	PickingStatusCompleted  PickingStatus = "COMPLETED"
	PickingStatusCancelled  PickingStatus = "CANCELLED"
)

type PickingItem struct {
	domain.AggregateRoot[string]
	PickingID  string
	ProductID  string
	LocationID string
	Quantity   float64
	Picked     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Picking struct {
	domain.AggregateRoot[string]
	OrderID   string
	Status    PickingStatus
	Note      string
	Items     []*PickingItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPicking(orderID, note string) (*Picking, error) {
	if orderID == "" {
		return nil, dderr.New("INVALID_ORDER_ID", "order_id is required", nil)
	}

	return &Picking{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		OrderID:       orderID,
		Status:        PickingStatusPending,
		Note:          note,
		Items:         []*PickingItem{},
	}, nil
}

func RestorePicking(id, orderID string, status PickingStatus, note string, createdAt, updatedAt time.Time) *Picking {
	return &Picking{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		OrderID:       orderID,
		Status:        status,
		Note:          note,
		Items:         []*PickingItem{},
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func NewPickingItem(pickingID, productID, locationID string, quantity float64) (*PickingItem, error) {
	if pickingID == "" {
		return nil, dderr.New("INVALID_PICKING_ID", "picking_id is required", nil)
	}
	if productID == "" {
		return nil, dderr.New("INVALID_PRODUCT_ID", "product_id is required", nil)
	}
	if locationID == "" {
		return nil, dderr.New("INVALID_LOCATION_ID", "location_id is required", nil)
	}
	if quantity <= 0 {
		return nil, dderr.New("INVALID_QUANTITY", "quantity must be greater than zero", nil)
	}

	return &PickingItem{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		PickingID:     pickingID,
		ProductID:     productID,
		LocationID:    locationID,
		Quantity:      quantity,
		Picked:        false,
	}, nil
}

func RestorePickingItem(id, pickingID, productID, locationID string, quantity float64, picked bool, createdAt, updatedAt time.Time) *PickingItem {
	return &PickingItem{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		PickingID:     pickingID,
		ProductID:     productID,
		LocationID:    locationID,
		Quantity:      quantity,
		Picked:        picked,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (p *Picking) Start() error {
	if p.Status != PickingStatusPending {
		return dderr.New("INVALID_STATUS", "only pending pickings can be started", nil)
	}

	p.Status = PickingStatusInProgress

	return nil
}

func (p *Picking) Complete() error {
	if p.Status != PickingStatusInProgress {
		return dderr.New("INVALID_STATUS", "only in_progress pickings can be completed", nil)
	}

	for _, item := range p.Items {
		if !item.Picked {
			return dderr.New("INCOMPLETE_PICKING", "all items must be picked before completing", nil)
		}
	}

	p.Status = PickingStatusCompleted

	return nil
}

func (p *Picking) Cancel() error {
	if p.Status == PickingStatusCompleted {
		return dderr.New("INVALID_STATUS", "completed pickings cannot be cancelled", nil)
	}

	p.Status = PickingStatusCancelled

	return nil
}

func (item *PickingItem) Pick() {
	item.Picked = true
}
