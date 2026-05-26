package application

import (
	"context"

	dderr "github.com/aleodoni/go-ddd/errors"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/domain"
)

type ReceiptService struct {
	receiptRepo domain.ReceiptRepository
	poRepo      domain.PurchaseOrderRepository
}

func NewReceiptService(receiptRepo domain.ReceiptRepository, poRepo domain.PurchaseOrderRepository) *ReceiptService {
	return &ReceiptService{
		receiptRepo: receiptRepo,
		poRepo:      poRepo,
	}
}

type ReceiveItemInput struct {
	ReceiptID        string
	ProductID        string
	ReceivedQuantity float64
}

func (s *ReceiptService) Create(ctx context.Context, purchaseOrderID string) (*domain.Receipt, error) {
	po, err := s.poRepo.FindByID(ctx, purchaseOrderID)
	if err != nil {
		return nil, err
	}

	if po.Status != domain.PurchaseOrderStatusConfirmed {
		return nil, dderr.New("INVALID_STATUS", "only confirmed purchase orders can have a receipt", nil)
	}

	receipt, err := domain.NewReceipt(purchaseOrderID, "")
	if err != nil {
		return nil, err
	}

	if err := s.receiptRepo.Create(ctx, receipt); err != nil {
		return nil, err
	}

	items, err := s.poRepo.FindAllItems(ctx, purchaseOrderID)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		receiptItem, err := domain.NewReceiptItem(receipt.ID, item.ProductID, item.Quantity)
		if err != nil {
			return nil, err
		}

		if err := s.receiptRepo.AddItem(ctx, receiptItem); err != nil {
			return nil, err
		}

		receipt.Items = append(receipt.Items, receiptItem)
	}

	po.StartReceiving()
	if err := s.poRepo.UpdateStatus(ctx, po); err != nil {
		return nil, err
	}

	return receipt, nil
}

func (s *ReceiptService) GetByID(ctx context.Context, id string) (*domain.Receipt, error) {
	receipt, err := s.receiptRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.receiptRepo.FindAllItems(ctx, id)
	if err != nil {
		return nil, err
	}

	receipt.Items = items

	return receipt, nil
}

func (s *ReceiptService) Start(ctx context.Context, id string) (*domain.Receipt, error) {
	receipt, err := s.receiptRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := receipt.Start(); err != nil {
		return nil, err
	}

	if err := s.receiptRepo.UpdateStatus(ctx, receipt); err != nil {
		return nil, err
	}

	return receipt, nil
}

func (s *ReceiptService) ReceiveItem(ctx context.Context, input ReceiveItemInput) (*domain.ReceiptItem, error) {
	receipt, err := s.receiptRepo.FindByID(ctx, input.ReceiptID)
	if err != nil {
		return nil, err
	}

	if receipt.Status != domain.ReceiptStatusInProgress {
		return nil, dderr.New("INVALID_STATUS", "only in_progress receipts can receive items", nil)
	}

	items, err := s.receiptRepo.FindAllItems(ctx, input.ReceiptID)
	if err != nil {
		return nil, err
	}

	var target *domain.ReceiptItem
	for _, item := range items {
		if item.ProductID == input.ProductID {
			target = item
			break
		}
	}

	if target == nil {
		return nil, dderr.ErrNotFound
	}

	if err := target.Receive(input.ReceivedQuantity); err != nil {
		return nil, err
	}

	if err := s.receiptRepo.UpdateItemReceivedQuantity(ctx, target); err != nil {
		return nil, err
	}

	return target, nil
}

func (s *ReceiptService) Complete(ctx context.Context, id string) (*domain.Receipt, error) {
	receipt, err := s.receiptRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := receipt.Complete(); err != nil {
		return nil, err
	}

	if err := s.receiptRepo.UpdateStatus(ctx, receipt); err != nil {
		return nil, err
	}

	items, err := s.receiptRepo.FindAllItems(ctx, id)
	if err != nil {
		return nil, err
	}

	receipt.Items = items

	return receipt, nil
}

func (s *ReceiptService) Cancel(ctx context.Context, id string) (*domain.Receipt, error) {
	receipt, err := s.receiptRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := receipt.Cancel(); err != nil {
		return nil, err
	}

	if err := s.receiptRepo.UpdateStatus(ctx, receipt); err != nil {
		return nil, err
	}

	po, err := s.poRepo.FindByID(ctx, receipt.PurchaseOrderID)
	if err != nil {
		return nil, err
	}

	if err := po.Cancel(); err != nil {
		return nil, err
	}

	if err := s.poRepo.UpdateStatus(ctx, po); err != nil {
		return nil, err
	}

	return receipt, nil
}
