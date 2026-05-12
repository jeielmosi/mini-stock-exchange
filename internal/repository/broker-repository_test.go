package repository

import (
	"context"
	"testing"

	"mini-stock-exchange/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBrokerRepository(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := SetupTestDB(ctx)
	assert.NoError(t, err)
	defer cleanup()

	repo, err := NewBrokerRepository(db)
	assert.NoError(t, err)

	t.Run("GetByID", func(t *testing.T) {
		broker := entity.Broker{
			ID:   uuid.New(),
			Name: "Another Broker",
		}

		err := repo.Insert(broker)
		require.NoError(t, err)

		fetched, err := repo.GetByID(broker.ID)
		require.NoError(t, err)
		assert.Equal(t, broker.ID, fetched.ID)
		assert.Equal(t, broker.Name, fetched.Name)
	})
}
