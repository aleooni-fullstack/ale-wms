package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/domain"
)

func TestNewLocation(t *testing.T) {
	t.Run("should create a valid location", func(t *testing.T) {
		l, err := domain.NewLocation("zone-1", "A-01", "Position A-01")

		require.NoError(t, err)
		assert.NotEmpty(t, l.ID)
		assert.Equal(t, "zone-1", l.ZoneID)
		assert.Equal(t, "A-01", l.Code)
		assert.True(t, l.Active)
	})

	t.Run("should fail when zone_id is empty", func(t *testing.T) {
		l, err := domain.NewLocation("", "A-01", "Position A-01")

		assert.Error(t, err)
		assert.Nil(t, l)
	})

	t.Run("should fail when code is empty", func(t *testing.T) {
		l, err := domain.NewLocation("zone-1", "", "Position A-01")

		assert.Error(t, err)
		assert.Nil(t, l)
	})

	t.Run("should fail when name is empty", func(t *testing.T) {
		l, err := domain.NewLocation("zone-1", "A-01", "")

		assert.Error(t, err)
		assert.Nil(t, l)
	})
}

func TestLocationUpdate(t *testing.T) {
	t.Run("should update a location", func(t *testing.T) {
		l, _ := domain.NewLocation("zone-1", "A-01", "Position A-01")

		err := l.Update("A-02", "Position A-02")

		require.NoError(t, err)
		assert.Equal(t, "A-02", l.Code)
		assert.Equal(t, "Position A-02", l.Name)
	})

	t.Run("should fail when code is empty", func(t *testing.T) {
		l, _ := domain.NewLocation("zone-1", "A-01", "Position A-01")

		err := l.Update("", "Position A-02")

		assert.Error(t, err)
		assert.Equal(t, "A-01", l.Code)
	})
}

func TestLocationDeactivate(t *testing.T) {
	t.Run("should deactivate a location", func(t *testing.T) {
		l, _ := domain.NewLocation("zone-1", "A-01", "Position A-01")

		l.Deactivate()

		assert.False(t, l.Active)
	})
}
