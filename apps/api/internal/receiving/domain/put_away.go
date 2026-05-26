package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/nrednav/cuid2"
)

type PutAwayStatus string

const (
	PutAwayStatusPending    PutAwayStatus = "PENDING"
	PutAwayStatusInProgress PutAwayStatus = "IN_PROGRESS"
	PutAwayStatusCompleted  PutAwayStatus = "COMPLETED"
	PutAwayStatusCancelled  PutAwayStatus = "CANCELLED"
)

type PutAwayItem struct {
	domain.AggregateRoot[string]
	PutAwayID  string
	ProductID  string
	LocationID string
	Quantity   float64
	PutAway    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type PutAway struct {
	domain.AggregateRoot[string]
	ReceiptID string
	Status    PutAwayStatus
	Note      string
	Items     []*PutAwayItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPutAway(receiptID, note string) (*PutAway, error) {
	if receiptID == "" {
		return nil, dderr.New("INVALID_RECEIPT_ID", "receipt_id is required", nil)
	}

	return &PutAway{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		ReceiptID:     receiptID,
		Status:        PutAwayStatusPending,
		Note:          note,
		Items:         []*PutAwayItem{},
	}, nil
}

func RestorePutAway(id, receiptID string, status PutAwayStatus, note string, createdAt, updatedAt time.Time) *PutAway {
	return &PutAway{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		ReceiptID:     receiptID,
		Status:        status,
		Note:          note,
		Items:         []*PutAwayItem{},
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func NewPutAwayItem(putAwayID, productID, locationID string, quantity float64) (*PutAwayItem, error) {
	if putAwayID == "" {
		return nil, dderr.New("INVALID_PUT_AWAY_ID", "put_away_id is required", nil)
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

	return &PutAwayItem{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		PutAwayID:     putAwayID,
		ProductID:     productID,
		LocationID:    locationID,
		Quantity:      quantity,
		PutAway:       false,
	}, nil
}

func RestorePutAwayItem(id, putAwayID, productID, locationID string, quantity float64, putAway bool, createdAt, updatedAt time.Time) *PutAwayItem {
	return &PutAwayItem{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		PutAwayID:     putAwayID,
		ProductID:     productID,
		LocationID:    locationID,
		Quantity:      quantity,
		PutAway:       putAway,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (p *PutAway) Start() error {
	if p.Status != PutAwayStatusPending {
		return dderr.New("INVALID_STATUS", "only pending put aways can be started", nil)
	}

	p.Status = PutAwayStatusInProgress

	return nil
}

func (p *PutAway) Complete() error {
	if p.Status != PutAwayStatusInProgress {
		return dderr.New("INVALID_STATUS", "only in_progress put aways can be completed", nil)
	}

	for _, item := range p.Items {
		if !item.PutAway {
			return dderr.New("INCOMPLETE_PUT_AWAY", "all items must be put away before completing", nil)
		}
	}

	p.Status = PutAwayStatusCompleted

	return nil
}

func (p *PutAway) Cancel() error {
	if p.Status == PutAwayStatusCompleted {
		return dderr.New("INVALID_STATUS", "completed put aways cannot be cancelled", nil)
	}

	p.Status = PutAwayStatusCancelled

	return nil
}

func (item *PutAwayItem) Store() {
	item.PutAway = true
}
