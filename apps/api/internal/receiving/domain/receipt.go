package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/nrednav/cuid2"
)

type ReceiptStatus string

const (
	ReceiptStatusPending    ReceiptStatus = "PENDING"
	ReceiptStatusInProgress ReceiptStatus = "IN_PROGRESS"
	ReceiptStatusCompleted  ReceiptStatus = "COMPLETED"
	ReceiptStatusCancelled  ReceiptStatus = "CANCELLED"
)

type ReceiptItem struct {
	domain.AggregateRoot[string]
	ReceiptID        string
	ProductID        string
	ExpectedQuantity float64
	ReceivedQuantity *float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Receipt struct {
	domain.AggregateRoot[string]
	PurchaseOrderID string
	Status          ReceiptStatus
	Note            string
	Items           []*ReceiptItem
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewReceipt(purchaseOrderID, note string) (*Receipt, error) {
	if purchaseOrderID == "" {
		return nil, dderr.New("INVALID_PURCHASE_ORDER_ID", "purchase_order_id is required", nil)
	}

	return &Receipt{
		AggregateRoot:   domain.NewAggregateRoot[string](cuid2.Generate()),
		PurchaseOrderID: purchaseOrderID,
		Status:          ReceiptStatusPending,
		Note:            note,
		Items:           []*ReceiptItem{},
	}, nil
}

func RestoreReceipt(id, purchaseOrderID string, status ReceiptStatus, note string, createdAt, updatedAt time.Time) *Receipt {
	return &Receipt{
		AggregateRoot:   domain.NewAggregateRoot[string](id),
		PurchaseOrderID: purchaseOrderID,
		Status:          status,
		Note:            note,
		Items:           []*ReceiptItem{},
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

func NewReceiptItem(receiptID, productID string, expectedQuantity float64) (*ReceiptItem, error) {
	if receiptID == "" {
		return nil, dderr.New("INVALID_RECEIPT_ID", "receipt_id is required", nil)
	}
	if productID == "" {
		return nil, dderr.New("INVALID_PRODUCT_ID", "product_id is required", nil)
	}
	if expectedQuantity <= 0 {
		return nil, dderr.New("INVALID_QUANTITY", "expected_quantity must be greater than zero", nil)
	}

	return &ReceiptItem{
		AggregateRoot:    domain.NewAggregateRoot[string](cuid2.Generate()),
		ReceiptID:        receiptID,
		ProductID:        productID,
		ExpectedQuantity: expectedQuantity,
	}, nil
}

func RestoreReceiptItem(id, receiptID, productID string, expectedQuantity float64, receivedQuantity *float64, createdAt, updatedAt time.Time) *ReceiptItem {
	return &ReceiptItem{
		AggregateRoot:    domain.NewAggregateRoot[string](id),
		ReceiptID:        receiptID,
		ProductID:        productID,
		ExpectedQuantity: expectedQuantity,
		ReceivedQuantity: receivedQuantity,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}

func (r *Receipt) Start() error {
	if r.Status != ReceiptStatusPending {
		return dderr.New("INVALID_STATUS", "only pending receipts can be started", nil)
	}

	r.Status = ReceiptStatusInProgress

	return nil
}

func (r *Receipt) Complete() error {
	if r.Status != ReceiptStatusInProgress {
		return dderr.New("INVALID_STATUS", "only in_progress receipts can be completed", nil)
	}

	r.Status = ReceiptStatusCompleted

	return nil
}

func (r *Receipt) Cancel() error {
	if r.Status == ReceiptStatusCompleted {
		return dderr.New("INVALID_STATUS", "completed receipts cannot be cancelled", nil)
	}

	r.Status = ReceiptStatusCancelled

	return nil
}

func (item *ReceiptItem) Receive(quantity float64) error {
	if quantity <= 0 {
		return dderr.New("INVALID_QUANTITY", "received_quantity must be greater than zero", nil)
	}

	item.ReceivedQuantity = &quantity

	return nil
}

func (item *ReceiptItem) Difference() *float64 {
	if item.ReceivedQuantity == nil {
		return nil
	}

	diff := *item.ReceivedQuantity - item.ExpectedQuantity

	return &diff
}
