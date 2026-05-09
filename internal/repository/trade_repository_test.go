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

func TestTradeRepository(t *testing.T) {
	ctx := context.Background()
	db, cleanup := SetupTestDB(ctx)
	defer cleanup()

	orderRepo := NewOrderRepository(db)
	tradeRepo := NewTradeRepository(db)

	t.Run("Create and GetByOrderID", func(t *testing.T) {
		buyOrder := domain.Order{
			ID:                uuid.New(),
			BrokerID:          "b1",
			OwnerDoc:          "d1",
			Type:              domain.Bid,
			Symbol:            "AAPL",
			Price:             decimal.NewFromFloat(150.0),
			Quantity:          100,
			RemainingQuantity: 100,
			ValidUntil:        time.Now().Add(time.Hour),
			Status:            domain.Pending,
			CreatedAt:         time.Now().UTC().Truncate(time.Microsecond),
		}
		sellOrder := domain.Order{
			ID:                uuid.New(),
			BrokerID:          "b2",
			OwnerDoc:          "d2",
			Type:              domain.Ask,
			Symbol:            "AAPL",
			Price:             decimal.NewFromFloat(150.0),
			Quantity:          100,
			RemainingQuantity: 100,
			ValidUntil:        time.Now().Add(time.Hour),
			Status:            domain.Pending,
			CreatedAt:         time.Now().UTC().Truncate(time.Microsecond),
		}

		require.NoError(t, orderRepo.Insert(buyOrder))
		require.NoError(t, orderRepo.Insert(sellOrder))

		trade := domain.Trade{
			ID:          uuid.New(),
			BuyOrderID:  buyOrder.ID,
			SellOrderID: sellOrder.ID,
			Symbol:      "AAPL",
			Price:       decimal.NewFromFloat(150.0),
			Quantity:    10,
			ExecutedAt:  time.Now().UTC().Truncate(time.Microsecond),
		}

		err := tradeRepo.Create(trade)
		require.NoError(t, err)

		trades, err := tradeRepo.GetByOrderID(buyOrder.ID)
		require.NoError(t, err)
		assert.Len(t, trades, 1)
		assert.Equal(t, trade.ID, trades[0].ID)

		trades, err = tradeRepo.GetByOrderID(sellOrder.ID)
		require.NoError(t, err)
		assert.Len(t, trades, 1)
		assert.Equal(t, trade.ID, trades[0].ID)
	})
}
