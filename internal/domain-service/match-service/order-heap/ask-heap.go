package order_heaps

import (
	"errors"
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"
)

func lessQt(lhs *queryTrigger, rhs *queryTrigger) bool {
	if lhs == nil {
		return true
	}
	if rhs == nil {
		return false
	}

	if lhs.Price.Equal(rhs.Price) {
		return lhs.CreatedAt.Before(rhs.CreatedAt)
	}

	return lhs.Price.LessThan(rhs.Price)
}

func lessOrder(lhs *entity.Order, rhs *entity.Order) bool {
	if lhs == nil {
		return true
	}
	if rhs == nil {
		return false
	}

	return lessQt(newQueryTrigger(*lhs), newQueryTrigger(*rhs))
}

type AskHeap struct {
	heap      *orderHeap
	orderRepo repository.OrderRepository
	qt        *queryTrigger
	symbol    string
}

func NewAskHeap(symbol string, capacity int, orderRepo repository.OrderRepository) OrderHeap {
	return &AskHeap{
		heap:      newOrderHeap(capacity, lessOrder),
		orderRepo: orderRepo,
		symbol:    symbol,
	}
}

func (a *AskHeap) Push(order entity.Order) {
	if a.heap.Len() == a.heap.Cap() {
		var qt *queryTrigger = nil
		last, ok := a.heap.PopBack()
		if ok {
			qt = newQueryTrigger(last)
		}
		if lessQt(qt, a.qt) {
			a.qt = qt
		}
	}
	a.heap.Push(order)
}

func (a *AskHeap) Pop(order entity.Order) (*MatchDTO, error) {
	err := a.fill()
	if err != nil {
		return nil, err
	}
	dto := a.heap.Pop(order)

	return &dto, nil
}

func (a *AskHeap) fill() error {
	if a.heap.Len() == 0 {
		orders, err := a.orderRepo.GetAsks(a.symbol, a.heap.Cap())
		if err != nil {
			return err
		}
		for _, o := range orders {
			a.heap.Push(o)
		}
		return nil
	}
	top, ok := a.heap.Peek()
	if !ok {
		return errors.New("heap is not empty, but can not peek")
	}
	if lessQt(newQueryTrigger(top), a.qt) {
		return nil
	}

	orders, err := a.orderRepo.GetAsksGT(top, a.heap.Cap()-a.heap.Len())
	if err != nil {
		return err
	}
	for _, o := range orders {
		a.heap.Push(o)
	}
	return nil
}
