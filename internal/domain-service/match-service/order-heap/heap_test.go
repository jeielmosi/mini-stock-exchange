package order_heaps

import (
	"testing"

	"mini-stock-exchange/internal/entity"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func newTestHeap(cmp func(a *entity.Order, b *entity.Order) bool) OrderQueue {
	return &orderHeap{
		heap: []entity.Order{},
		cmp:  cmp,
	}
}

func TestOrderHeap_AddPeekDrop(t *testing.T) {
	cmp := func(a, b *entity.Order) bool {
		return a.Price.GreaterThan(b.Price)
	}
	oh := newTestHeap(cmp)

	// PopBack on empty heap
	_, ok := oh.PeekBack()
	assert.False(t, ok)

	order := entity.Order{Price: decimal.NewFromInt(100), ID: uuid.New()}
	oh.Push(order)

	popped, ok := oh.PeekBack()
	assert.True(t, ok)
	assert.Equal(t, order.ID, popped.ID)
	assert.Equal(t, 1, oh.Len())

	oh.DropBack()
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
