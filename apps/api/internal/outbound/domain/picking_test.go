package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/domain"
)

func TestPickingComplete(t *testing.T) {
	t.Run("should complete when all items are picked", func(t *testing.T) {
		picking, _ := domain.NewPicking("order-1", "")
		item, _ := domain.NewPickingItem(picking.ID, "product-1", "location-1", 10)
		item.Pick()
		picking.Items = append(picking.Items, item)
		picking.Start()

		err := picking.Complete()

		require.NoError(t, err)
		assert.Equal(t, domain.PickingStatusCompleted, picking.Status)
	})

	t.Run("should fail when not all items are picked", func(t *testing.T) {
		picking, _ := domain.NewPicking("order-1", "")
		item, _ := domain.NewPickingItem(picking.ID, "product-1", "location-1", 10)
		picking.Items = append(picking.Items, item)
		picking.Start()

		err := picking.Complete()

		assert.Error(t, err)
		assert.Equal(t, domain.PickingStatusInProgress, picking.Status)
	})

	t.Run("should fail when status is not in_progress", func(t *testing.T) {
		picking, _ := domain.NewPicking("order-1", "")

		err := picking.Complete()

		assert.Error(t, err)
	})
}

func TestPickingCancel(t *testing.T) {
	t.Run("should cancel a pending picking", func(t *testing.T) {
		picking, _ := domain.NewPicking("order-1", "")

		err := picking.Cancel()

		require.NoError(t, err)
		assert.Equal(t, domain.PickingStatusCancelled, picking.Status)
	})

	t.Run("should not cancel a completed picking", func(t *testing.T) {
		picking, _ := domain.NewPicking("order-1", "")
		item, _ := domain.NewPickingItem(picking.ID, "product-1", "location-1", 10)
		item.Pick()
		picking.Items = append(picking.Items, item)
		picking.Start()
		picking.Complete()

		err := picking.Cancel()

		assert.Error(t, err)
	})
}

func TestPickingItemPick(t *testing.T) {
	t.Run("should mark item as picked", func(t *testing.T) {
		item, _ := domain.NewPickingItem("picking-1", "product-1", "location-1", 10)

		item.Pick()

		assert.True(t, item.Picked)
	})
}
