package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
)

func TestNewStockMovement(t *testing.T) {
	t.Run("should create a valid IN movement", func(t *testing.T) {
		m, err := domain.NewStockMovement("product-1", "location-1", domain.MovementTypeIn, 10, "note")

		require.NoError(t, err)
		assert.NotEmpty(t, m.ID)
		assert.Equal(t, domain.MovementTypeIn, m.Type)
		assert.Equal(t, 10.0, m.Quantity)
	})

	t.Run("should fail when product_id is empty", func(t *testing.T) {
		m, err := domain.NewStockMovement("", "location-1", domain.MovementTypeIn, 10, "")

		assert.Error(t, err)
		assert.Nil(t, m)
	})

	t.Run("should fail when location_id is empty", func(t *testing.T) {
		m, err := domain.NewStockMovement("product-1", "", domain.MovementTypeIn, 10, "")

		assert.Error(t, err)
		assert.Nil(t, m)
	})

	t.Run("should fail when quantity is zero", func(t *testing.T) {
		m, err := domain.NewStockMovement("product-1", "location-1", domain.MovementTypeIn, 0, "")

		assert.Error(t, err)
		assert.Nil(t, m)
	})

	t.Run("should fail when quantity is negative", func(t *testing.T) {
		m, err := domain.NewStockMovement("product-1", "location-1", domain.MovementTypeIn, -5, "")

		assert.Error(t, err)
		assert.Nil(t, m)
	})

	t.Run("should fail when movement type is invalid", func(t *testing.T) {
		m, err := domain.NewStockMovement("product-1", "location-1", "INVALID", 10, "")

		assert.Error(t, err)
		assert.Nil(t, m)
	})
}
