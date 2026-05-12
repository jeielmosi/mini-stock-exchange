package match_service

import (
	"context"
	"testing"
	"time"

	"mini-stock-exchange/internal/config"
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"
	"mini-stock-exchange/internal/usecase"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

	configEnv(10)
	executor := newExecutor("AAPL", orderRepo)
	defer executor.Stop()

	order := entity.Order{
		ID:                mustUUID(),
		Symbol:            "AAPL",
		Price:             decimal.NewFromInt(150),
		Quantity:          10,
		RemainingQuantity: 10,
		Type:              entity.Bid,
		Status:            entity.Pending,
	}

	err = executor.ProcessOrder(context.Background(), order)
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

	configEnv(10)
	executor := newExecutor("AAPL", orderRepo)
	defer executor.Stop()

	// 1. Add an Ask order to the heap
	ask := &entity.Order{
		ID:                mustUUID(),
		Symbol:            "AAPL",
		Price:             decimal.NewFromInt(150),
		Quantity:          10,
		RemainingQuantity: 10,
		Type:              entity.Ask,
		Status:            entity.Pending,
		OwnerDoc:          "broker1",
		ValidUntil:        time.Now().Add(time.Hour),
	}
	executor.ProcessOrder(context.Background(), *ask)

	// Small sleep to ensure the order is processed by the goroutine
	time.Sleep(100 * time.Millisecond)

	// 2. Add a Bid order that matches the Ask
	bid := &entity.Order{
		ID:                mustUUID(),
		Symbol:            "AAPL",
		Price:             decimal.NewFromInt(150),
		Quantity:          10,
		RemainingQuantity: 10,
		Status:            entity.Pending,
		Type:              entity.Bid,
		OwnerDoc:          "broker2",
	}

	err = executor.ProcessOrder(context.Background(), *bid)
	if err != nil {
		t.Log(err.Error())
	}
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
}
