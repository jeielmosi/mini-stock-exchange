package service

import (
	"context"
	"testing"
	"time"

	"mini-stock-exchange/internal/domain"

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

type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) Create(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepo) GetByID(id uuid.UUID) (*domain.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepo) Update(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepo) FindMatches(order domain.Order) ([]domain.Order, error) {
	args := m.Called(order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Order), args.Error(1)
}

func (m *MockOrderRepo) UpdateRemainingQuantity(id uuid.UUID, quantity int, status domain.OrderStatus) error {
	args := m.Called(id, quantity, status)
	return args.Error(0)
}

type MockTradeRepo struct {
	mock.Mock
}

func (m *MockTradeRepo) Create(trade *domain.Trade) error {
	args := m.Called(trade)
	return args.Error(0)
}

func (m *MockTradeRepo) GetByOrderID(orderID uuid.UUID) ([]*domain.Trade, error) {
	args := m.Called(orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Trade), args.Error(1)
}

func TestOrderService_SubmitOrder_Match(t *testing.T) {
	orderRepo := new(MockOrderRepo)
	tradeRepo := new(MockTradeRepo)
	svc := NewOrderService(orderRepo, tradeRepo)

	symbol := "AAPL"
	bidPrice := decimal.NewFromInt(150)
	askPrice := decimal.NewFromInt(140)

	bidOrder := &domain.Order{
		Symbol:   symbol,
		Price:    bidPrice,
		Quantity: 10,
		Type:     domain.Bid,
	}

	askOrder := &domain.Order{
		ID:                mustNewV7(),
		Symbol:            symbol,
		Price:             askPrice,
		Quantity:          10,
		RemainingQuantity: 10,
		Type:              domain.Ask,
		Status:            domain.Pending,
	}

	orderRepo.On("Create", mock.Anything).Return(nil)
	orderRepo.On("FindMatches", mock.Anything).Return([]domain.Order{*askOrder}, nil).Once()
	tradeRepo.On("Create", mock.Anything).Return(nil)
	orderRepo.On("Update", mock.Anything).Return(nil).Times(2)

	err := svc.SubmitOrder(context.Background(), bidOrder)

	time.Sleep(10 * time.Millisecond)

	assert.NoError(t, err)
	orderRepo.AssertExpectations(t)
	tradeRepo.AssertExpectations(t)
}

func TestOrderService_SubmitOrder_PartialMatch(t *testing.T) {
	orderRepo := new(MockOrderRepo)
	tradeRepo := new(MockTradeRepo)
	svc := NewOrderService(orderRepo, tradeRepo)

	symbol := "AAPL"
	bidPrice := decimal.NewFromInt(150)
	askPrice := decimal.NewFromInt(140)

	bidOrder := &domain.Order{
		Symbol:   symbol,
		Price:    bidPrice,
		Quantity: 10,
		Type:     domain.Bid,
	}

	askOrder := &domain.Order{
		ID:                mustNewV7(),
		Symbol:            symbol,
		Price:             askPrice,
		Quantity:          5,
		RemainingQuantity: 5,
		Type:              domain.Ask,
		Status:            domain.Pending,
	}

	orderRepo.On("Create", mock.Anything).Return(nil)
	orderRepo.On("FindMatches", mock.Anything).Return([]domain.Order{*askOrder}, nil).Once()
	tradeRepo.On("Create", mock.Anything).Return(nil)
	orderRepo.On("Update", mock.Anything).Return(nil).Times(2)
	orderRepo.On("FindMatches", mock.Anything).Return([]domain.Order{}, nil).Once()

	err := svc.SubmitOrder(context.Background(), bidOrder)

	time.Sleep(10 * time.Millisecond)

	assert.NoError(t, err)
	assert.Equal(t, 5, bidOrder.RemainingQuantity)
	assert.Equal(t, domain.Partial, bidOrder.Status)
	orderRepo.AssertExpectations(t)
	tradeRepo.AssertExpectations(t)
}

func TestOrderService_SubmitOrder_NoMatch(t *testing.T) {
	orderRepo := new(MockOrderRepo)
	tradeRepo := new(MockTradeRepo)
	svc := NewOrderService(orderRepo, tradeRepo)

	symbol := "AAPL"
	bidPrice := decimal.NewFromInt(100)

	bidOrder := &domain.Order{
		Symbol:   symbol,
		Price:    bidPrice,
		Quantity: 10,
		Type:     domain.Bid,
	}

	orderRepo.On("Create", mock.Anything).Return(nil)
	orderRepo.On("FindMatches", mock.Anything).Return([]domain.Order{}, nil).Once()

	err := svc.SubmitOrder(context.Background(), bidOrder)

	time.Sleep(10 * time.Millisecond)

	assert.NoError(t, err)
	assert.Equal(t, 10, bidOrder.RemainingQuantity)
	assert.Equal(t, domain.Pending, bidOrder.Status)
	orderRepo.AssertExpectations(t)
}

func TestOrderService_SubmitOrder_MultipleMatches(t *testing.T) {
	orderRepo := new(MockOrderRepo)
	tradeRepo := new(MockTradeRepo)
	svc := NewOrderService(orderRepo, tradeRepo)

	symbol := "AAPL"
	bidPrice := decimal.NewFromInt(150)
	askPrice := decimal.NewFromInt(140)

	bidOrder := &domain.Order{
		Symbol:   symbol,
		Price:    bidPrice,
		Quantity: 20,
		Type:     domain.Bid,
	}

	askOrder1 := &domain.Order{
		ID:                mustNewV7(),
		Symbol:            symbol,
		Price:             askPrice,
		Quantity:          10,
		RemainingQuantity: 10,
		Type:              domain.Ask,
		Status:            domain.Pending,
	}

	askOrder2 := &domain.Order{
		ID:                mustNewV7(),
		Symbol:            symbol,
		Price:             askPrice,
		Quantity:          10,
		RemainingQuantity: 10,
		Type:              domain.Ask,
		Status:            domain.Pending,
	}

	orderRepo.On("Create", mock.Anything).Return(nil).Once()
	orderRepo.On("FindMatches", mock.Anything).Return([]domain.Order{*askOrder1}, nil).Once()
	tradeRepo.On("Create", mock.Anything).Return(nil).Times(2)
	orderRepo.On("Update", mock.Anything).Return(nil).Times(4)

	orderRepo.On("FindMatches", mock.Anything).Return([]domain.Order{*askOrder2}, nil).Once()

	err := svc.SubmitOrder(context.Background(), bidOrder)

	time.Sleep(10 * time.Millisecond)

	assert.NoError(t, err)
	assert.Equal(t, 0, bidOrder.RemainingQuantity)
	assert.Equal(t, domain.Filled, bidOrder.Status)
	orderRepo.AssertExpectations(t)
	tradeRepo.AssertExpectations(t)
}

func TestOrderService_SubmitOrder_AskMatch(t *testing.T) {
	orderRepo := new(MockOrderRepo)
	tradeRepo := new(MockTradeRepo)
	svc := NewOrderService(orderRepo, tradeRepo)

	symbol := "AAPL"
	bidPrice := decimal.NewFromInt(140)
	askPrice := decimal.NewFromInt(130)

	askOrder := &domain.Order{
		Symbol:   symbol,
		Price:    askPrice,
		Quantity: 10,
		Type:     domain.Ask,
	}

	bidOrder := &domain.Order{
		ID:                mustNewV7(),
		Symbol:            symbol,
		Price:             bidPrice,
		Quantity:          10,
		RemainingQuantity: 10,
		Type:              domain.Bid,
		Status:            domain.Pending,
	}

	orderRepo.On("Create", mock.Anything).Return(nil)
	orderRepo.On("FindMatches", mock.Anything).Return([]domain.Order{*bidOrder}, nil).Once()
	tradeRepo.On("Create", mock.Anything).Return(nil)
	orderRepo.On("Update", mock.Anything).Return(nil).Times(2)

	err := svc.SubmitOrder(context.Background(), askOrder)

	time.Sleep(10 * time.Millisecond)

	assert.NoError(t, err)
	assert.Equal(t, 0, askOrder.RemainingQuantity)
	assert.Equal(t, domain.Filled, askOrder.Status)
	orderRepo.AssertExpectations(t)
	tradeRepo.AssertExpectations(t)
}
