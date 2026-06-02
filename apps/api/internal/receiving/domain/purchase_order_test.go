package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/domain"
)

func TestNewPurchaseOrder(t *testing.T) {
	t.Run("should create a valid purchase order", func(t *testing.T) {
		po, err := domain.NewPurchaseOrder("PO-001", "Supplier A", "note")

		require.NoError(t, err)
		assert.NotEmpty(t, po.ID)
		assert.Equal(t, "PO-001", po.Reference)
		assert.Equal(t, "Supplier A", po.Supplier)
		assert.Equal(t, domain.PurchaseOrderStatusDraft, po.Status)
		assert.Empty(t, po.Items)
	})

	t.Run("should fail when reference is empty", func(t *testing.T) {
		po, err := domain.NewPurchaseOrder("", "Supplier A", "note")

		assert.Error(t, err)
		assert.Nil(t, po)
	})
}

func TestPurchaseOrderConfirm(t *testing.T) {
	t.Run("should confirm with items", func(t *testing.T) {
		po, _ := domain.NewPurchaseOrder("PO-001", "Supplier A", "")
		item, _ := domain.NewPurchaseOrderItem(po.ID, "product-1", 10)
		po.Items = append(po.Items, item)

		err := po.Confirm()

		require.NoError(t, err)
		assert.Equal(t, domain.PurchaseOrderStatusConfirmed, po.Status)
	})

	t.Run("should fail when no items", func(t *testing.T) {
		po, _ := domain.NewPurchaseOrder("PO-001", "Supplier A", "")

		err := po.Confirm()

		assert.Error(t, err)
		assert.Equal(t, domain.PurchaseOrderStatusDraft, po.Status)
	})

	t.Run("should fail when not draft", func(t *testing.T) {
		po, _ := domain.NewPurchaseOrder("PO-001", "Supplier A", "")
		item, _ := domain.NewPurchaseOrderItem(po.ID, "product-1", 10)
		po.Items = append(po.Items, item)
		po.Confirm()

		err := po.Confirm()

		assert.Error(t, err)
	})
}

func TestPurchaseOrderStatusFlow(t *testing.T) {
	t.Run("should follow the correct status flow", func(t *testing.T) {
		po, _ := domain.NewPurchaseOrder("PO-001", "Supplier A", "")
		item, _ := domain.NewPurchaseOrderItem(po.ID, "product-1", 10)
		po.Items = append(po.Items, item)

		require.NoError(t, po.Confirm())
		assert.Equal(t, domain.PurchaseOrderStatusConfirmed, po.Status)

		require.NoError(t, po.StartReceiving())
		assert.Equal(t, domain.PurchaseOrderStatusReceiving, po.Status)

		require.NoError(t, po.Complete())
		assert.Equal(t, domain.PurchaseOrderStatusCompleted, po.Status)
	})

	t.Run("should not skip steps", func(t *testing.T) {
		po, _ := domain.NewPurchaseOrder("PO-001", "Supplier A", "")
		item, _ := domain.NewPurchaseOrderItem(po.ID, "product-1", 10)
		po.Items = append(po.Items, item)

		err := po.StartReceiving()

		assert.Error(t, err)
	})
}

func TestPurchaseOrderCancel(t *testing.T) {
	t.Run("should cancel a draft purchase order", func(t *testing.T) {
		po, _ := domain.NewPurchaseOrder("PO-001", "Supplier A", "")

		err := po.Cancel()

		require.NoError(t, err)
		assert.Equal(t, domain.PurchaseOrderStatusCancelled, po.Status)
	})

	t.Run("should not cancel a completed purchase order", func(t *testing.T) {
		po, _ := domain.NewPurchaseOrder("PO-001", "Supplier A", "")
		item, _ := domain.NewPurchaseOrderItem(po.ID, "product-1", 10)
		po.Items = append(po.Items, item)
		po.Confirm()
		po.StartReceiving()
		po.Complete()

		err := po.Cancel()

		assert.Error(t, err)
		assert.Equal(t, domain.PurchaseOrderStatusCompleted, po.Status)
	})
}
