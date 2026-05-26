package application

import (
	"context"

	dderr "github.com/aleodoni/go-ddd/errors"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/domain"
)

type PurchaseOrderService struct {
	repo domain.PurchaseOrderRepository
}

func NewPurchaseOrderService(repo domain.PurchaseOrderRepository) *PurchaseOrderService {
	return &PurchaseOrderService{repo: repo}
}

type CreatePurchaseOrderInput struct {
	Reference string
	Supplier  string
	Note      string
}

type AddPurchaseOrderItemInput struct {
	PurchaseOrderID string
	ProductID       string
	Quantity        float64
}

type ListPurchaseOrdersInput struct {
	Page    int32
	PerPage int32
}

type ListPurchaseOrdersOutput struct {
	Data    []*domain.PurchaseOrder
	Page    int32
	PerPage int32
}

func (s *PurchaseOrderService) Create(ctx context.Context, input CreatePurchaseOrderInput) (*domain.PurchaseOrder, error) {
	_, err := s.repo.FindByReference(ctx, input.Reference)
	if err == nil {
		return nil, dderr.New("REFERENCE_ALREADY_EXISTS", "a purchase order with this reference already exists", nil)
	}
	if err != dderr.ErrNotFound {
		return nil, err
	}

	po, err := domain.NewPurchaseOrder(input.Reference, input.Supplier, input.Note)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, po); err != nil {
		return nil, err
	}

	return po, nil
}

func (s *PurchaseOrderService) GetByID(ctx context.Context, id string) (*domain.PurchaseOrder, error) {
	po, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.FindAllItems(ctx, id)
	if err != nil {
		return nil, err
	}

	po.Items = items

	return po, nil
}

func (s *PurchaseOrderService) List(ctx context.Context, input ListPurchaseOrdersInput) (*ListPurchaseOrdersOutput, error) {
	if input.PerPage == 0 {
		input.PerPage = 20
	}
	if input.Page == 0 {
		input.Page = 1
	}

	offset := (input.Page - 1) * input.PerPage

	pos, err := s.repo.FindAll(ctx, input.PerPage, offset)
	if err != nil {
		return nil, err
	}

	return &ListPurchaseOrdersOutput{
		Data:    pos,
		Page:    input.Page,
		PerPage: input.PerPage,
	}, nil
}

func (s *PurchaseOrderService) AddItem(ctx context.Context, input AddPurchaseOrderItemInput) (*domain.PurchaseOrderItem, error) {
	po, err := s.repo.FindByID(ctx, input.PurchaseOrderID)
	if err != nil {
		return nil, err
	}

	if po.Status != domain.PurchaseOrderStatusDraft {
		return nil, dderr.New("INVALID_STATUS", "items can only be added to draft purchase orders", nil)
	}

	item, err := domain.NewPurchaseOrderItem(input.PurchaseOrderID, input.ProductID, input.Quantity)
	if err != nil {
		return nil, err
	}

	if err := s.repo.AddItem(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *PurchaseOrderService) Confirm(ctx context.Context, id string) (*domain.PurchaseOrder, error) {
	po, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.FindAllItems(ctx, id)
	if err != nil {
		return nil, err
	}

	po.Items = items

	if err := po.Confirm(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateStatus(ctx, po); err != nil {
		return nil, err
	}

	return po, nil
}

func (s *PurchaseOrderService) Cancel(ctx context.Context, id string) (*domain.PurchaseOrder, error) {
	po, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := po.Cancel(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateStatus(ctx, po); err != nil {
		return nil, err
	}

	return po, nil
}
