package order_heaps

import (
	"testing"
	"time"

	"mini-stock-exchange/internal/entity"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func newTestHeap(cmp func(a *entity.Order, b *entity.Order) bool) *orderHeap {
	return &orderHeap{
		heap: []entity.Order{},
		cmp:  cmp,
	}
}

func TestOrderHeap_PopBack(t *testing.T) {
	cmp := func(a, b *entity.Order) bool {
		return a.Price.GreaterThan(b.Price)
	}
	oh := newTestHeap(cmp)

	// PopBack on empty heap
	_, ok := oh.PopBack()
	assert.False(t, ok)

	order := entity.Order{Price: decimal.NewFromInt(100), ID: uuid.New()}
	oh.heap = append(oh.heap, order)

	popped, ok := oh.PopBack()
	assert.True(t, ok)
	assert.Equal(t, order.ID, popped.ID)
	assert.Equal(t, 0, oh.Len())
}

func TestOrderHeap_PushPopPeek(t *testing.T) {
	// Max heap by price
	cmp := func(a, b *entity.Order) bool {
		return a.Price.GreaterThan(b.Price)
	}
	oh := newTestHeap(cmp)

	orders := []entity.Order{
		{Price: decimal.NewFromInt(100), ID: uuid.New()},
		{Price: decimal.NewFromInt(200), ID: uuid.New()},
		{Price: decimal.NewFromInt(150), ID: uuid.New()},
	}

	for _, o := range orders {
		oh.Push(o)
	}

	assert.Equal(t, 3, oh.Len())

	peeked, ok := oh.Peek()
	assert.True(t, ok)
	assert.Equal(t, decimal.NewFromInt(200), peeked.Price)

	oh.Drop()
	peeked, ok = oh.Peek()
	assert.True(t, ok)
	assert.Equal(t, decimal.NewFromInt(150), peeked.Price)

	oh.Drop()
	peeked, ok = oh.Peek()
	assert.True(t, ok)
	assert.Equal(t, decimal.NewFromInt(100), peeked.Price)

	oh.Drop()
	_, ok = oh.Peek()
	assert.False(t, ok)
	assert.Equal(t, 0, oh.Len())
}

func TestCouldMatch(t *testing.T) {
	bidPrice100 := decimal.NewFromInt(100)
	bidPrice110 := decimal.NewFromInt(110)
	askPrice100 := decimal.NewFromInt(100)

	bid := entity.Order{Type: entity.Bid, Price: bidPrice100}
	ask := entity.Order{Type: entity.Ask, Price: askPrice100}

	t.Run("Bid and Ask same price", func(t *testing.T) {
		assert.True(t, couldMatch(&bid, &ask))
	})

	t.Run("Bid price higher than Ask", func(t *testing.T) {
		bid.Price = bidPrice110
		assert.True(t, couldMatch(&bid, &ask))
	})

	t.Run("Bid price lower than Ask", func(t *testing.T) {
		bid.Price = decimal.NewFromInt(90)
		assert.False(t, couldMatch(&bid, &ask))
	})

	t.Run("Ask and Bid (reversed args)", func(t *testing.T) {
		bid.Price = bidPrice110
		assert.True(t, couldMatch(&ask, &bid))
	})

	t.Run("Same type (Bid, Bid)", func(t *testing.T) {
		bid2 := entity.Order{Type: entity.Bid, Price: bidPrice100}
		assert.False(t, couldMatch(&bid, &bid2))
	})

	t.Run("Nil order", func(t *testing.T) {
		assert.False(t, couldMatch(nil, &ask))
		assert.False(t, couldMatch(&bid, nil))
	})
}

func TestOrderHeap_MatchMake(t *testing.T) {
	// For matching Bids, we want a Min-Heap of Asks (lowest price first)
	cmpAsk := func(a, b *entity.Order) bool {
		return a.Price.LessThan(b.Price)
	}

	t.Run("Successful match", func(t *testing.T) {
		oh := newTestHeap(cmpAsk)
		ask := entity.Order{
			Type:       entity.Ask,
			Price:      decimal.NewFromInt(100),
			Quantity:   10,
			ValidUntil: time.Now().Add(time.Hour),
			OwnerDoc:   "owner1",
		}
		oh.Push(ask)

		bid := entity.Order{
			Type:     entity.Bid,
			Price:    decimal.NewFromInt(110),
			Quantity: 10,
			OwnerDoc: "owner2",
		}

		res := oh.Pop(bid)
		assert.Len(t, res.Matches, 1)
		assert.Equal(t, ask.ID, res.Matches[0].ID)
	})

	t.Run("No match due to price", func(t *testing.T) {
		oh := newTestHeap(cmpAsk)
		ask := entity.Order{
			Type:       entity.Ask,
			Price:      decimal.NewFromInt(120),
			Quantity:   10,
			ValidUntil: time.Now().Add(time.Hour),
			OwnerDoc:   "owner1",
		}
		oh.Push(ask)

		bid := entity.Order{
			Type:     entity.Bid,
			Price:    decimal.NewFromInt(110),
			Quantity: 10,
			OwnerDoc: "owner2",
		}

		res := oh.Pop(bid)
		assert.Len(t, res.Matches, 0)
	})

	t.Run("Expired order", func(t *testing.T) {
		oh := newTestHeap(cmpAsk)
		ask := entity.Order{
			Type:       entity.Ask,
			Price:      decimal.NewFromInt(100),
			Quantity:   10,
			ValidUntil: time.Now().Add(-time.Hour),
			OwnerDoc:   "owner1",
		}
		oh.Push(ask)

		bid := entity.Order{
			Type:     entity.Bid,
			Price:    decimal.NewFromInt(110),
			Quantity: 10,
			OwnerDoc: "owner2",
		}

		res := oh.Pop(bid)
		assert.Len(t, res.Matches, 0)
		assert.Len(t, res.Expired, 1)
	})

	t.Run("Same owner match prevention", func(t *testing.T) {
		oh := newTestHeap(cmpAsk)
		ask := entity.Order{
			Type:       entity.Ask,
			Price:      decimal.NewFromInt(100),
			Quantity:   10,
			ValidUntil: time.Now().Add(time.Hour),
			OwnerDoc:   "owner1",
		}
		oh.Push(ask)

		bid := entity.Order{
			Type:     entity.Bid,
			Price:    decimal.NewFromInt(110),
			Quantity: 10,
			OwnerDoc: "owner1",
		}

		res := oh.Pop(bid)
		assert.Len(t, res.Matches, 0)

		// Verify the order was pushed back
		peeked, ok := oh.Peek()
		assert.True(t, ok)
		assert.Equal(t, ask.ID, peeked.ID)
	})

	t.Run("Partial match quantity", func(t *testing.T) {
		oh := newTestHeap(cmpAsk)
		ask := entity.Order{
			Type:       entity.Ask,
			Price:      decimal.NewFromInt(100),
			Quantity:   5,
			ValidUntil: time.Now().Add(time.Hour),
			OwnerDoc:   "owner1",
		}
		oh.Push(ask)

		bid := entity.Order{
			Type:     entity.Bid,
			Price:    decimal.NewFromInt(110),
			Quantity: 10,
			OwnerDoc: "owner2",
		}

		res := oh.Pop(bid)
		assert.Len(t, res.Matches, 1)
		assert.Equal(t, ask, res.Matches[0])
	})

	t.Run("Multiple matches", func(t *testing.T) {
		oh := newTestHeap(cmpAsk)
		asks := []entity.Order{
			{Type: entity.Ask, Price: decimal.NewFromInt(100), Quantity: 5, ValidUntil: time.Now().Add(time.Hour), OwnerDoc: "owner1"},
			{Type: entity.Ask, Price: decimal.NewFromInt(90), Quantity: 5, ValidUntil: time.Now().Add(time.Hour), OwnerDoc: "owner2"},
		}
		for _, a := range asks {
			oh.Push(a)
		}

		bid := entity.Order{
			Type:     entity.Bid,
			Price:    decimal.NewFromInt(110),
			Quantity: 15,
			OwnerDoc: "owner3",
		}

		res := oh.Pop(bid)
		assert.Len(t, res.Matches, 2)
	})
}
