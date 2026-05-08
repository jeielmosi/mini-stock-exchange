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

		err := repo.Create(order)
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
		err := repo.Create(order)
		require.NoError(t, err)

		order.RemainingQuantity = 50
		order.Status = domain.Partial
		err = repo.Update(order)
		require.NoError(t, err)

		fetched, err := repo.GetByID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, 50, fetched.RemainingQuantity)
		assert.Equal(t, domain.Partial, fetched.Status)
	})

	t.Run("FindMatches", func(t *testing.T) {
		// Clear or use specific symbols for this test
		symbol := "MSFT"

		// Case 1: BID order looking for ASKs (price <= bid_price)
		ask1 := &domain.Order{
			ID:                uuid.New(),
			BrokerID:          "b1",
			OwnerDoc:          "d1",
			Type:              domain.Ask,
			Symbol:            symbol,
			Price:             decimal.NewFromFloat(100.0),
			Quantity:          10,
			RemainingQuantity: 10,
			ValidUntil:        time.Now().Add(time.Hour),
			Status:            domain.Pending,
			CreatedAt:         time.Now().UTC().Truncate(time.Microsecond),
		}
		ask2 := &domain.Order{
			ID:                uuid.New(),
			BrokerID:          "b2",
			OwnerDoc:          "d2",
			Type:              domain.Ask,
			Symbol:            symbol,
			Price:             decimal.NewFromFloat(110.0),
			Quantity:          10,
			RemainingQuantity: 10,
			ValidUntil:        time.Now().Add(time.Hour),
			Status:            domain.Pending,
			CreatedAt:         time.Now().UTC().Truncate(time.Microsecond),
		}
		askExpired := &domain.Order{
			ID:                uuid.New(),
			BrokerID:          "b3",
			OwnerDoc:          "d3",
			Type:              domain.Ask,
			Symbol:            symbol,
			Price:             decimal.NewFromFloat(90.0),
			Quantity:          10,
			RemainingQuantity: 10,
			ValidUntil:        time.Now().Add(-time.Hour),
			Status:            domain.Pending,
			CreatedAt:         time.Now().UTC().Truncate(time.Microsecond),
		}

		require.NoError(t, repo.Create(ask1))
		require.NoError(t, repo.Create(ask2))
		require.NoError(t, repo.Create(askExpired))

		matches, err := repo.FindMatches(domain.Order{
			Symbol:   symbol,
			Type:     domain.Bid,
			Price:    decimal.NewFromFloat(105.0),
			Quantity: 10,
		})
		require.NoError(t, err)
		assert.Len(t, matches, 1)
		assert.Equal(t, ask1.ID, matches[0].ID)
		t.Log(matches)

		matches, err = repo.FindMatches(domain.Order{
			Symbol:   symbol,
			Type:     domain.Bid,
			Price:    decimal.NewFromFloat(115.0),
			Quantity: 11,
		})
		require.NoError(t, err)
		assert.Len(t, matches, 2)
		t.Log(matches)

		// Case 2: ASK order looking for BIDs (price >= ask_price)
		bid1 := &domain.Order{
			ID:                uuid.New(),
			BrokerID:          "b4",
			OwnerDoc:          "d4",
			Type:              domain.Bid,
			Symbol:            symbol,
			Price:             decimal.NewFromFloat(100.0),
			Quantity:          10,
			RemainingQuantity: 10,
			ValidUntil:        time.Now().Add(time.Hour),
			Status:            domain.Pending,
			CreatedAt:         time.Now().Add(-time.Minute).UTC().Truncate(time.Microsecond),
		}
		bid2 := &domain.Order{
			ID:                uuid.New(),
			BrokerID:          "b5",
			OwnerDoc:          "d5",
			Type:              domain.Bid,
			Symbol:            symbol,
			Price:             decimal.NewFromFloat(110.0),
			Quantity:          10,
			RemainingQuantity: 10,
			ValidUntil:        time.Now().Add(time.Hour),
			Status:            domain.Pending,
			CreatedAt:         time.Now().UTC().Truncate(time.Microsecond),
		}

		require.NoError(t, repo.Create(bid1))
		require.NoError(t, repo.Create(bid2))

		matches, err = repo.FindMatches(domain.Order{
			Symbol:   symbol,
			Type:     domain.Ask,
			Price:    decimal.NewFromFloat(105.0),
			Quantity: 10,
		})
		require.NoError(t, err)
		assert.Len(t, matches, 1)
		assert.Equal(t, bid2.ID, matches[0].ID)

		// Case 3: Price and Time Priority for ASK order
		// Two Bids with same price, different times
		bid3 := &domain.Order{
			ID:                uuid.New(),
			BrokerID:          "b6",
			OwnerDoc:          "d6",
			Type:              domain.Bid,
			Symbol:            symbol,
			Price:             decimal.NewFromFloat(120.0),
			Quantity:          10,
			RemainingQuantity: 10,
			ValidUntil:        time.Now().Add(time.Hour),
			Status:            domain.Pending,
			CreatedAt:         time.Now().Add(-2 * time.Minute).UTC().Truncate(time.Microsecond),
		}
		bid4 := &domain.Order{
			ID:                uuid.New(),
			BrokerID:          "b7",
			OwnerDoc:          "d7",
			Type:              domain.Bid,
			Symbol:            symbol,
			Price:             decimal.NewFromFloat(120.0),
			Quantity:          10,
			RemainingQuantity: 10,
			ValidUntil:        time.Now().Add(time.Hour),
			Status:            domain.Pending,
			CreatedAt:         time.Now().Add(-1 * time.Minute).UTC().Truncate(time.Microsecond),
		}
		require.NoError(t, repo.Create(bid3))
		require.NoError(t, repo.Create(bid4))

		matches, err = repo.FindMatches(domain.Order{
			Symbol:   symbol,
			Type:     domain.Ask,
			Price:    decimal.NewFromFloat(120.0),
			Quantity: 20,
		})
		require.NoError(t, err)
		assert.Len(t, matches, 2)
		assert.Equal(t, bid3.ID, matches[0].ID, "Oldest bid should be matched first")
		assert.Equal(t, bid4.ID, matches[1].ID)
	})

	t.Run("UpdateRemainingQuantity", func(t *testing.T) {
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
		err := repo.Create(order)
		require.NoError(t, err)

		err = repo.UpdateRemainingQuantity(order.ID, 20, domain.Partial)
		require.NoError(t, err)

		fetched, err := repo.GetByID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, 20, fetched.RemainingQuantity)
		assert.Equal(t, domain.Partial, fetched.Status)
	})
}
