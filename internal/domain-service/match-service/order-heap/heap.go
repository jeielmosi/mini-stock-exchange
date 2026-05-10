package order_heaps

import (
	"mini-stock-exchange/internal/entity"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderHeap interface {
	Push(order entity.Order)
	Pop(order entity.Order) (*MatchDTO, error)
}

type MatchDTO struct {
	Matches []entity.Order
	Expired []uuid.UUID
}

type queryTrigger struct {
	Price     decimal.Decimal
	CreatedAt time.Time
}

func newQueryTrigger(order entity.Order) *queryTrigger {
	return &queryTrigger{
		Price:     order.Price,
		CreatedAt: order.CreatedAt,
	}
}

type orderHeap struct {
	heap []entity.Order
	cmp  func(a *entity.Order, b *entity.Order) bool
}

func newOrderHeap(capacity int, cmp func(a *entity.Order, b *entity.Order) bool) *orderHeap {
	return &orderHeap{
		heap: make([]entity.Order, 0, capacity),
		cmp:  cmp,
	}
}

func (oh *orderHeap) Len() int {
	return len(oh.heap)
}

func (oh *orderHeap) Cap() int {
	return cap(oh.heap)
}

func (oh *orderHeap) PopBack() (entity.Order, bool) {
	if len(oh.heap) == 0 {
		return entity.Order{}, false
	}
	last := oh.heap[len(oh.heap)-1]
	oh.heap = oh.heap[:len(oh.heap)-1]
	return last, true
}

func (oh *orderHeap) Push(order entity.Order) {
	child := len(oh.heap)
	oh.heap = append(oh.heap, order)

	for 0 < child {
		parent := (child - 1) / 2
		if oh.cmp(&oh.heap[child], &oh.heap[parent]) {
			oh.heap[child], oh.heap[parent] = oh.heap[parent], oh.heap[child]
			child = parent
		} else {
			break
		}
	}
}

func (oh *orderHeap) Peek() (entity.Order, bool) {
	if len(oh.heap) == 0 {
		return entity.Order{}, false
	}
	return oh.heap[0], true
}

func (oh *orderHeap) Drop() {
	if len(oh.heap) == 0 {
		return
	}
	last := oh.heap[len(oh.heap)-1]
	oh.heap = oh.heap[:len(oh.heap)-1]
	if len(oh.heap) == 0 {
		return
	}
	oh.heap[0] = last

	curr := 0
	for curr < len(oh.heap) {
		left := 2*curr + 1
		if len(oh.heap) <= left {
			break
		}
		right := 2*curr + 2
		if len(oh.heap) <= right {
			if oh.cmp(&oh.heap[left], &oh.heap[curr]) {
				oh.heap[left], oh.heap[curr] = oh.heap[curr], oh.heap[left]
			}
			break
		}

		up := right
		if oh.cmp(&oh.heap[left], &oh.heap[right]) {
			up = left
		}
		if oh.cmp(&oh.heap[up], &oh.heap[curr]) {
			oh.heap[up], oh.heap[curr] = oh.heap[curr], oh.heap[up]
			curr = up
		} else {
			break
		}
	}
}

func couldMatch(bid *entity.Order, ask *entity.Order) bool {
	if bid == nil || ask == nil {
		return false
	}

	if (bid.Type == entity.Ask) && (ask.Type == entity.Bid) {
		return couldMatch(ask, bid)
	}

	if (bid.Type != entity.Bid) || (ask.Type != entity.Ask) {
		return false
	}

	return bid.Price.GreaterThanOrEqual(ask.Price)
}

func (oh *orderHeap) Pop(order entity.Order) MatchDTO {
	now := time.Now()
	retry := []entity.Order{}
	matches := []entity.Order{}
	expired := []uuid.UUID{}

	quantity := order.Quantity
	for 0 < quantity {
		match, ok := oh.Peek()
		if !ok || !couldMatch(&order, &match) {
			break
		}
		oh.Drop()
		if order.OwnerDoc == match.OwnerDoc {
			retry = append(retry, match)
			continue
		}
		if match.ValidUntil.Before(now) {
			expired = append(expired, match.ID)
			continue
		}
		matches = append(matches, match)
		quantity -= match.Quantity
	}
	for _, o := range retry {
		oh.Push(o)
	}

	return MatchDTO{
		Matches: matches,
		Expired: expired,
	}
}
