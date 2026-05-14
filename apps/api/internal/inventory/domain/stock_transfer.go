package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/nrednav/cuid2"
)

type TransferStatus string

const (
	TransferStatusPending   TransferStatus = "PENDING"
	TransferStatusCompleted TransferStatus = "COMPLETED"
	TransferStatusCancelled TransferStatus = "CANCELLED"
)

type StockTransfer struct {
	domain.AggregateRoot[string]
	ProductID      string
	FromLocationID string
	ToLocationID   string
	Quantity       float64
	Status         TransferStatus
	Note           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewStockTransfer(productID, fromLocationID, toLocationID string, quantity float64, note string) (*StockTransfer, error) {
	if productID == "" {
		return nil, dderr.New("INVALID_PRODUCT_ID", "product_id is required", nil)
	}
	if fromLocationID == "" {
		return nil, dderr.New("INVALID_FROM_LOCATION_ID", "from_location_id is required", nil)
	}
	if toLocationID == "" {
		return nil, dderr.New("INVALID_TO_LOCATION_ID", "to_location_id is required", nil)
	}
	if fromLocationID == toLocationID {
		return nil, dderr.New("INVALID_TRANSFER", "from_location_id and to_location_id must be different", nil)
	}
	if quantity <= 0 {
		return nil, dderr.New("INVALID_QUANTITY", "quantity must be greater than zero", nil)
	}

	return &StockTransfer{
		AggregateRoot:  domain.NewAggregateRoot[string](cuid2.Generate()),
		ProductID:      productID,
		FromLocationID: fromLocationID,
		ToLocationID:   toLocationID,
		Quantity:       quantity,
		Status:         TransferStatusPending,
		Note:           note,
	}, nil
}

func RestoreStockTransfer(id, productID, fromLocationID, toLocationID string, quantity float64, status TransferStatus, note string, createdAt, updatedAt time.Time) *StockTransfer {
	return &StockTransfer{
		AggregateRoot:  domain.NewAggregateRoot[string](id),
		ProductID:      productID,
		FromLocationID: fromLocationID,
		ToLocationID:   toLocationID,
		Quantity:       quantity,
		Status:         status,
		Note:           note,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

func (t *StockTransfer) Complete() error {
	if t.Status != TransferStatusPending {
		return dderr.New("INVALID_STATUS", "only pending transfers can be completed", nil)
	}

	t.Status = TransferStatusCompleted

	return nil
}

func (t *StockTransfer) Cancel() error {
	if t.Status != TransferStatusPending {
		return dderr.New("INVALID_STATUS", "only pending transfers can be cancelled", nil)
	}

	t.Status = TransferStatusCancelled

	return nil
}
