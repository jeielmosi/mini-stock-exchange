package match_service

import (
	"mini-stock-exchange/internal/config"
	"mini-stock-exchange/internal/repository"
)

func NewMockExecutor(symbol string, capacity int, orderRepo repository.OrderRepository) Executor {
	config.LoadTest(capacity)
	return NewExecutor(symbol, orderRepo)
}

/*
type MockExecutor struct {
	mock.Mock
}

func (m *MockExecutor) Insert(order entity.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockExecutor) GetByID(id uuid.UUID) (entity.Order, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Order), args.Error(1)
}

func (m *MockExecutor) Update(order entity.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockExecutor) GetBids(symbol string) ([]entity.Order, error) {
	args := m.Called(symbol)
	return args.Get(0).([]entity.Order), args.Error(1)
}

func (m *MockExecutor) GetAsks(symbol string) ([]entity.Order, error) {
	args := m.Called(symbol)
	return args.Get(0).([]entity.Order), args.Error(1)
}

func (m *MockExecutor) ProcessOrder(ctx context.Context, order entity.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockExecutor) Stop() {
	m.Called()
}
*/
