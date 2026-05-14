package match_service

import (
	"context"
	"testing"
	"time"

	"mini-stock-exchange/internal/config"
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"
	"mini-stock-exchange/internal/usecase"

	"math/big"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustUUID() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return id
}

func configEnv(capacity int) {
	config.LoadTest(capacity)
}

func newExecutor(symbol string, orderRepo repository.OrderRepository) Executor {
	return NewExecutor(symbol, orderRepo, usecase.NewOrderMatchUsecase(), usecase.NewCreateTradeUsecase())
}

func TestSymbolExecutor_ProcessOrder(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := repository.SetupTestDB(ctx)
	assert.NoError(t, err)
	defer cleanup()
	orderRepo, err := repository.NewOrderRepository(db)
	require.NoError(t, err)
	defer orderRepo.Stop()
	brokerRepo, err := repository.NewBrokerRepository(db)
	require.NoError(t, err)
	brokerID := uuid.New()
	err = brokerRepo.Insert(entity.Broker{ID: brokerID, Name: "broker1"})
	require.NoError(t, err)

	configEnv(10)
	executor := newExecutor("AAPL", orderRepo)
	defer executor.Stop()

	order := entity.Order{
		ID:                mustUUID(),
		Symbol:            "AAPL",
		Price:             new(big.Rat).SetInt64(150),
		Quantity:          10,
		RemainingQuantity: 10,
		Type:              entity.Bid,
		Status:            entity.Pending,
		BrokerID:          brokerID,
	}

	// Insert must happen before ProcessOrder for FK constraints
	err = orderRepo.Insert(order)
	require.NoError(t, err)
	err = executor.ProcessOrder(context.Background(), &order)
	if err != nil {
		t.Log(err.Error())
	}
	assert.NoError(t, err)
}

func TestSymbolExecutor_MatchMaking(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := repository.SetupTestDB(ctx)
	assert.NoError(t, err)
	defer cleanup()
	orderRepo, err := repository.NewOrderRepository(db)
	assert.NoError(t, err)
	defer orderRepo.Stop()
	brokerRepo, err := repository.NewBrokerRepository(db)
	require.NoError(t, err)
	b1 := uuid.New()
	b2 := uuid.New()
	err = brokerRepo.Insert(entity.Broker{ID: b1, Name: "broker1"})
	require.NoError(t, err)
	err = brokerRepo.Insert(entity.Broker{ID: b2, Name: "broker2"})
	require.NoError(t, err)

	configEnv(10)
	executor := newExecutor("AAPL", orderRepo)
	defer executor.Stop()

	// 1. Add an Ask order - must be in DB first for FK constraints
	ask := &entity.Order{
		ID:                mustUUID(),
		Symbol:            "AAPL",
		Price:             new(big.Rat).SetInt64(150),
		Quantity:          10,
		RemainingQuantity: 10,
		Type:              entity.Ask,
		Status:            entity.Pending,
		OwnerDoc:          "broker1",
		ValidUntil:        time.Now().Add(time.Hour),
		BrokerID:          b1,
	}
	err = orderRepo.Insert(*ask)
	require.NoError(t, err)
	executor.ProcessOrder(context.Background(), ask)

	// Small sleep to ensure the order is processed by the goroutine
	time.Sleep(100 * time.Millisecond)

	// 2. Add a Bid order that matches the Ask
	bid := &entity.Order{
		ID:                mustUUID(),
		Symbol:            "AAPL",
		Price:             new(big.Rat).SetInt64(150),
		Quantity:          10,
		RemainingQuantity: 10,
		Status:            entity.Pending,
		Type:              entity.Bid,
		OwnerDoc:          "broker2",
		BrokerID:          b2,
		ValidUntil:        time.Now().Add(time.Hour),
	}

	err = orderRepo.Insert(*bid)
	require.NoError(t, err)
	err = executor.ProcessOrder(context.Background(), bid)
	if err != nil {
		t.Log(err.Error())
	}
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
}
