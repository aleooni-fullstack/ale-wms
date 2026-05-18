package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/nrednav/cuid2"
)

type PackingStatus string

const (
	PackingStatusPending    PackingStatus = "PENDING"
	PackingStatusInProgress PackingStatus = "IN_PROGRESS"
	PackingStatusCompleted  PackingStatus = "COMPLETED"
	PackingStatusCancelled  PackingStatus = "CANCELLED"
)

type Packing struct {
	domain.AggregateRoot[string]
	OrderID   string
	PickingID string
	Status    PackingStatus
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPacking(orderID, pickingID, note string) (*Packing, error) {
	if orderID == "" {
		return nil, dderr.New("INVALID_ORDER_ID", "order_id is required", nil)
	}
	if pickingID == "" {
		return nil, dderr.New("INVALID_PICKING_ID", "picking_id is required", nil)
	}

	return &Packing{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		OrderID:       orderID,
		PickingID:     pickingID,
		Status:        PackingStatusPending,
		Note:          note,
	}, nil
}

func RestorePacking(id, orderID, pickingID string, status PackingStatus, note string, createdAt, updatedAt time.Time) *Packing {
	return &Packing{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		OrderID:       orderID,
		PickingID:     pickingID,
		Status:        status,
		Note:          note,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (p *Packing) Start() error {
	if p.Status != PackingStatusPending {
		return dderr.New("INVALID_STATUS", "only pending packings can be started", nil)
	}

	p.Status = PackingStatusInProgress

	return nil
}

func (p *Packing) Complete() error {
	if p.Status != PackingStatusInProgress {
		return dderr.New("INVALID_STATUS", "only in_progress packings can be completed", nil)
	}

	p.Status = PackingStatusCompleted

	return nil
}

func (p *Packing) Cancel() error {
	if p.Status == PackingStatusCompleted {
		return dderr.New("INVALID_STATUS", "completed packings cannot be cancelled", nil)
	}

	p.Status = PackingStatusCancelled

	return nil
}
