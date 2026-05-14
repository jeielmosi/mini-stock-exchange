package repository

import (
	"context"
	"testing"
	"time"

	"mini-stock-exchange/internal/entity"

	"github.com/google/uuid"
	"math/big"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTradeRepository_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer cleanup()

	repo := &tradeRepository{db: db}

	id := uuid.New()
	_, err = repo.GetByID(id)
	assert.Error(t, err)
}

func TestTradeRepository_GetByID_Found(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer cleanup()

	repo := &tradeRepository{db: db}
	orderRepo, err := NewOrderRepository(db)
	require.NoError(t, err)
	brokerRepo, err := NewBrokerRepository(db)
	require.NoError(t, err)

	broker1ID := uuid.New()
	broker2ID := uuid.New()
	err = brokerRepo.Insert(entity.Broker{ID: broker1ID, Name: "broker1"})
	require.NoError(t, err)
	err = brokerRepo.Insert(entity.Broker{ID: broker2ID, Name: "broker2"})
	require.NoError(t, err)

	// Create orders first to satisfy foreign key constraints
	buyOrderID := uuid.New()
	sellOrderID := uuid.New()
	symbol := "AAPL"
	now := time.Now()

	buyOrder := &entity.Order{
		ID:                buyOrderID,
		BrokerID:          broker1ID,
		OwnerDoc:          "doc1",
		Type:              entity.Bid,
		Symbol:            symbol,
		Price:             new(big.Rat).SetFloat64(150.0),
		Quantity:          100,
		RemainingQuantity: 100,
		ValidUntil:        now.Add(time.Hour),
		Status:            entity.Pending,
		CreatedAt:         now,
	}
	sellOrder := &entity.Order{
		ID:                sellOrderID,
		BrokerID:          broker2ID,
		OwnerDoc:          "doc2",
		Type:              entity.Ask,
		Symbol:            symbol,
		Price:             new(big.Rat).SetFloat64(150.0),
		Quantity:          100,
		RemainingQuantity: 100,
		ValidUntil:        now.Add(time.Hour),
		Status:            entity.Pending,
		CreatedAt:         now,
	}

	err = orderRepo.Insert(*buyOrder)
	require.NoError(t, err)
	err = orderRepo.Insert(*sellOrder)
	require.NoError(t, err)

	// Insert a trade
	tradeID := uuid.New()
	tradePrice := new(big.Rat).SetFloat64(150.0)
	tradeQuantity := 100
	executedAt := now

	query := `INSERT INTO trades (id, symbol, price, quantity, executed_at, buy_order_id, sell_order_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.Exec(query, tradeID, symbol, tradePrice.FloatString(8), tradeQuantity, executedAt, buyOrderID, sellOrderID)
	require.NoError(t, err)

	// Get the trade
	retrievedTrade, err := repo.GetByID(tradeID)
	require.NoError(t, err)

	// Assertions
	assert.Equal(t, tradeID, retrievedTrade.ID)
	assert.Equal(t, symbol, retrievedTrade.Symbol)
	assert.True(t, tradePrice.Cmp(retrievedTrade.Price) == 0)
	assert.Equal(t, tradeQuantity, retrievedTrade.Quantity)
	// Use tolerance for time comparison since PostgreSQL truncates to microseconds
	assert.True(t, executedAt.Sub(retrievedTrade.ExecutedAt).Abs() < time.Microsecond)
	assert.Equal(t, buyOrderID, retrievedTrade.BuyOrderID)
	assert.Equal(t, sellOrderID, retrievedTrade.SellOrderID)
}
