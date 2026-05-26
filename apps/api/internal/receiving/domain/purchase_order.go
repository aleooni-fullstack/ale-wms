package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/nrednav/cuid2"
)

type PurchaseOrderStatus string

const (
	PurchaseOrderStatusDraft     PurchaseOrderStatus = "DRAFT"
	PurchaseOrderStatusConfirmed PurchaseOrderStatus = "CONFIRMED"
	PurchaseOrderStatusReceiving PurchaseOrderStatus = "RECEIVING"
	PurchaseOrderStatusCompleted PurchaseOrderStatus = "COMPLETED"
	PurchaseOrderStatusCancelled PurchaseOrderStatus = "CANCELLED"
)

type PurchaseOrderItem struct {
	domain.AggregateRoot[string]
	PurchaseOrderID string
	ProductID       string
	Quantity        float64
	CreatedAt       time.Time
}

type PurchaseOrder struct {
	domain.AggregateRoot[string]
	Reference string
	Supplier  string
	Status    PurchaseOrderStatus
	Note      string
	Items     []*PurchaseOrderItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPurchaseOrder(reference, supplier, note string) (*PurchaseOrder, error) {
	if reference == "" {
		return nil, dderr.New("INVALID_REFERENCE", "reference is required", nil)
	}

	return &PurchaseOrder{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		Reference:     reference,
		Supplier:      supplier,
		Status:        PurchaseOrderStatusDraft,
		Note:          note,
		Items:         []*PurchaseOrderItem{},
	}, nil
}

func RestorePurchaseOrder(id, reference, supplier string, status PurchaseOrderStatus, note string, createdAt, updatedAt time.Time) *PurchaseOrder {
	return &PurchaseOrder{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		Reference:     reference,
		Supplier:      supplier,
		Status:        status,
		Note:          note,
		Items:         []*PurchaseOrderItem{},
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func NewPurchaseOrderItem(purchaseOrderID, productID string, quantity float64) (*PurchaseOrderItem, error) {
	if purchaseOrderID == "" {
		return nil, dderr.New("INVALID_PURCHASE_ORDER_ID", "purchase_order_id is required", nil)
	}
	if productID == "" {
		return nil, dderr.New("INVALID_PRODUCT_ID", "product_id is required", nil)
	}
	if quantity <= 0 {
		return nil, dderr.New("INVALID_QUANTITY", "quantity must be greater than zero", nil)
	}

	return &PurchaseOrderItem{
		AggregateRoot:   domain.NewAggregateRoot[string](cuid2.Generate()),
		PurchaseOrderID: purchaseOrderID,
		ProductID:       productID,
		Quantity:        quantity,
	}, nil
}

func RestorePurchaseOrderItem(id, purchaseOrderID, productID string, quantity float64, createdAt time.Time) *PurchaseOrderItem {
	return &PurchaseOrderItem{
		AggregateRoot:   domain.NewAggregateRoot[string](id),
		PurchaseOrderID: purchaseOrderID,
		ProductID:       productID,
		Quantity:        quantity,
		CreatedAt:       createdAt,
	}
}

func (po *PurchaseOrder) Confirm() error {
	if po.Status != PurchaseOrderStatusDraft {
		return dderr.New("INVALID_STATUS", "only draft purchase orders can be confirmed", nil)
	}
	if len(po.Items) == 0 {
		return dderr.New("INVALID_PURCHASE_ORDER", "purchase order must have at least one item", nil)
	}

	po.Status = PurchaseOrderStatusConfirmed

	return nil
}

func (po *PurchaseOrder) StartReceiving() error {
	if po.Status != PurchaseOrderStatusConfirmed {
		return dderr.New("INVALID_STATUS", "only confirmed purchase orders can start receiving", nil)
	}

	po.Status = PurchaseOrderStatusReceiving

	return nil
}

func (po *PurchaseOrder) Complete() error {
	if po.Status != PurchaseOrderStatusReceiving {
		return dderr.New("INVALID_STATUS", "only receiving purchase orders can be completed", nil)
	}

	po.Status = PurchaseOrderStatusCompleted

	return nil
}

func (po *PurchaseOrder) Cancel() error {
	if po.Status == PurchaseOrderStatusCompleted {
		return dderr.New("INVALID_STATUS", "completed purchase orders cannot be cancelled", nil)
	}

	po.Status = PurchaseOrderStatusCancelled

	return nil
}
