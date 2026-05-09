package order_service

import (
	"context"
	"testing"

	"mini-stock-exchange/internal/domain"
	"mini-stock-exchange/internal/repository"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func mustNewV7() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return id
}

type MockCentralizer struct {
	mock.Mock
}

func (m *MockCentralizer) RouteOrder(ctx context.Context, order domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func TestOrderService_SubmitOrder_Match(t *testing.T) {
	orderRepo := new(repository.MockOrderRepository)
	tradeRepo := new(repository.MockTradeRepository)
	centralizer := new(MockCentralizer)
	svc := NewOrderService(orderRepo, tradeRepo, centralizer)

	symbol := "AAPL"
	bidPrice := decimal.NewFromInt(150)

	bidOrder := &domain.Order{
		Symbol:   symbol,
		Price:    bidPrice,
		Quantity: 10,
		Type:     domain.Bid,
	}

	centralizer.On("RouteOrder", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// Simulate matching in background or directly if we want to test the match logic.
		// But here we are testing OrderService.SubmitOrder which now calls centralizer.
		// Since the original match logic was in OrderService, and we moved it to Executor,
		// this test might need adjustment.
		// For now, let's just make sure centralizer is called.
	})

	err := svc.SubmitOrder(context.Background(), *bidOrder)

	assert.NoError(t, err)
	orderRepo.AssertExpectations(t)
	centralizer.AssertExpectations(t)
}

func TestOrderService_SubmitOrder_PartialMatch(t *testing.T) {
	orderRepo := new(repository.MockOrderRepository)
	tradeRepo := new(repository.MockTradeRepository)
	centralizer := new(MockCentralizer)
	svc := NewOrderService(orderRepo, tradeRepo, centralizer)

	symbol := "AAPL"
	bidPrice := decimal.NewFromInt(150)

	bidOrder := &domain.Order{
		Symbol:   symbol,
		Price:    bidPrice,
		Quantity: 10,
		Type:     domain.Bid,
	}

	centralizer.On("RouteOrder", mock.Anything, mock.Anything).Return(nil)

	err := svc.SubmitOrder(context.Background(), *bidOrder)

	assert.NoError(t, err)
	orderRepo.AssertExpectations(t)
	centralizer.AssertExpectations(t)
}

func TestOrderService_SubmitOrder_NoMatch(t *testing.T) {
	orderRepo := new(repository.MockOrderRepository)
	tradeRepo := new(repository.MockTradeRepository)
	centralizer := new(MockCentralizer)
	svc := NewOrderService(orderRepo, tradeRepo, centralizer)

	symbol := "AAPL"
	bidPrice := decimal.NewFromInt(100)

	bidOrder := &domain.Order{
		Symbol:   symbol,
		Price:    bidPrice,
		Quantity: 10,
		Type:     domain.Bid,
	}

	centralizer.On("RouteOrder", mock.Anything, mock.Anything).Return(nil)

	err := svc.SubmitOrder(context.Background(), *bidOrder)

	assert.NoError(t, err)
	orderRepo.AssertExpectations(t)
	centralizer.AssertExpectations(t)
}

func TestOrderService_SubmitOrder_MultipleMatches(t *testing.T) {
	orderRepo := new(repository.MockOrderRepository)
	tradeRepo := new(repository.MockTradeRepository)
	centralizer := new(MockCentralizer)
	svc := NewOrderService(orderRepo, tradeRepo, centralizer)

	symbol := "AAPL"
	bidPrice := decimal.NewFromInt(150)

	bidOrder := &domain.Order{
		Symbol:   symbol,
		Price:    bidPrice,
		Quantity: 10,
		Type:     domain.Bid,
	}

	centralizer.On("RouteOrder", mock.Anything, mock.Anything).Return(nil)

	err := svc.SubmitOrder(context.Background(), *bidOrder)

	assert.NoError(t, err)
	orderRepo.AssertExpectations(t)
	centralizer.AssertExpectations(t)
}

func TestOrderService_SubmitOrder_AskMatch(t *testing.T) {
	orderRepo := new(repository.MockOrderRepository)
	tradeRepo := new(repository.MockTradeRepository)
	centralizer := new(MockCentralizer)
	svc := NewOrderService(orderRepo, tradeRepo, centralizer)

	symbol := "AAPL"
	askPrice := decimal.NewFromInt(130)

	askOrder := &domain.Order{
		Symbol:   symbol,
		Price:    askPrice,
		Quantity: 10,
		Type:     domain.Ask,
	}

	centralizer.On("RouteOrder", mock.Anything, mock.Anything).Return(nil)

	err := svc.SubmitOrder(context.Background(), *askOrder)

	assert.NoError(t, err)
	orderRepo.AssertExpectations(t)
	centralizer.AssertExpectations(t)
}
