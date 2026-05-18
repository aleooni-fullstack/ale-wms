package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/nrednav/cuid2"
)

type ShippingStatus string

const (
	ShippingStatusPending   ShippingStatus = "PENDING"
	ShippingStatusShipped   ShippingStatus = "SHIPPED"
	ShippingStatusCancelled ShippingStatus = "CANCELLED"
)

type Shipping struct {
	domain.AggregateRoot[string]
	OrderID      string
	PackingID    string
	Status       ShippingStatus
	TrackingCode string
	Note         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewShipping(orderID, packingID, note string) (*Shipping, error) {
	if orderID == "" {
		return nil, dderr.New("INVALID_ORDER_ID", "order_id is required", nil)
	}
	if packingID == "" {
		return nil, dderr.New("INVALID_PACKING_ID", "packing_id is required", nil)
	}

	return &Shipping{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		OrderID:       orderID,
		PackingID:     packingID,
		Status:        ShippingStatusPending,
		Note:          note,
	}, nil
}

func RestoreShipping(id, orderID, packingID string, status ShippingStatus, trackingCode, note string, createdAt, updatedAt time.Time) *Shipping {
	return &Shipping{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		OrderID:       orderID,
		PackingID:     packingID,
		Status:        status,
		TrackingCode:  trackingCode,
		Note:          note,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (s *Shipping) Ship(trackingCode string) error {
	if s.Status != ShippingStatusPending {
		return dderr.New("INVALID_STATUS", "only pending shippings can be shipped", nil)
	}

	s.Status = ShippingStatusShipped
	s.TrackingCode = trackingCode

	return nil
}

func (s *Shipping) Cancel() error {
	if s.Status == ShippingStatusShipped {
		return dderr.New("INVALID_STATUS", "shipped shippings cannot be cancelled", nil)
	}

	s.Status = ShippingStatusCancelled

	return nil
}
