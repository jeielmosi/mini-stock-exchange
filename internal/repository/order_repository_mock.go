package repository

import (
	"mini-stock-exchange/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Insert(order domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetByID(id uuid.UUID) (domain.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return domain.Order{}, args.Error(1)
	}
	return args.Get(0).(domain.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(order domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetBids(symbol string) ([]domain.Order, error) {
	args := m.Called(symbol)
	return args.Get(0).([]domain.Order), args.Error(1)
}

func (m *MockOrderRepository) GetAsks(symbol string) ([]domain.Order, error) {
	args := m.Called(symbol)
	return args.Get(0).([]domain.Order), args.Error(1)
}
