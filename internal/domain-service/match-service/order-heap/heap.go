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

type OrderQueue interface {
	Push(order entity.Order)
	Peek() (entity.Order, bool)
	Drop()
	PeekBack() (entity.Order, bool)
	DropBack()
	Len() int
	Cap() int
}

type MatchDTO struct {
	Matches []entity.Order
	Expired []uuid.UUID
	OK      bool
}

type QueryTrigger struct {
	Price     decimal.Decimal
	CreatedAt time.Time
}

func NewQueryTrigger(order entity.Order) *QueryTrigger {
	return &QueryTrigger{
		Price:     order.Price,
		CreatedAt: order.CreatedAt,
	}
}

type orderHeap struct {
	heap []entity.Order
	cmp  func(a *entity.Order, b *entity.Order) bool
}

func NewOrderHeap(capacity int, cmp func(a *entity.Order, b *entity.Order) bool) OrderQueue {
	return &orderHeap{
		heap: make([]entity.Order, 0, capacity),
		cmp:  cmp,
	}
}

func (oh *orderHeap) Len() int {
	return len(oh.heap)
}

func (oh *orderHeap) Cap() int {
	return cap(oh.heap) - len(oh.heap)
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

func (oh *orderHeap) PeekBack() (entity.Order, bool) {
	if len(oh.heap) == 0 {
		return entity.Order{}, false
	}
	return oh.heap[len(oh.heap)-1], true
}

func (oh *orderHeap) DropBack() {
	oh.heap = oh.heap[:len(oh.heap)-1]
}
