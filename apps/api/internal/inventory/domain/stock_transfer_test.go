package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
)

func TestNewStockTransfer(t *testing.T) {
	t.Run("should create a valid transfer", func(t *testing.T) {
		transfer, err := domain.NewStockTransfer("product-1", "location-1", "location-2", 10, "note")

		require.NoError(t, err)
		assert.NotEmpty(t, transfer.ID)
		assert.Equal(t, domain.TransferStatusPending, transfer.Status)
	})

	t.Run("should fail when from and to location are the same", func(t *testing.T) {
		transfer, err := domain.NewStockTransfer("product-1", "location-1", "location-1", 10, "")

		assert.Error(t, err)
		assert.Nil(t, transfer)
	})

	t.Run("should fail when quantity is zero", func(t *testing.T) {
		transfer, err := domain.NewStockTransfer("product-1", "location-1", "location-2", 0, "")

		assert.Error(t, err)
		assert.Nil(t, transfer)
	})
}

func TestStockTransferComplete(t *testing.T) {
	t.Run("should complete a pending transfer", func(t *testing.T) {
		transfer, _ := domain.NewStockTransfer("product-1", "location-1", "location-2", 10, "")

		err := transfer.Complete()

		require.NoError(t, err)
		assert.Equal(t, domain.TransferStatusCompleted, transfer.Status)
	})

	t.Run("should fail when transfer is already completed", func(t *testing.T) {
		transfer, _ := domain.NewStockTransfer("product-1", "location-1", "location-2", 10, "")
		transfer.Complete()

		err := transfer.Complete()

		assert.Error(t, err)
	})

	t.Run("should fail when transfer is cancelled", func(t *testing.T) {
		transfer, _ := domain.NewStockTransfer("product-1", "location-1", "location-2", 10, "")
		transfer.Cancel()

		err := transfer.Complete()

		assert.Error(t, err)
	})
}

func TestStockTransferCancel(t *testing.T) {
	t.Run("should cancel a pending transfer", func(t *testing.T) {
		transfer, _ := domain.NewStockTransfer("product-1", "location-1", "location-2", 10, "")

		err := transfer.Cancel()

		require.NoError(t, err)
		assert.Equal(t, domain.TransferStatusCancelled, transfer.Status)
	})

	t.Run("should fail when transfer is already completed", func(t *testing.T) {
		transfer, _ := domain.NewStockTransfer("product-1", "location-1", "location-2", 10, "")
		transfer.Complete()

		err := transfer.Cancel()

		assert.Error(t, err)
	})
}
