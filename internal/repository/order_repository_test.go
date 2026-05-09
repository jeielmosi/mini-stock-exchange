package repository

import (
	"context"
	"testing"
	"time"

	"mini-stock-exchange/internal/domain"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderRepository(t *testing.T) {
	ctx := context.Background()
	db, cleanup := SetupTestDB(ctx)
	defer cleanup()

	repo := NewOrderRepository(db)

	t.Run("Create and GetByID", func(t *testing.T) {
		order := &domain.Order{
			ID:                uuid.New(),
			BrokerID:          "broker1",
			OwnerDoc:          "doc1",
			Type:              domain.Bid,
			Symbol:            "AAPL",
			Price:             decimal.NewFromFloat(150.0),
			Quantity:          100,
			RemainingQuantity: 100,
			ValidUntil:        time.Now().Add(time.Hour),
			Status:            domain.Pending,
			CreatedAt:         time.Now().UTC().Truncate(time.Microsecond),
		}

		err := repo.Insert(*order)
		require.NoError(t, err)

		fetched, err := repo.GetByID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, order.ID, fetched.ID)
		assert.Equal(t, order.Symbol, fetched.Symbol)
		assert.True(t, order.Price.Equal(fetched.Price))
	})

	t.Run("Update", func(t *testing.T) {
		order := &domain.Order{
			ID:                uuid.New(),
			BrokerID:          "broker1",
			OwnerDoc:          "doc1",
			Type:              domain.Bid,
			Symbol:            "AAPL",
			Price:             decimal.NewFromFloat(150.0),
			Quantity:          100,
			RemainingQuantity: 100,
			ValidUntil:        time.Now().Add(time.Hour),
			Status:            domain.Pending,
			CreatedAt:         time.Now().UTC().Truncate(time.Microsecond),
		}
		err := repo.Insert(*order)
		require.NoError(t, err)

		order.RemainingQuantity = 50
		order.Status = domain.Partial
		err = repo.Update(*order)
		require.NoError(t, err)

		fetched, err := repo.GetByID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, 50, fetched.RemainingQuantity)
		assert.Equal(t, domain.Partial, fetched.Status)
	})

	t.Run("GetBids and GetAsks", func(t *testing.T) {
		symbol := "SYM_" + uuid.New().String()[:10]
		orders := []domain.Order{
			{
				ID:                uuid.New(),
				BrokerID:          "b1",
				OwnerDoc:          "d1",
				Type:              domain.Bid,
				Symbol:            symbol,
				Price:             decimal.NewFromFloat(150.0),
				Quantity:          10,
				RemainingQuantity: 10,
				ValidUntil:        time.Now().Add(time.Hour),
				Status:            domain.Pending,
				CreatedAt:         time.Now().UTC().Truncate(time.Microsecond),
			},
			{
				ID:                uuid.New(),
				BrokerID:          "b2",
				OwnerDoc:          "d2",
				Type:              domain.Ask,
				Symbol:            symbol,
				Price:             decimal.NewFromFloat(151.0),
				Quantity:          10,
				RemainingQuantity: 10,
				ValidUntil:        time.Now().Add(time.Hour),
				Status:            domain.Pending,
				CreatedAt:         time.Now().UTC().Truncate(time.Microsecond),
			},
		}
		for _, o := range orders {
			err := repo.Insert(o)
			require.NoError(t, err)
		}

		bids, err := repo.GetBids(symbol)
		require.NoError(t, err)
		assert.Equal(t, 1, len(bids))
		assert.Equal(t, domain.Bid, bids[0].Type)

		asks, err := repo.GetAsks(symbol)
		require.NoError(t, err)
		assert.Equal(t, 1, len(asks))
		assert.Equal(t, domain.Ask, asks[0].Type)
	})
}
