package order_heaps

import (
	"errors"
	"testing"
	"time"

	"mini-stock-exchange/internal/domain"
	"mini-stock-exchange/internal/repository"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewBidHeap(t *testing.T) {
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
				Type:      domain.Bid,
			},
			{
				ID:        uuid.New(),
				Symbol:    symbol,
				Price:     decimal.NewFromInt(200),
				CreatedAt: now,
				Type:      domain.Bid,
			},
			{
				ID:        uuid.New(),
				Symbol:    symbol,
				Price:     decimal.NewFromInt(200),
				CreatedAt: now.Add(-time.Minute),
				Type:      domain.Bid,
			},
		}
		// Note: In bid-heap.go, it incorrectly calls GetAsks instead of GetBids.
		// But I should test what's there or fix it.
		// Actually, looking at bid-heap.go:
		// 27: 	orders, err := orderRepo.GetAsks(symbol)
		// That looks like a bug in bid-heap.go.
		// I'll write the test to expect GetAsks for now as it's what's implemented,
		// or maybe I should fix it first.
		// Let's see what GetAsks returns.
		repo.On("GetBids", symbol).Return(orders, nil)

		heap, err := NewBidHeap(symbol, repo)

		assert.NoError(t, err)
		assert.NotNil(t, heap)

		// The heap should be ordered by Price (descending) and then CreatedAt (ascending)
		// 1st: Price 200, CreatedAt now - 1m
		// 2nd: Price 200, CreatedAt now
		// 3rd: Price 150, CreatedAt now + 1m

		o1, ok := heap.Peek()
		assert.True(t, ok)
		assert.True(t, o1.Price.Equal(decimal.NewFromInt(200)))
		assert.True(t, o1.CreatedAt.Before(now))

		heap.Pop()
		o2, ok := heap.Peek()
		assert.True(t, ok)
		assert.True(t, o2.Price.Equal(decimal.NewFromInt(200)))
		assert.True(t, o2.CreatedAt.Equal(now))

		heap.Pop()
		o3, ok := heap.Peek()
		assert.True(t, ok)
		assert.True(t, o3.Price.Equal(decimal.NewFromInt(150)))

		repo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repo := new(repository.MockOrderRepository)
		repo.On("GetBids", symbol).Return([]domain.Order{}, errors.New("db error"))

		heap, err := NewBidHeap(symbol, repo)

		assert.Error(t, err)
		assert.Nil(t, heap)
		assert.Equal(t, "db error", err.Error())

		repo.AssertExpectations(t)
	})
}
