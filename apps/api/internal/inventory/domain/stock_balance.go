package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	"github.com/nrednav/cuid2"
)

type StockBalance struct {
	domain.AggregateRoot[string]
	ProductID  string
	LocationID string
	Quantity   float64
	UpdatedAt  time.Time
}

func NewStockBalance(productID, locationID string, quantity float64) *StockBalance {
	return &StockBalance{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		ProductID:     productID,
		LocationID:    locationID,
		Quantity:      quantity,
	}
}

func RestoreStockBalance(id, productID, locationID string, quantity float64, updatedAt time.Time) *StockBalance {
	return &StockBalance{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		ProductID:     productID,
		LocationID:    locationID,
		Quantity:      quantity,
		UpdatedAt:     updatedAt,
	}
}

func (b *StockBalance) Apply(m *StockMovement) {
	switch m.Type {
	case MovementTypeIn:
		b.Quantity += m.Quantity
	case MovementTypeOut:
		b.Quantity -= m.Quantity
	case MovementTypeAdjustment:
		b.Quantity = m.Quantity
	}
}
