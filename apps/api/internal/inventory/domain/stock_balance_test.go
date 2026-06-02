package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
)

func TestStockBalanceAvailableQuantity(t *testing.T) {
	t.Run("should return quantity minus reserved", func(t *testing.T) {
		b := domain.NewStockBalance("product-1", "location-1", 100)
		b.ReservedQuantity = 30

		assert.Equal(t, 70.0, b.AvailableQuantity())
	})

	t.Run("should return zero when all is reserved", func(t *testing.T) {
		b := domain.NewStockBalance("product-1", "location-1", 100)
		b.ReservedQuantity = 100

		assert.Equal(t, 0.0, b.AvailableQuantity())
	})
}

func TestStockBalanceApply(t *testing.T) {
	t.Run("should increase quantity on IN movement", func(t *testing.T) {
		b := domain.NewStockBalance("product-1", "location-1", 50)
		m, _ := domain.NewStockMovement("product-1", "location-1", domain.MovementTypeIn, 30, "")

		b.Apply(m)

		assert.Equal(t, 80.0, b.Quantity)
	})

	t.Run("should decrease quantity on OUT movement", func(t *testing.T) {
		b := domain.NewStockBalance("product-1", "location-1", 50)
		m, _ := domain.NewStockMovement("product-1", "location-1", domain.MovementTypeOut, 20, "")

		b.Apply(m)

		assert.Equal(t, 30.0, b.Quantity)
	})

	t.Run("should set quantity on ADJUSTMENT movement", func(t *testing.T) {
		b := domain.NewStockBalance("product-1", "location-1", 50)
		m, _ := domain.NewStockMovement("product-1", "location-1", domain.MovementTypeAdjustment, 25, "")

		b.Apply(m)

		assert.Equal(t, 25.0, b.Quantity)
	})
}

func TestStockBalanceReserve(t *testing.T) {
	t.Run("should increase reserved quantity", func(t *testing.T) {
		b := domain.NewStockBalance("product-1", "location-1", 100)

		b.Reserve(30)

		assert.Equal(t, 30.0, b.ReservedQuantity)
		assert.Equal(t, 70.0, b.AvailableQuantity())
	})
}

func TestStockBalanceRelease(t *testing.T) {
	t.Run("should decrease reserved quantity", func(t *testing.T) {
		b := domain.NewStockBalance("product-1", "location-1", 100)
		b.Reserve(30)

		b.Release(10)

		assert.Equal(t, 20.0, b.ReservedQuantity)
	})

	t.Run("should not go below zero", func(t *testing.T) {
		b := domain.NewStockBalance("product-1", "location-1", 100)
		b.Reserve(10)

		b.Release(50)

		assert.Equal(t, 0.0, b.ReservedQuantity)
	})
}
