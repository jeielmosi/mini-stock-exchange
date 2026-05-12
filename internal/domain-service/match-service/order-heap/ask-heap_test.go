package order_heaps

import (
	"context"
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
		heap:      NewOrderHeap(capacity, lessOrder),
		orderRepo: orderRepo,
		symbol:    symbol,
	}
}

func TestNewAskHeap(t *testing.T) {
	symbol := "AAPL"
	capacity := 10
	var repo repository.OrderRepository
	ah := NewAskHeap(symbol, capacity, repo)
	assert.NotNil(t, ah)
}

func TestAskHeap(t *testing.T) {
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

	t.Run("Pop Fill from Repository", func(t *testing.T) {
		ctx := context.Background()
		db, cleanup, err := repository.SetupTestDB(ctx)
		assert.NoError(t, err)
		defer cleanup()
		repo, err := repository.NewOrderRepository(db)
		assert.NoError(t, err)
		defer repo.Stop()
		brokerRepo, err := repository.NewBrokerRepository(db)
		assert.NoError(t, err)

		brokerID := uuid.New()
		err = brokerRepo.Insert(entity.Broker{ID: brokerID, Name: "broker1"})
		assert.NoError(t, err)

		ah := newAskHeap(symbol, 10, repo)

		order := entity.Order{
			ID:                uuid.New(),
			BrokerID:          brokerID,
			OwnerDoc:          "doc1",
			Symbol:            symbol,
			Type:              entity.Ask,
			Price:             decimal.NewFromInt(100),
			Quantity:          100,
			RemainingQuantity: 100,
			ValidUntil:        now.Add(24 * time.Hour),
			Status:            entity.Pending,
			CreatedAt:         now,
		}
		err = repo.Insert(order)
		assert.NoError(t, err)

		inserted, err := repo.GetByID(order.ID)
		assert.NoError(t, err)
		assert.True(t, inserted.ID == order.ID)

		asks, err := repo.GetAsks(symbol, 10)
		assert.NoError(t, err)
		t.Logf("Asks from repo: %d", len(asks))

		bidOrder := entity.Order{
			ID:                uuid.New(),
			Price:             decimal.NewFromInt(200),
			Quantity:          100,
			RemainingQuantity: 100,
			OwnerDoc:          "doc2",
		}
		res, err := ah.Pop(bidOrder)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		t.Logf("Matches length: %d", len(res.Matches))
		t.Log("Expired: ", res.Expired)
		assert.True(t, len(res.Matches) > 0)
		assert.True(t, res.Matches[0].ID == order.ID)
	})

	t.Run("Pop Empty Repository", func(t *testing.T) {

		ctx := context.Background()
		db, cleanup, err := repository.SetupTestDB(ctx)
		assert.NoError(t, err)
		defer cleanup()
		repo, err := repository.NewOrderRepository(db)
		assert.NoError(t, err)
		defer repo.Stop()

		ah := newAskHeap(symbol, 10, repo)
		_, err = ah.Pop(entity.Order{})

		assert.NoError(t, err)
	})

	t.Run("Push", func(t *testing.T) {
		ctx := context.Background()
		db, cleanup, err := repository.SetupTestDB(ctx)
		assert.NoError(t, err)
		defer cleanup()
		repo, err := repository.NewOrderRepository(db)
		assert.NoError(t, err)
		ah := newAskHeap(symbol, 5, repo)
		order1 := entity.Order{
			ID:                uuid.New(),
			Symbol:            symbol,
			Price:             decimal.NewFromInt(100),
			CreatedAt:         now,
			Type:              entity.Ask,
			ValidUntil:        now.Add(24 * time.Hour),
			Status:            entity.Pending,
			Quantity:          100,
			RemainingQuantity: 100,
		}
		order2 := entity.Order{
			ID:                uuid.New(),
			Symbol:            symbol,
			Price:             decimal.NewFromInt(90),
			CreatedAt:         now,
			Type:              entity.Ask,
			ValidUntil:        now.Add(24 * time.Hour),
			Status:            entity.Pending,
			Quantity:          100,
			RemainingQuantity: 100,
		}

		ah.Push(order1)
		ah.Push(order2)

		o, ok := ah.heap.Peek()
		assert.True(t, ok)
		assert.True(t, o.Price.Equal(decimal.NewFromInt(90)))
	})
	t.Run("Peek", func(t *testing.T) {
		ctx := context.Background()
		db, cleanup, err := repository.SetupTestDB(ctx)
		assert.NoError(t, err)
		defer cleanup()
		repo, err := repository.NewOrderRepository(db)
		assert.NoError(t, err)
		ah := newAskHeap(symbol, 5, repo)
		order1 := entity.Order{
			ID:                uuid.New(),
			Symbol:            symbol,
			Price:             decimal.NewFromInt(100),
			CreatedAt:         now,
			Type:              entity.Ask,
			ValidUntil:        now.Add(24 * time.Hour),
			Status:            entity.Pending,
			Quantity:          100,
			RemainingQuantity: 100,
		}
		order2 := entity.Order{
			ID:                uuid.New(),
			Symbol:            symbol,
			Price:             decimal.NewFromInt(90),
			CreatedAt:         now,
			Type:              entity.Ask,
			ValidUntil:        now.Add(24 * time.Hour),
			Status:            entity.Pending,
			Quantity:          100,
			RemainingQuantity: 100,
		}

		ah.Push(order1)
		ah.Push(order2)

		o, ok := ah.heap.Peek()
		assert.True(t, ok)
		assert.True(t, o.Price.Equal(decimal.NewFromInt(90)))
	})

	t.Run("Push Full Heap and Trigger Fill", func(t *testing.T) {
		ctx := context.Background()
		db, cleanup, err := repository.SetupTestDB(ctx)
		assert.NoError(t, err)
		defer cleanup()
		repo, err := repository.NewOrderRepository(db)
		assert.NoError(t, err)
		defer repo.Stop()
		brokerRepo, err := repository.NewBrokerRepository(db)
		assert.NoError(t, err)

		b0 := uuid.New()
		b1 := uuid.New()
		b2 := uuid.New()
		err = brokerRepo.Insert(entity.Broker{ID: b0, Name: "broker0"})
		assert.NoError(t, err)
		err = brokerRepo.Insert(entity.Broker{ID: b1, Name: "broker1"})
		assert.NoError(t, err)
		err = brokerRepo.Insert(entity.Broker{ID: b2, Name: "broker2"})
		assert.NoError(t, err)

		ah := newAskHeap(symbol, 2, repo)

		// Order in repo that should be picked up by fill
		o0 := entity.Order{
			ID:                uuid.New(),
			BrokerID:          b0,
			OwnerDoc:          "doc0",
			Symbol:            symbol,
			Price:             decimal.NewFromInt(5),
			CreatedAt:         now,
			Type:              entity.Ask,
			ValidUntil:        now.Add(24 * time.Hour),
			Status:            entity.Pending,
			Quantity:          100,
			RemainingQuantity: 100,
		}
		ah.Push(o0)
		assert.True(t, ah.heap.Len() == 1)
		err = repo.Insert(o0)
		assert.NoError(t, err)

		o1 := entity.Order{
			ID:                uuid.New(),
			BrokerID:          b1,
			OwnerDoc:          "doc1",
			Symbol:            symbol,
			Price:             decimal.NewFromInt(10),
			CreatedAt:         now,
			Type:              entity.Ask,
			ValidUntil:        now.Add(24 * time.Hour),
			Status:            entity.Pending,
			Quantity:          100,
			RemainingQuantity: 100,
		}
		ah.Push(o1) // Heap: [5, 10]
		assert.True(t, ah.heap.Len() == 2)
		err = repo.Insert(o1)
		assert.NoError(t, err)

		o2 := entity.Order{
			ID:                uuid.New(),
			BrokerID:          b2,
			OwnerDoc:          "doc2",
			Symbol:            symbol,
			Price:             decimal.NewFromInt(20),
			CreatedAt:         now,
			Type:              entity.Ask,
			ValidUntil:        now.Add(24 * time.Hour),
			Status:            entity.Pending,
			Quantity:          100,
			RemainingQuantity: 100,
		}
		ah.Push(o2) // Heap full
		err = repo.Insert(o2)
		assert.NoError(t, err)

		// Pop should trigger fill because qt(10) < top(20)
		match := entity.Order{
			ID:                uuid.New(),
			Symbol:            symbol,
			Price:             decimal.NewFromInt(200),
			Quantity:          100,
			RemainingQuantity: 100,
			Type:              entity.Bid,
			OwnerDoc:          "match",
			BrokerID:          uuid.New(),
		}

		res, err := ah.Pop(match)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		// The fill should have pushed repoOrder(5) into the heap
		// Heap was [20], now it's [5, 20].
		// Pop returns 5.
		assert.True(t, len(res.Matches) == 1)
		err = repo.Expire([]uuid.UUID{res.Matches[0].ID})
		assert.NoError(t, err)
		t.Log(res.Matches[0].ID == o0.ID, "match: ", res.Matches[0].OwnerDoc, "order: ", o0.OwnerDoc)
		assert.True(t, res.Matches[0].ID == o0.ID)

		res, err = ah.Pop(match)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.True(t, len(res.Matches) == 1)
		err = repo.Expire([]uuid.UUID{res.Matches[0].ID})
		assert.NoError(t, err)
		t.Log(res.Matches[0].ID == o1.ID, "match: ", res.Matches[0].OwnerDoc, "order: ", o1.OwnerDoc)
		assert.True(t, res.Matches[0].ID == o1.ID)

		res, err = ah.Pop(match)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.True(t, len(res.Matches) == 1)
		err = repo.Expire([]uuid.UUID{res.Matches[0].ID})
		assert.NoError(t, err)
		assert.True(t, res.Matches[0].ID == o2.ID)
		t.Log(res.Matches[0].ID == o2.ID, "match: ", res.Matches[0].OwnerDoc, "order: ", o2.OwnerDoc)
	})
}
