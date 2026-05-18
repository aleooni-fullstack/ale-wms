package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/nrednav/cuid2"
)

type OrderStatus string

const (
	OrderStatusDraft     OrderStatus = "DRAFT"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusPicking   OrderStatus = "PICKING"
	OrderStatusPacking   OrderStatus = "PACKING"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type OrderItem struct {
	domain.AggregateRoot[string]
	OrderID    string
	ProductID  string
	LocationID string
	Quantity   float64
	CreatedAt  time.Time
}

type Order struct {
	domain.AggregateRoot[string]
	Reference string
	Status    OrderStatus
	Note      string
	Items     []*OrderItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOrder(reference, note string) (*Order, error) {
	if reference == "" {
		return nil, dderr.New("INVALID_REFERENCE", "reference is required", nil)
	}

	return &Order{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		Reference:     reference,
		Status:        OrderStatusDraft,
		Note:          note,
		Items:         []*OrderItem{},
	}, nil
}

func RestoreOrder(id, reference string, status OrderStatus, note string, createdAt, updatedAt time.Time) *Order {
	return &Order{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		Reference:     reference,
		Status:        status,
		Note:          note,
		Items:         []*OrderItem{},
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func NewOrderItem(orderID, productID, locationID string, quantity float64) (*OrderItem, error) {
	if orderID == "" {
		return nil, dderr.New("INVALID_ORDER_ID", "order_id is required", nil)
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

	return &OrderItem{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		OrderID:       orderID,
		ProductID:     productID,
		LocationID:    locationID,
		Quantity:      quantity,
	}, nil
}

func RestoreOrderItem(id, orderID, productID, locationID string, quantity float64, createdAt time.Time) *OrderItem {
	return &OrderItem{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		OrderID:       orderID,
		ProductID:     productID,
		LocationID:    locationID,
		Quantity:      quantity,
		CreatedAt:     createdAt,
	}
}

func (o *Order) Confirm() error {
	if o.Status != OrderStatusDraft {
		return dderr.New("INVALID_STATUS", "only draft orders can be confirmed", nil)
	}
	if len(o.Items) == 0 {
		return dderr.New("INVALID_ORDER", "order must have at least one item", nil)
	}

	o.Status = OrderStatusConfirmed

	return nil
}

func (o *Order) StartPicking() error {
	if o.Status != OrderStatusConfirmed {
		return dderr.New("INVALID_STATUS", "only confirmed orders can start picking", nil)
	}

	o.Status = OrderStatusPicking

	return nil
}

func (o *Order) StartPacking() error {
	if o.Status != OrderStatusPicking {
		return dderr.New("INVALID_STATUS", "only picking orders can start packing", nil)
	}

	o.Status = OrderStatusPacking

	return nil
}

func (o *Order) Ship() error {
	if o.Status != OrderStatusPacking {
		return dderr.New("INVALID_STATUS", "only packing orders can be shipped", nil)
	}

	o.Status = OrderStatusShipped

	return nil
}

func (o *Order) Cancel() error {
	if o.Status == OrderStatusShipped {
		return dderr.New("INVALID_STATUS", "shipped orders cannot be cancelled", nil)
	}

	o.Status = OrderStatusCancelled

	return nil
}
