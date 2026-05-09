package repository

import (
	"mini-stock-exchange/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTradeRepository struct {
	mock.Mock
}

func (m *MockTradeRepository) Create(trade domain.Trade) error {
	args := m.Called(trade)
	if args.Get(0) == nil {
		return args.Error(0)
	}
	return args.Error(0)
}

func (m *MockTradeRepository) GetByOrderID(orderID uuid.UUID) ([]domain.Trade, error) {
	args := m.Called(orderID)
	return args.Get(0).([]domain.Trade), args.Error(1)
}
