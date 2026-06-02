package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/domain"
)

func TestNewPutAway(t *testing.T) {
	t.Run("should create a valid put away", func(t *testing.T) {
		p, err := domain.NewPutAway("receipt-1", "note")

		require.NoError(t, err)
		assert.NotEmpty(t, p.ID)
		assert.Equal(t, "receipt-1", p.ReceiptID)
		assert.Equal(t, domain.PutAwayStatusPending, p.Status)
	})

	t.Run("should fail when receipt_id is empty", func(t *testing.T) {
		p, err := domain.NewPutAway("", "note")

		assert.Error(t, err)
		assert.Nil(t, p)
	})
}

func TestPutAwayComplete(t *testing.T) {
	t.Run("should complete when all items are stored", func(t *testing.T) {
		p, _ := domain.NewPutAway("receipt-1", "")
		item, _ := domain.NewPutAwayItem(p.ID, "product-1", "location-1", 10)
		item.Store()
		p.Items = append(p.Items, item)
		p.Start()

		err := p.Complete()

		require.NoError(t, err)
		assert.Equal(t, domain.PutAwayStatusCompleted, p.Status)
	})

	t.Run("should fail when not all items are stored", func(t *testing.T) {
		p, _ := domain.NewPutAway("receipt-1", "")
		item, _ := domain.NewPutAwayItem(p.ID, "product-1", "location-1", 10)
		p.Items = append(p.Items, item)
		p.Start()

		err := p.Complete()

		assert.Error(t, err)
		assert.Equal(t, domain.PutAwayStatusInProgress, p.Status)
	})

	t.Run("should fail when status is not in_progress", func(t *testing.T) {
		p, _ := domain.NewPutAway("receipt-1", "")

		err := p.Complete()

		assert.Error(t, err)
	})
}

func TestPutAwayCancel(t *testing.T) {
	t.Run("should cancel a pending put away", func(t *testing.T) {
		p, _ := domain.NewPutAway("receipt-1", "")

		err := p.Cancel()

		require.NoError(t, err)
		assert.Equal(t, domain.PutAwayStatusCancelled, p.Status)
	})

	t.Run("should not cancel a completed put away", func(t *testing.T) {
		p, _ := domain.NewPutAway("receipt-1", "")
		item, _ := domain.NewPutAwayItem(p.ID, "product-1", "location-1", 10)
		item.Store()
		p.Items = append(p.Items, item)
		p.Start()
		p.Complete()

		err := p.Cancel()

		assert.Error(t, err)
	})
}

func TestPutAwayItemStore(t *testing.T) {
	t.Run("should mark item as stored", func(t *testing.T) {
		item, _ := domain.NewPutAwayItem("put-away-1", "product-1", "location-1", 10)

		item.Store()

		assert.True(t, item.PutAway)
	})
}
