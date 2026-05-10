package repository

import (
	"context"
	"testing"
	"time"

	"mini-stock-exchange/internal/entity"

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
		order := &entity.Order{
			ID:                uuid.New(),
			BrokerID:          "broker1",
			OwnerDoc:          "doc1",
			Type:              entity.Bid,
			Symbol:            "AAPL",
			Price:             decimal.NewFromFloat(150.0),
			Quantity:          100,
			RemainingQuantity: 100,
			ValidUntil:        time.Now().Add(time.Hour),
			Status:            entity.Pending,
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

	t.Run("GetBids and GetAsks", func(t *testing.T) {
		symbol := "SYM_" + uuid.New().String()[:10]
		orders := []entity.Order{
			{
				ID:                uuid.New(),
				BrokerID:          "b1",
				OwnerDoc:          "d1",
				Type:              entity.Bid,
				Symbol:            symbol,
				Price:             decimal.NewFromFloat(150.0),
				Quantity:          10,
				RemainingQuantity: 10,
				ValidUntil:        time.Now().Add(time.Hour),
				Status:            entity.Pending,
				CreatedAt:         time.Now().UTC().Truncate(time.Microsecond),
			},
			{
				ID:                uuid.New(),
				BrokerID:          "b2",
				OwnerDoc:          "d2",
				Type:              entity.Ask,
				Symbol:            symbol,
				Price:             decimal.NewFromFloat(151.0),
				Quantity:          10,
				RemainingQuantity: 10,
				ValidUntil:        time.Now().Add(time.Hour),
				Status:            entity.Pending,
				CreatedAt:         time.Now().UTC().Truncate(time.Microsecond),
			},
		}
		for _, o := range orders {
			err := repo.Insert(o)
			require.NoError(t, err)
		}

		bids, err := repo.GetBids(symbol, 100)
		require.NoError(t, err)
		assert.Equal(t, 1, len(bids))
		assert.Equal(t, entity.Bid, bids[0].Type)

		asks, err := repo.GetAsks(symbol, 100)
		require.NoError(t, err)
		assert.Equal(t, 1, len(asks))
		assert.Equal(t, entity.Ask, asks[0].Type)
	})

	t.Run("Update", func(t *testing.T) {
		askID := uuid.New()
		bidID := uuid.New()
		tradeID := uuid.New()
		symbol := "AAPL"
		now := time.Now().UTC().Truncate(time.Microsecond)

		ask := entity.Order{
			ID:                askID,
			BrokerID:          "broker1",
			OwnerDoc:          "doc1",
			Type:              entity.Ask,
			Symbol:            symbol,
			Price:             decimal.NewFromFloat(150.0),
			Quantity:          100,
			RemainingQuantity: 100,
			ValidUntil:        time.Now().Add(time.Hour),
			Status:            entity.Pending,
			CreatedAt:         now,
		}

		bid := entity.Order{
			ID:                bidID,
			BrokerID:          "broker2",
			OwnerDoc:          "doc2",
			Type:              entity.Bid,
			Symbol:            symbol,
			Price:             decimal.NewFromFloat(150.0),
			Quantity:          100,
			RemainingQuantity: 100,
			ValidUntil:        time.Now().Add(time.Hour),
			Status:            entity.Pending,
			CreatedAt:         now,
		}

		err := repo.Insert(ask)
		require.NoError(t, err)
		err = repo.Insert(bid)
		require.NoError(t, err)

		ask.RemainingQuantity = 0
		ask.Status = entity.Filled

		bid.RemainingQuantity = 0
		bid.Status = entity.Filled

		trade := entity.Trade{
			ID:          tradeID,
			Symbol:      symbol,
			Price:       decimal.NewFromFloat(150.0),
			Quantity:    100,
			ExecutedAt:  now,
			BuyOrderID:  bidID,
			SellOrderID: askID,
		}

		match := MatchDTO{
			Ask:   ask,
			Bid:   bid,
			Trade: trade,
		}

		ctx := context.Background()
		err = repo.Match(ctx, match)
		require.NoError(t, err)

		updatedAsk, err := repo.GetByID(askID)
		require.NoError(t, err)
		assert.Equal(t, 0, updatedAsk.RemainingQuantity)
		assert.Equal(t, entity.Filled, updatedAsk.Status)

		updatedBid, err := repo.GetByID(bidID)
		require.NoError(t, err)
		assert.Equal(t, 0, updatedBid.RemainingQuantity)
		assert.Equal(t, entity.Filled, updatedBid.Status)
	})
}
