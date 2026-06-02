package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/domain"
)

func TestNewZone(t *testing.T) {
	t.Run("should create a valid zone", func(t *testing.T) {
		z, err := domain.NewZone("warehouse-1", "A", "Zone A")

		require.NoError(t, err)
		assert.NotEmpty(t, z.ID)
		assert.Equal(t, "warehouse-1", z.WarehouseID)
		assert.Equal(t, "A", z.Code)
		assert.True(t, z.Active)
	})

	t.Run("should fail when warehouse_id is empty", func(t *testing.T) {
		z, err := domain.NewZone("", "A", "Zone A")

		assert.Error(t, err)
		assert.Nil(t, z)
	})

	t.Run("should fail when code is empty", func(t *testing.T) {
		z, err := domain.NewZone("warehouse-1", "", "Zone A")

		assert.Error(t, err)
		assert.Nil(t, z)
	})

	t.Run("should fail when name is empty", func(t *testing.T) {
		z, err := domain.NewZone("warehouse-1", "A", "")

		assert.Error(t, err)
		assert.Nil(t, z)
	})
}

func TestZoneUpdate(t *testing.T) {
	t.Run("should update a zone", func(t *testing.T) {
		z, _ := domain.NewZone("warehouse-1", "A", "Zone A")

		err := z.Update("B", "Zone B")

		require.NoError(t, err)
		assert.Equal(t, "B", z.Code)
		assert.Equal(t, "Zone B", z.Name)
	})

	t.Run("should fail when code is empty", func(t *testing.T) {
		z, _ := domain.NewZone("warehouse-1", "A", "Zone A")

		err := z.Update("", "Zone B")

		assert.Error(t, err)
		assert.Equal(t, "A", z.Code)
	})
}

func TestZoneDeactivate(t *testing.T) {
	t.Run("should deactivate a zone", func(t *testing.T) {
		z, _ := domain.NewZone("warehouse-1", "A", "Zone A")

		z.Deactivate()

		assert.False(t, z.Active)
	})
}
