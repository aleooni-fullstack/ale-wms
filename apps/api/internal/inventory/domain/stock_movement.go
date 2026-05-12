package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/nrednav/cuid2"
)

type MovementType string

const (
	MovementTypeIn         MovementType = "IN"
	MovementTypeOut        MovementType = "OUT"
	MovementTypeAdjustment MovementType = "ADJUSTMENT"
)

type StockMovement struct {
	domain.AggregateRoot[string]
	ProductID  string
	LocationID string
	Type       MovementType
	Quantity   float64
	Note       string
	CreatedAt  time.Time
}

func NewStockMovement(productID, locationID string, movementType MovementType, quantity float64, note string) (*StockMovement, error) {
	if productID == "" {
		return nil, dderr.New("INVALID_PRODUCT_ID", "product_id is required", nil)
	}
	if locationID == "" {
		return nil, dderr.New("INVALID_LOCATION_ID", "location_id is required", nil)
	}
	if quantity <= 0 {
		return nil, dderr.New("INVALID_QUANTITY", "quantity must be greater than zero", nil)
	}
	if movementType != MovementTypeIn && movementType != MovementTypeOut && movementType != MovementTypeAdjustment {
		return nil, dderr.New("INVALID_MOVEMENT_TYPE", "movement type must be IN, OUT or ADJUSTMENT", nil)
	}

	return &StockMovement{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		ProductID:     productID,
		LocationID:    locationID,
		Type:          movementType,
		Quantity:      quantity,
		Note:          note,
	}, nil
}

func RestoreStockMovement(id, productID, locationID string, movementType MovementType, quantity float64, note string, createdAt time.Time) *StockMovement {
	return &StockMovement{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		ProductID:     productID,
		LocationID:    locationID,
		Type:          movementType,
		Quantity:      quantity,
		Note:          note,
		CreatedAt:     createdAt,
	}
}
