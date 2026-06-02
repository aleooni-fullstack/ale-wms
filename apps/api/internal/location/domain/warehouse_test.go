package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/domain"
)

func TestNewWarehouse(t *testing.T) {
	t.Run("should create a valid warehouse", func(t *testing.T) {
		w, err := domain.NewWarehouse("WH-01", "Main Warehouse", "Address")

		require.NoError(t, err)
		assert.NotEmpty(t, w.ID)
		assert.Equal(t, "WH-01", w.Code)
		assert.Equal(t, "Main Warehouse", w.Name)
		assert.True(t, w.Active)
	})

	t.Run("should fail when code is empty", func(t *testing.T) {
		w, err := domain.NewWarehouse("", "Main Warehouse", "Address")

		assert.Error(t, err)
		assert.Nil(t, w)
	})

	t.Run("should fail when name is empty", func(t *testing.T) {
		w, err := domain.NewWarehouse("WH-01", "", "Address")

		assert.Error(t, err)
		assert.Nil(t, w)
	})
}

func TestWarehouseUpdate(t *testing.T) {
	t.Run("should update a warehouse", func(t *testing.T) {
		w, _ := domain.NewWarehouse("WH-01", "Main Warehouse", "Address")

		err := w.Update("WH-02", "Secondary Warehouse", "New Address")

		require.NoError(t, err)
		assert.Equal(t, "WH-02", w.Code)
		assert.Equal(t, "Secondary Warehouse", w.Name)
		assert.Equal(t, "New Address", w.Address)
	})

	t.Run("should fail when code is empty", func(t *testing.T) {
		w, _ := domain.NewWarehouse("WH-01", "Main Warehouse", "Address")

		err := w.Update("", "Secondary Warehouse", "New Address")

		assert.Error(t, err)
		assert.Equal(t, "WH-01", w.Code)
	})
}

func TestWarehouseDeactivate(t *testing.T) {
	t.Run("should deactivate a warehouse", func(t *testing.T) {
		w, _ := domain.NewWarehouse("WH-01", "Main Warehouse", "Address")

		w.Deactivate()

		assert.False(t, w.Active)
	})
}
