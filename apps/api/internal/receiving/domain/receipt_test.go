package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/domain"
)

func TestNewReceipt(t *testing.T) {
	t.Run("should create a valid receipt", func(t *testing.T) {
		r, err := domain.NewReceipt("po-1", "note")

		require.NoError(t, err)
		assert.NotEmpty(t, r.ID)
		assert.Equal(t, "po-1", r.PurchaseOrderID)
		assert.Equal(t, domain.ReceiptStatusPending, r.Status)
	})

	t.Run("should fail when purchase_order_id is empty", func(t *testing.T) {
		r, err := domain.NewReceipt("", "note")

		assert.Error(t, err)
		assert.Nil(t, r)
	})
}

func TestReceiptStatusFlow(t *testing.T) {
	t.Run("should follow the correct status flow", func(t *testing.T) {
		r, _ := domain.NewReceipt("po-1", "")

		require.NoError(t, r.Start())
		assert.Equal(t, domain.ReceiptStatusInProgress, r.Status)

		require.NoError(t, r.Complete())
		assert.Equal(t, domain.ReceiptStatusCompleted, r.Status)
	})

	t.Run("should not complete when pending", func(t *testing.T) {
		r, _ := domain.NewReceipt("po-1", "")

		err := r.Complete()

		assert.Error(t, err)
	})
}

func TestReceiptCancel(t *testing.T) {
	t.Run("should cancel a pending receipt", func(t *testing.T) {
		r, _ := domain.NewReceipt("po-1", "")

		err := r.Cancel()

		require.NoError(t, err)
		assert.Equal(t, domain.ReceiptStatusCancelled, r.Status)
	})

	t.Run("should not cancel a completed receipt", func(t *testing.T) {
		r, _ := domain.NewReceipt("po-1", "")
		r.Start()
		r.Complete()

		err := r.Cancel()

		assert.Error(t, err)
	})
}

func TestReceiptItemReceive(t *testing.T) {
	t.Run("should set received quantity", func(t *testing.T) {
		item, _ := domain.NewReceiptItem("receipt-1", "product-1", 50)

		err := item.Receive(48)

		require.NoError(t, err)
		require.NotNil(t, item.ReceivedQuantity)
		assert.Equal(t, 48.0, *item.ReceivedQuantity)
	})

	t.Run("should calculate negative difference", func(t *testing.T) {
		item, _ := domain.NewReceiptItem("receipt-1", "product-1", 50)
		item.Receive(48)

		diff := item.Difference()

		require.NotNil(t, diff)
		assert.Equal(t, -2.0, *diff)
	})

	t.Run("should calculate positive difference", func(t *testing.T) {
		item, _ := domain.NewReceiptItem("receipt-1", "product-1", 50)
		item.Receive(55)

		diff := item.Difference()

		require.NotNil(t, diff)
		assert.Equal(t, 5.0, *diff)
	})

	t.Run("should return nil difference when not received", func(t *testing.T) {
		item, _ := domain.NewReceiptItem("receipt-1", "product-1", 50)

		diff := item.Difference()

		assert.Nil(t, diff)
	})

	t.Run("should fail when received quantity is zero", func(t *testing.T) {
		item, _ := domain.NewReceiptItem("receipt-1", "product-1", 50)

		err := item.Receive(0)

		assert.Error(t, err)
	})
}
