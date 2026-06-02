package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/domain"
)

func TestNewOrder(t *testing.T) {
	t.Run("should create a valid order", func(t *testing.T) {
		order, err := domain.NewOrder("PED-001", "note")

		require.NoError(t, err)
		assert.NotEmpty(t, order.ID)
		assert.Equal(t, "PED-001", order.Reference)
		assert.Equal(t, domain.OrderStatusDraft, order.Status)
		assert.Empty(t, order.Items)
	})

	t.Run("should fail when reference is empty", func(t *testing.T) {
		order, err := domain.NewOrder("", "note")

		assert.Error(t, err)
		assert.Nil(t, order)
	})
}

func TestOrderConfirm(t *testing.T) {
	t.Run("should confirm an order with items", func(t *testing.T) {
		order, _ := domain.NewOrder("PED-001", "")
		item, _ := domain.NewOrderItem(order.ID, "product-1", "location-1", 10)
		order.Items = append(order.Items, item)

		err := order.Confirm()

		require.NoError(t, err)
		assert.Equal(t, domain.OrderStatusConfirmed, order.Status)
	})

	t.Run("should fail when order has no items", func(t *testing.T) {
		order, _ := domain.NewOrder("PED-001", "")

		err := order.Confirm()

		assert.Error(t, err)
		assert.Equal(t, domain.OrderStatusDraft, order.Status)
	})

	t.Run("should fail when order is not draft", func(t *testing.T) {
		order, _ := domain.NewOrder("PED-001", "")
		item, _ := domain.NewOrderItem(order.ID, "product-1", "location-1", 10)
		order.Items = append(order.Items, item)
		order.Confirm()

		err := order.Confirm()

		assert.Error(t, err)
	})
}

func TestOrderStatusFlow(t *testing.T) {
	t.Run("should follow the correct status flow", func(t *testing.T) {
		order, _ := domain.NewOrder("PED-001", "")
		item, _ := domain.NewOrderItem(order.ID, "product-1", "location-1", 10)
		order.Items = append(order.Items, item)

		require.NoError(t, order.Confirm())
		assert.Equal(t, domain.OrderStatusConfirmed, order.Status)

		require.NoError(t, order.StartPicking())
		assert.Equal(t, domain.OrderStatusPicking, order.Status)

		require.NoError(t, order.StartPacking())
		assert.Equal(t, domain.OrderStatusPacking, order.Status)

		require.NoError(t, order.Ship())
		assert.Equal(t, domain.OrderStatusShipped, order.Status)
	})

	t.Run("should not skip steps", func(t *testing.T) {
		order, _ := domain.NewOrder("PED-001", "")
		item, _ := domain.NewOrderItem(order.ID, "product-1", "location-1", 10)
		order.Items = append(order.Items, item)
		order.Confirm()

		err := order.StartPacking()

		assert.Error(t, err)
	})
}

func TestOrderCancel(t *testing.T) {
	t.Run("should cancel a draft order", func(t *testing.T) {
		order, _ := domain.NewOrder("PED-001", "")

		err := order.Cancel()

		require.NoError(t, err)
		assert.Equal(t, domain.OrderStatusCancelled, order.Status)
	})

	t.Run("should not cancel a shipped order", func(t *testing.T) {
		order, _ := domain.NewOrder("PED-001", "")
		item, _ := domain.NewOrderItem(order.ID, "product-1", "location-1", 10)
		order.Items = append(order.Items, item)
		order.Confirm()
		order.StartPicking()
		order.StartPacking()
		order.Ship()

		err := order.Cancel()

		assert.Error(t, err)
		assert.Equal(t, domain.OrderStatusShipped, order.Status)
	})
}
