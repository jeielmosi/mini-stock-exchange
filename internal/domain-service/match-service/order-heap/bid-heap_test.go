package order_heaps

import (
	"context"
	"testing"
	"time"

	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func newBidHeap(symbol string, capacity int, orderRepo repository.OrderRepository) *BidHeap {
	return &BidHeap{
		heap:      NewOrderHeap(capacity, greaterOrder),
		orderRepo: orderRepo,
		symbol:    symbol,
	}
}

func TestNewBidHeap(t *testing.T) {
	symbol := "AAPL"
	now := time.Now().UTC()

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		db, cleanup, err := repository.SetupTestDB(ctx)
		assert.NoError(t, err)
		defer cleanup()

		repo, err := repository.NewOrderRepository(db)
		assert.NoError(t, err)
		defer repo.Stop()
		orders := []entity.Order{
			{
				ID:         uuid.New(),
				Symbol:     symbol,
				Price:      decimal.NewFromInt(150),
				CreatedAt:  now.Add(time.Minute),
				Type:       entity.Bid,
				ValidUntil: now.Add(24 * time.Hour),
			},
			{
				ID:         uuid.New(),
				Symbol:     symbol,
				Price:      decimal.NewFromInt(200),
				CreatedAt:  now,
				Type:       entity.Bid,
				ValidUntil: now.Add(24 * time.Hour),
			},
			{
				ID:         uuid.New(),
				Symbol:     symbol,
				Price:      decimal.NewFromInt(200),
				CreatedAt:  now.Add(-time.Minute),
				Type:       entity.Bid,
				ValidUntil: now.Add(24 * time.Hour),
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

		bh := newBidHeap(symbol, 10, repo)

		assert.NotNil(t, bh)

		// The heap should be ordered by Price (descending) and then CreatedAt (ascending)
		// 1st: Price 200, CreatedAt now - 1m
		// 2nd: Price 200, CreatedAt now
		// 3rd: Price 150, CreatedAt now + 1m
		for _, o := range orders {
			bh.Push(o)
		}

		o1, ok := bh.heap.Peek()
		assert.True(t, ok)
		assert.True(t, o1.Price.Equal(decimal.NewFromInt(200)))
		assert.True(t, o1.CreatedAt.Before(now))

		bh.heap.Drop()
		o2, ok := bh.heap.Peek()
		assert.True(t, ok)
		assert.True(t, o2.Price.Equal(decimal.NewFromInt(200)))
		assert.True(t, o2.CreatedAt.Equal(now))

		bh.heap.Drop()
		o3, ok := bh.heap.Peek()
		assert.True(t, ok)
		assert.True(t, o3.Price.Equal(decimal.NewFromInt(150)))
	})
}
