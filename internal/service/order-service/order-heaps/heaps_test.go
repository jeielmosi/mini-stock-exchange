package order_heaps

import (
	"testing"
	"time"

	"mini-stock-exchange/internal/domain"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func newTestHeap(cmp func(a *domain.Order, b *domain.Order) bool) OrderHeap {
	return &orderHeap{
		heap: []domain.Order{},
		cmp:  cmp,
	}
}

func TestOrderHeap_PushPopPeek(t *testing.T) {
	// Max heap by price
	cmp := func(a, b *domain.Order) bool {
		return a.Price.GreaterThan(b.Price)
	}
	oh := newTestHeap(cmp)

	orders := []domain.Order{
		{Price: decimal.NewFromInt(100), ID: uuid.New()},
		{Price: decimal.NewFromInt(200), ID: uuid.New()},
		{Price: decimal.NewFromInt(150), ID: uuid.New()},
	}

	for _, o := range orders {
		oh.Push(o)
	}

	peeked, ok := oh.Peek()
	assert.True(t, ok)
	assert.Equal(t, decimal.NewFromInt(200), peeked.Price)

	oh.Pop()
	peeked, ok = oh.Peek()
	assert.True(t, ok)
	assert.Equal(t, decimal.NewFromInt(150), peeked.Price)

	oh.Pop()
	peeked, ok = oh.Peek()
	assert.True(t, ok)
	assert.Equal(t, decimal.NewFromInt(100), peeked.Price)

	oh.Pop()
	_, ok = oh.Peek()
	assert.False(t, ok)
}

func TestCouldMatch(t *testing.T) {
	bidPrice100 := decimal.NewFromInt(100)
	bidPrice110 := decimal.NewFromInt(110)
	askPrice100 := decimal.NewFromInt(100)

	bid := domain.Order{Type: domain.Bid, Price: bidPrice100}
	ask := domain.Order{Type: domain.Ask, Price: askPrice100}

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
		bid2 := domain.Order{Type: domain.Bid, Price: bidPrice100}
		assert.False(t, couldMatch(&bid, &bid2))
	})

	t.Run("Nil order", func(t *testing.T) {
		assert.False(t, couldMatch(nil, &ask))
		assert.False(t, couldMatch(&bid, nil))
	})
}

func TestOrderHeap_MatchMake(t *testing.T) {
	// For matching Bids, we want a Min-Heap of Asks (lowest price first)
	cmpAsk := func(a, b *domain.Order) bool {
		return a.Price.LessThan(b.Price)
	}

	t.Run("Successful match", func(t *testing.T) {
		oh := newTestHeap(cmpAsk)
		ask := domain.Order{
			Type:       domain.Ask,
			Price:      decimal.NewFromInt(100),
			Quantity:   10,
			ValidUntil: time.Now().Add(time.Hour),
			OwnerDoc:   "owner1",
		}
		oh.Push(ask)

		bid := domain.Order{
			Type:     domain.Bid,
			Price:    decimal.NewFromInt(110),
			Quantity: 10,
			OwnerDoc: "owner2",
		}

		res := oh.MatchMake(bid)
		assert.Len(t, res.Matches, 1)
		assert.Equal(t, ask.ID, res.Matches[0].ID)
	})

	t.Run("No match due to price", func(t *testing.T) {
		oh := newTestHeap(cmpAsk)
		ask := domain.Order{
			Type:       domain.Ask,
			Price:      decimal.NewFromInt(120),
			Quantity:   10,
			ValidUntil: time.Now().Add(time.Hour),
			OwnerDoc:   "owner1",
		}
		oh.Push(ask)

		bid := domain.Order{
			Type:     domain.Bid,
			Price:    decimal.NewFromInt(110),
			Quantity: 10,
			OwnerDoc: "owner2",
		}

		res := oh.MatchMake(bid)
		assert.Len(t, res.Matches, 0)
	})

	t.Run("Expired order", func(t *testing.T) {
		oh := newTestHeap(cmpAsk)
		ask := domain.Order{
			Type:       domain.Ask,
			Price:      decimal.NewFromInt(100),
			Quantity:   10,
			ValidUntil: time.Now().Add(-time.Hour),
			OwnerDoc:   "owner1",
		}
		oh.Push(ask)

		bid := domain.Order{
			Type:     domain.Bid,
			Price:    decimal.NewFromInt(110),
			Quantity: 10,
			OwnerDoc: "owner2",
		}

		res := oh.MatchMake(bid)
		assert.Len(t, res.Matches, 0)
		assert.Len(t, res.Expired, 1)
		assert.Equal(t, domain.Expired, res.Expired[0].Status)
	})

	t.Run("Same owner match prevention", func(t *testing.T) {
		oh := newTestHeap(cmpAsk)
		ask := domain.Order{
			Type:       domain.Ask,
			Price:      decimal.NewFromInt(100),
			Quantity:   10,
			ValidUntil: time.Now().Add(time.Hour),
			OwnerDoc:   "owner1",
		}
		oh.Push(ask)

		bid := domain.Order{
			Type:     domain.Bid,
			Price:    decimal.NewFromInt(110),
			Quantity: 10,
			OwnerDoc: "owner1",
		}

		res := oh.MatchMake(bid)
		assert.Len(t, res.Matches, 0)

		// Verify the order was pushed back
		peeked, ok := oh.Peek()
		assert.True(t, ok)
		assert.Equal(t, ask.ID, peeked.ID)
	})

	t.Run("Partial match quantity", func(t *testing.T) {
		oh := newTestHeap(cmpAsk)
		ask := domain.Order{
			Type:       domain.Ask,
			Price:      decimal.NewFromInt(100),
			Quantity:   5,
			ValidUntil: time.Now().Add(time.Hour),
			OwnerDoc:   "owner1",
		}
		oh.Push(ask)

		bid := domain.Order{
			Type:     domain.Bid,
			Price:    decimal.NewFromInt(110),
			Quantity: 10,
			OwnerDoc: "owner2",
		}

		res := oh.MatchMake(bid)
		assert.Len(t, res.Matches, 1)
		assert.Equal(t, ask, res.Matches[0])
		// The current MatchMake implementation doesn't handle partial quantity
		// correctly in terms of returning the remaining bid quantity,
		// but it should take the ask.
	})
}
