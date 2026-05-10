package order_heaps

import (
	"testing"
	"time"

	"mini-stock-exchange/internal/entity"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"mini-stock-exchange/internal/repository"
)

func newAskHeap(symbol string, capacity int, orderRepo repository.OrderRepository) *AskHeap {
	return &AskHeap{
		heap:      newOrderHeap(capacity, lessOrder),
		orderRepo: orderRepo,
		symbol:    symbol,
	}
}

func TestAskHeap(t *testing.T) {
	symbol := "AAPL"
	now := time.Now().UTC()

	t.Run("Success", func(t *testing.T) {
		repo, cleanup := repository.NewMockOrderRepository()
		defer cleanup()
		orders := []entity.Order{
			{
				ID:         uuid.New(),
				Symbol:     symbol,
				Price:      decimal.NewFromInt(150),
				CreatedAt:  now.Add(time.Minute),
				Type:       entity.Ask,
				ValidUntil: now.Add(24 * time.Hour),
			},
			{
				ID:         uuid.New(),
				Symbol:     symbol,
				Price:      decimal.NewFromInt(100),
				CreatedAt:  now,
				Type:       entity.Ask,
				ValidUntil: now.Add(24 * time.Hour),
			},
			{
				ID:         uuid.New(),
				Symbol:     symbol,
				Price:      decimal.NewFromInt(100),
				CreatedAt:  now.Add(-time.Minute),
				Type:       entity.Ask,
				ValidUntil: now.Add(24 * time.Hour),
			},
		}

		ah := newAskHeap(symbol, 10, repo)
		assert.NotNil(t, ah)

		// The heap should be ordered by Price (ascending) and then CreatedAt (ascending)
		// 1st: Price 100, CreatedAt now - 1m
		// 2nd: Price 100, CreatedAt now
		// 3rd: Price 150, CreatedAt now + 1m
		for _, o := range orders {
			ah.Push(o)
		}

		o1, ok := ah.heap.Peek()
		assert.True(t, ok)
		assert.True(t, o1.Price.Equal(decimal.NewFromInt(100)))
		assert.True(t, o1.CreatedAt.Before(now))

		ah.heap.Drop()
		o2, ok := ah.heap.Peek()
		assert.True(t, ok)
		assert.True(t, o2.Price.Equal(decimal.NewFromInt(100)))
		assert.True(t, o2.CreatedAt.Equal(now))

		ah.heap.Drop()
		o3, ok := ah.heap.Peek()
		assert.True(t, ok)
		assert.True(t, o3.Price.Equal(decimal.NewFromInt(150)))
	})

	t.Run("Pop Empty Repository", func(t *testing.T) {
		repo, cleanup := repository.NewMockOrderRepository()
		defer cleanup()

		ah := newAskHeap(symbol, 10, repo)
		_, err := ah.Pop(entity.Order{})

		assert.NoError(t, err)
	})

	t.Run("Push", func(t *testing.T) {
		repo, cleanup := repository.NewMockOrderRepository()
		defer cleanup()
		ah := newAskHeap(symbol, 5, repo)

		order1 := entity.Order{
			ID:         uuid.New(),
			Symbol:     symbol,
			Price:      decimal.NewFromInt(100),
			CreatedAt:  now,
			Type:       entity.Ask,
			ValidUntil: now.Add(24 * time.Hour),
		}
		order2 := entity.Order{
			ID:         uuid.New(),
			Symbol:     symbol,
			Price:      decimal.NewFromInt(90),
			CreatedAt:  now,
			Type:       entity.Ask,
			ValidUntil: now.Add(24 * time.Hour),
		}

		ah.Push(order1)
		ah.Push(order2)

		o, ok := ah.heap.Peek()
		assert.True(t, ok)
		assert.True(t, o.Price.Equal(decimal.NewFromInt(90)))
	})
	t.Run("Peek", func(t *testing.T) {
		repo, cleanup := repository.NewMockOrderRepository()
		defer cleanup()
		ah := newAskHeap(symbol, 5, repo)
		order1 := entity.Order{
			ID:         uuid.New(),
			Symbol:     symbol,
			Price:      decimal.NewFromInt(100),
			CreatedAt:  now,
			Type:       entity.Ask,
			ValidUntil: now.Add(24 * time.Hour),
		}
		order2 := entity.Order{
			ID:         uuid.New(),
			Symbol:     symbol,
			Price:      decimal.NewFromInt(90),
			CreatedAt:  now,
			Type:       entity.Ask,
			ValidUntil: now.Add(24 * time.Hour),
		}

		ah.Push(order1)
		ah.Push(order2)

		o, ok := ah.heap.Peek()
		assert.True(t, ok)
		assert.True(t, o.Price.Equal(decimal.NewFromInt(90)))
	})
}
