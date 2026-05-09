package order_service

import (
	"context"
	"testing"
	"time"

	"mini-stock-exchange/internal/domain"
	"mini-stock-exchange/internal/repository"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func mustUUID() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return id
}

func TestSymbolExecutor_ProcessOrder(t *testing.T) {
	orderRepo := new(repository.MockOrderRepository)
	tradeRepo := new(repository.MockTradeRepository)

	orderRepo.On("GetBids", "AAPL").Return([]domain.Order{}, nil)
	orderRepo.On("GetAsks", "AAPL").Return([]domain.Order{}, nil)

	executor, err := NewExecutor("AAPL", orderRepo, tradeRepo)
	assert.NoError(t, err)
	defer executor.Stop()

	order := domain.Order{
		ID:                mustUUID(),
		Symbol:            "AAPL",
		Price:             decimal.NewFromInt(150),
		Quantity:          10,
		RemainingQuantity: 10,
		Type:              domain.Bid,
	}

	orderRepo.On("Insert", order).Return(nil)

	err = executor.ProcessOrder(context.Background(), order)
	assert.NoError(t, err)
	orderRepo.AssertExpectations(t)
}

func TestSymbolExecutor_MatchMaking(t *testing.T) {
	orderRepo := new(repository.MockOrderRepository)
	tradeRepo := new(repository.MockTradeRepository)

	orderRepo.On("GetBids", "AAPL").Return([]domain.Order{}, nil)
	orderRepo.On("GetAsks", "AAPL").Return([]domain.Order{}, nil)

	executor, err := NewExecutor("AAPL", orderRepo, tradeRepo)
	defer executor.Stop()
	assert.NoError(t, err)

	// 1. Add an Ask order to the heap
	ask := &domain.Order{
		ID:                mustUUID(),
		Symbol:            "AAPL",
		Price:             decimal.NewFromInt(150),
		Quantity:          10,
		RemainingQuantity: 10,
		Type:              domain.Ask,
		OwnerDoc:          "broker1",
		ValidUntil:        time.Now().Add(time.Hour),
	}
	orderRepo.On("Insert", *ask).Return(nil)
	executor.ProcessOrder(context.Background(), *ask)

	// Small sleep to ensure the order is processed by the goroutine
	time.Sleep(100 * time.Millisecond)

	// 2. Add a Bid order that matches the Ask
	bid := &domain.Order{
		ID:                mustUUID(),
		Symbol:            "AAPL",
		Price:             decimal.NewFromInt(150),
		Quantity:          10,
		RemainingQuantity: 10,
		Type:              domain.Bid,
		OwnerDoc:          "broker2",
	}
	orderRepo.On("Insert", *bid).Return(nil)

	// Expectations for the match
	tradeRepo.On("Create", mock.AnythingOfType("domain.Trade")).Return(nil)
	orderRepo.On("Update", mock.AnythingOfType("domain.Order")).Return(nil)

	err = executor.ProcessOrder(context.Background(), *bid)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	tradeRepo.AssertExpectations(t)
}
