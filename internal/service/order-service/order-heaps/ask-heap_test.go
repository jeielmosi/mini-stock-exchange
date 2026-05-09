package order_heaps

import (
	"errors"
	"testing"
	"time"

	"mini-stock-exchange/internal/domain"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"mini-stock-exchange/internal/repository"
)

func TestNewAskHeap(t *testing.T) {
	symbol := "AAPL"
	now := time.Now().UTC()

	t.Run("Success", func(t *testing.T) {
		repo := new(repository.MockOrderRepository)
		orders := []domain.Order{
			{
				ID:        uuid.New(),
				Symbol:    symbol,
				Price:     decimal.NewFromInt(150),
				CreatedAt: now.Add(time.Minute),
				Type:      domain.Ask,
			},
			{
				ID:        uuid.New(),
				Symbol:    symbol,
				Price:     decimal.NewFromInt(100),
				CreatedAt: now,
				Type:      domain.Ask,
			},
			{
				ID:        uuid.New(),
				Symbol:    symbol,
				Price:     decimal.NewFromInt(100),
				CreatedAt: now.Add(-time.Minute),
				Type:      domain.Ask,
			},
		}
		repo.On("GetAsks", symbol).Return(orders, nil)

		heap, err := NewAskHeap(symbol, repo)

		assert.NoError(t, err)
		assert.NotNil(t, heap)

		// The heap should be ordered by Price (ascending) and then CreatedAt (ascending)
		// 1st: Price 100, CreatedAt now - 1m
		// 2nd: Price 100, CreatedAt now
		// 3rd: Price 150, CreatedAt now + 1m

		o1, ok := heap.Peek()
		assert.True(t, ok)
		assert.True(t, o1.Price.Equal(decimal.NewFromInt(100)))
		assert.True(t, o1.CreatedAt.Before(now))

		heap.Pop()
		o2, ok := heap.Peek()
		assert.True(t, ok)
		assert.True(t, o2.Price.Equal(decimal.NewFromInt(100)))
		assert.True(t, o2.CreatedAt.Equal(now))

		heap.Pop()
		o3, ok := heap.Peek()
		assert.True(t, ok)
		assert.True(t, o3.Price.Equal(decimal.NewFromInt(150)))

		repo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repo := new(repository.MockOrderRepository)
		repo.On("GetAsks", symbol).Return([]domain.Order{}, errors.New("db error"))

		heap, err := NewAskHeap(symbol, repo)

		assert.Error(t, err)
		assert.Nil(t, heap)
		assert.Equal(t, "db error", err.Error())

		repo.AssertExpectations(t)
	})
}
