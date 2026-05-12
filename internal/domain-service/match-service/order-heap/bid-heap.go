package order_heaps

import (
	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"
	"time"
)

func greaterQt(lhs *QueryTrigger, rhs *QueryTrigger) bool {
	if lhs == nil {
		return false
	}
	if rhs == nil {
		return true
	}

	if lhs.Price.Equal(rhs.Price) {
		return lhs.CreatedAt.Before(rhs.CreatedAt)
	}
	return lhs.Price.GreaterThan(rhs.Price)
}

func greaterOrder(lhs *entity.Order, rhs *entity.Order) bool {
	if lhs == nil {
		return true
	}
	if rhs == nil {
		return false
	}

	return greaterQt(NewQueryTrigger(*lhs), NewQueryTrigger(*rhs))
}

type BidHeap struct {
	heap      OrderQueue
	orderRepo repository.OrderRepository
	qt        *QueryTrigger
	symbol    string
}

func NewBidHeap(symbol string, capacity int, orderRepo repository.OrderRepository) OrderHeap {
	return &BidHeap{
		heap:      NewOrderHeap(capacity, greaterOrder),
		orderRepo: orderRepo,
		symbol:    symbol,
	}
}

func (b *BidHeap) Push(order entity.Order) {
	if b.qt != nil && greaterQt(b.qt, NewQueryTrigger(order)) {
		return
	}
	if b.heap.Cap() == 0 {
		b.DropBack()
	}
	b.heap.Push(order)
}

// could return a valid dto and error
func (b *BidHeap) Pop(order entity.Order) (*MatchDTO, error) {
	now := time.Now()
	remaining := order.RemainingQuantity
	var retry []entity.Order
	var dto MatchDTO
	err := b.Init()
	if err != nil {
		return nil, err
	}

	for true {
		if remaining <= 0 {
			break
		}
		err = b.Fill()
		if err != nil {
			break
		}
		match, ok := b.heap.Peek()
		if !ok || order.Price.GreaterThan(match.Price) {
			break
		}
		b.heap.Drop()
		if match.OwnerDoc == order.OwnerDoc {
			retry = append(retry, match)
			continue
		}
		if match.ValidUntil.Before(now) {
			dto.Expired = append(dto.Expired, match.ID)
			continue
		}
		remaining -= match.RemainingQuantity
		dto.Matches = append(dto.Matches, match)
	}

	for _, r := range retry {
		b.heap.Push(r)
	}

	return &dto, err
}

func (b *BidHeap) DropBack() {
	var qt *QueryTrigger = nil
	last, ok := b.heap.PeekBack()
	if ok {
		qt = NewQueryTrigger(last)
		b.heap.DropBack()
	}
	if greaterQt(qt, b.qt) {
		b.qt = qt
	}
}

func (b *BidHeap) Init() error {
	if b.heap.Len() != 0 {
		return nil
	}
	orders, err := b.orderRepo.GetBids(b.symbol, b.heap.Cap())
	if err != nil {
		return err
	}
	for _, o := range orders {
		b.heap.Push(o)
	}
	return nil
}

func (b *BidHeap) Fill() error {
	top, ok := b.heap.Peek()
	if !ok || b.qt == nil {
		return nil
	}
	if greaterQt(b.qt, NewQueryTrigger(top)) {
		qf := dto.NewQueryFill(top.ID, top.Price, top.CreatedAt, b.symbol, b.heap.Cap())
		orders, err := b.orderRepo.GetBidsGT(qf)
		if err != nil {
			return err
		}
		for _, o := range orders {
			b.heap.Push(o)
		}
		b.qt = nil
	}
	return nil
}
