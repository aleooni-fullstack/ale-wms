package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/catalog/domain"
)

func TestNewProduct(t *testing.T) {
	t.Run("should create a valid product", func(t *testing.T) {
		p, err := domain.NewProduct("SKU-001", "Product A", "Description", "UN")

		require.NoError(t, err)
		assert.NotEmpty(t, p.ID)
		assert.Equal(t, "SKU-001", p.SKU)
		assert.Equal(t, "Product A", p.Name)
		assert.Equal(t, "Description", p.Description)
		assert.Equal(t, "UN", p.Unit)
		assert.True(t, p.Active)
	})

	t.Run("should fail when sku is empty", func(t *testing.T) {
		p, err := domain.NewProduct("", "Product A", "Description", "UN")

		assert.Error(t, err)
		assert.Nil(t, p)
	})

	t.Run("should fail when name is empty", func(t *testing.T) {
		p, err := domain.NewProduct("SKU-001", "", "Description", "UN")

		assert.Error(t, err)
		assert.Nil(t, p)
	})

	t.Run("should fail when unit is empty", func(t *testing.T) {
		p, err := domain.NewProduct("SKU-001", "Product A", "Description", "")

		assert.Error(t, err)
		assert.Nil(t, p)
	})
}

func TestProductUpdate(t *testing.T) {
	t.Run("should update a product", func(t *testing.T) {
		p, _ := domain.NewProduct("SKU-001", "Product A", "Description", "UN")

		err := p.Update("SKU-002", "Product B", "New Description", "KG")

		require.NoError(t, err)
		assert.Equal(t, "SKU-002", p.SKU)
		assert.Equal(t, "Product B", p.Name)
		assert.Equal(t, "New Description", p.Description)
		assert.Equal(t, "KG", p.Unit)
	})

	t.Run("should fail when sku is empty", func(t *testing.T) {
		p, _ := domain.NewProduct("SKU-001", "Product A", "Description", "UN")

		err := p.Update("", "Product B", "New Description", "KG")

		assert.Error(t, err)
		assert.Equal(t, "SKU-001", p.SKU)
	})

	t.Run("should fail when name is empty", func(t *testing.T) {
		p, _ := domain.NewProduct("SKU-001", "Product A", "Description", "UN")

		err := p.Update("SKU-002", "", "New Description", "KG")

		assert.Error(t, err)
	})
}

func TestProductDeactivate(t *testing.T) {
	t.Run("should deactivate a product", func(t *testing.T) {
		p, _ := domain.NewProduct("SKU-001", "Product A", "Description", "UN")

		p.Deactivate()

		assert.False(t, p.Active)
	})
}
