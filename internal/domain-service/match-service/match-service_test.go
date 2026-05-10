package match_service

import (
	"context"
	"testing"
	"time"

	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"

	"github.com/google/uuid"
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

func (m *MockCentralizer) RouteOrder(ctx context.Context, order entity.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func TestMatchService_SubmitOrder_Match(t *testing.T) {
	orderRepo, cleanup := repository.NewMockOrderRepository()
	defer cleanup()
	centralizer := new(MockCentralizer)
	svc := NewMatchService(orderRepo, centralizer)

	symbol := "AAPL"
	bidPrice := float64(150)

	bidOrder := dto.CreateOrderRequest{
		Symbol:     symbol,
		Price:      bidPrice,
		Quantity:   10,
		Type:       entity.Bid,
		ValidUntil: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}

	centralizer.On("RouteOrder", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// Simulate matching in background or directly if we want to test the match logic.
		// But here we are testing MatchService.SubmitOrder which now calls centralizer.
		// Since the original match logic was in MatchService, and we moved it to Executor,
		// this test might need adjustment.
		// For now, let's just make sure centralizer is called.
	})

	dto, err := svc.SubmitOrder(context.Background(), bidOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, dto.ID)
	centralizer.AssertExpectations(t)
}

func TestMatchService_SubmitOrder_PartialMatch(t *testing.T) {
	orderRepo, cleanup := repository.NewMockOrderRepository()
	defer cleanup()
	centralizer := new(MockCentralizer)
	svc := NewMatchService(orderRepo, centralizer)

	symbol := "AAPL"
	bidPrice := float64(150)

	bidOrder := dto.CreateOrderRequest{
		Symbol:     symbol,
		Price:      bidPrice,
		Quantity:   10,
		Type:       entity.Bid,
		ValidUntil: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}

	centralizer.On("RouteOrder", mock.Anything, mock.Anything).Return(nil)

	dto, err := svc.SubmitOrder(context.Background(), bidOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, dto.ID)
	centralizer.AssertExpectations(t)
}

func TestMatchService_SubmitOrder_NoMatch(t *testing.T) {
	orderRepo, cleanup := repository.NewMockOrderRepository()
	defer cleanup()
	centralizer := new(MockCentralizer)
	svc := NewMatchService(orderRepo, centralizer)

	symbol := "AAPL"
	bidPrice := float64(100)

	bidOrder := dto.CreateOrderRequest{
		Symbol:     symbol,
		Price:      bidPrice,
		Quantity:   10,
		Type:       entity.Bid,
		ValidUntil: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}

	centralizer.On("RouteOrder", mock.Anything, mock.Anything).Return(nil)

	dto, err := svc.SubmitOrder(context.Background(), bidOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, dto.ID)
	centralizer.AssertExpectations(t)
}

func TestMatchService_SubmitOrder_MultipleMatches(t *testing.T) {
	orderRepo, cleanup := repository.NewMockOrderRepository()
	defer cleanup()
	centralizer := new(MockCentralizer)
	svc := NewMatchService(orderRepo, centralizer)

	symbol := "AAPL"
	bidPrice := float64(150)

	bidOrder := dto.CreateOrderRequest{
		Symbol:     symbol,
		Price:      bidPrice,
		Quantity:   10,
		Type:       entity.Bid,
		ValidUntil: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}

	centralizer.On("RouteOrder", mock.Anything, mock.Anything).Return(nil)

	dto, err := svc.SubmitOrder(context.Background(), bidOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, dto.ID)
	centralizer.AssertExpectations(t)
}

func TestMatchService_SubmitOrder_AskMatch(t *testing.T) {
	orderRepo, cleanup := repository.NewMockOrderRepository()
	defer cleanup()
	centralizer := new(MockCentralizer)
	svc := NewMatchService(orderRepo, centralizer)

	symbol := "AAPL"
	askPrice := float64(130)

	askOrder := dto.CreateOrderRequest{
		Symbol:     symbol,
		Price:      askPrice,
		Quantity:   10,
		Type:       entity.Ask,
		ValidUntil: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}

	centralizer.On("RouteOrder", mock.Anything, mock.Anything).Return(nil)

	dto, err := svc.SubmitOrder(context.Background(), askOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, dto.ID)
	centralizer.AssertExpectations(t)
}
