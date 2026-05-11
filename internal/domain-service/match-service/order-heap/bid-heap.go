package order_heaps

import (
	"errors"
	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"
)

func greaterQt(lhs *queryTrigger, rhs *queryTrigger) bool {
	if lhs == nil {
		return true
	}
	if rhs == nil {
		return false
	}

	if lhs.Price.Equal(rhs.Price) {
		return lhs.CreatedAt.Before(rhs.CreatedAt)
	}
	return lhs.Price.GreaterThan(rhs.Price)
}

func greaterOrder(lhs *entity.Order, rhs *entity.Order) bool {
	if lhs == nil {
		return false
	}
	if rhs == nil {
		return true
	}

	return greaterQt(newQueryTrigger(*lhs), newQueryTrigger(*rhs))
}

type BidHeap struct {
	heap      *orderHeap
	orderRepo repository.OrderRepository
	qt        *queryTrigger
	symbol    string
}

func NewBidHeap(symbol string, capacity int, orderRepo repository.OrderRepository) OrderHeap {
	return &BidHeap{
		heap:      newOrderHeap(capacity, greaterOrder),
		orderRepo: orderRepo,
		symbol:    symbol,
	}
}

func (b *BidHeap) Push(order entity.Order) {
	if b.heap.Len() == b.heap.Cap() {
		var qt *queryTrigger = nil
		last, ok := b.heap.PopBack()
		if ok {
			qt = newQueryTrigger(last)
		}
		if greaterQt(qt, b.qt) {
			b.qt = qt
		}
	}
	b.heap.Push(order)
}

func (b *BidHeap) Pop(order entity.Order) (*MatchDTO, error) {
	err := b.fill()
	if err != nil {
		return nil, err
	}
	dto := b.heap.Pop(order)
	return &dto, nil
}

func (b *BidHeap) fill() error {
	if b.heap.Len() == 0 {
		orders, err := b.orderRepo.GetBids(b.symbol, b.heap.Cap())
		if err != nil {
			return err
		}
		for _, o := range orders {
			b.heap.Push(o)
		}
		return nil
	}
	top, ok := b.heap.Peek()
	if !ok {
		return errors.New("heap is not empty, but can not peek")
	}
	if greaterQt(b.qt, newQueryTrigger(top)) {
		qf := dto.NewQueryFill(b.qt.Price, b.qt.CreatedAt, b.symbol, b.heap.Cap()-b.heap.Len())
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
