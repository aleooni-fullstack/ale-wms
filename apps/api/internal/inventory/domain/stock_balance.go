package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	"github.com/nrednav/cuid2"
)

type StockBalance struct {
	domain.AggregateRoot[string]
	ProductID        string
	LocationID       string
	Quantity         float64
	ReservedQuantity float64
	UpdatedAt        time.Time
}

func NewStockBalance(productID, locationID string, quantity float64) *StockBalance {
	return &StockBalance{
		AggregateRoot:    domain.NewAggregateRoot[string](cuid2.Generate()),
		ProductID:        productID,
		LocationID:       locationID,
		Quantity:         quantity,
		ReservedQuantity: 0,
	}
}

func RestoreStockBalance(id, productID, locationID string, quantity, reservedQuantity float64, updatedAt time.Time) *StockBalance {
	return &StockBalance{
		AggregateRoot:    domain.NewAggregateRoot[string](id),
		ProductID:        productID,
		LocationID:       locationID,
		Quantity:         quantity,
		ReservedQuantity: reservedQuantity,
		UpdatedAt:        updatedAt,
	}
}

func (b *StockBalance) AvailableQuantity() float64 {
	return b.Quantity - b.ReservedQuantity
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

func (b *StockBalance) Reserve(quantity float64) {
	b.ReservedQuantity += quantity
}

func (b *StockBalance) Release(quantity float64) {
	b.ReservedQuantity -= quantity
	if b.ReservedQuantity < 0 {
		b.ReservedQuantity = 0
	}
}
