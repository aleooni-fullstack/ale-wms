package application

import (
	"context"

	dderr "github.com/aleodoni/go-ddd/errors"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/domain"
)

type OrderService struct {
	repo domain.OrderRepository
}

func NewOrderService(repo domain.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

type CreateOrderInput struct {
	Reference string
	Note      string
}

type AddOrderItemInput struct {
	OrderID    string
	ProductID  string
	LocationID string
	Quantity   float64
}

type ListOrdersInput struct {
	Page    int32
	PerPage int32
}

type ListOrdersOutput struct {
	Data    []*domain.Order
	Page    int32
	PerPage int32
}

func (s *OrderService) Create(ctx context.Context, input CreateOrderInput) (*domain.Order, error) {
	_, err := s.repo.FindByReference(ctx, input.Reference)
	if err == nil {
		return nil, dderr.New("REFERENCE_ALREADY_EXISTS", "an order with this reference already exists", nil)
	}
	if err != dderr.ErrNotFound {
		return nil, err
	}

	order, err := domain.NewOrder(input.Reference, input.Note)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.FindAllItems(ctx, id)
	if err != nil {
		return nil, err
	}

	order.Items = items

	return order, nil
}

func (s *OrderService) List(ctx context.Context, input ListOrdersInput) (*ListOrdersOutput, error) {
	if input.PerPage == 0 {
		input.PerPage = 20
	}
	if input.Page == 0 {
		input.Page = 1
	}

	offset := (input.Page - 1) * input.PerPage

	orders, err := s.repo.FindAll(ctx, input.PerPage, offset)
	if err != nil {
		return nil, err
	}

	return &ListOrdersOutput{
		Data:    orders,
		Page:    input.Page,
		PerPage: input.PerPage,
	}, nil
}

func (s *OrderService) AddItem(ctx context.Context, input AddOrderItemInput) (*domain.OrderItem, error) {
	order, err := s.repo.FindByID(ctx, input.OrderID)
	if err != nil {
		return nil, err
	}

	if order.Status != domain.OrderStatusDraft {
		return nil, dderr.New("INVALID_STATUS", "items can only be added to draft orders", nil)
	}

	item, err := domain.NewOrderItem(input.OrderID, input.ProductID, input.LocationID, input.Quantity)
	if err != nil {
		return nil, err
	}

	if err := s.repo.AddItem(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *OrderService) Confirm(ctx context.Context, id string) (*domain.Order, error) {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.FindAllItems(ctx, id)
	if err != nil {
		return nil, err
	}

	order.Items = items

	if err := order.Confirm(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateStatus(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) Cancel(ctx context.Context, id string) (*domain.Order, error) {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := order.Cancel(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateStatus(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}
