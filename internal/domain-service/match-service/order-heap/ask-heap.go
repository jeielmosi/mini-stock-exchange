package order_heaps

import (
	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"
	"time"
)

func lessQt(lhs *QueryTrigger, rhs *QueryTrigger) bool {
	if lhs == nil {
		return false
	}
	if rhs == nil {
		return true
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

	return lessQt(NewQueryTrigger(*lhs), NewQueryTrigger(*rhs))
}

type AskHeap struct {
	heap      OrderQueue
	orderRepo repository.OrderRepository
	qt        *QueryTrigger
	symbol    string
}

func NewAskHeap(symbol string, capacity int, orderRepo repository.OrderRepository) OrderHeap {
	return &AskHeap{
		heap:      NewOrderHeap(capacity, lessOrder),
		orderRepo: orderRepo,
		symbol:    symbol,
	}
}

func (a *AskHeap) Push(order entity.Order) {
	if a.qt != nil && lessQt(a.qt, NewQueryTrigger(order)) {
		return
	}
	if a.heap.Cap() == 0 {
		a.DropBack()
	}
	a.heap.Push(order)
}

// could return a valid dto and error
func (a *AskHeap) Pop(order entity.Order) (*MatchDTO, error) {
	now := time.Now()
	remaining := order.RemainingQuantity
	var retry []entity.Order
	var dto MatchDTO
	var err error
	for true {
		if remaining <= 0 {
			break
		}
		err = a.Fill()
		if err != nil {
			break
		}
		match, ok := a.heap.Peek()
		if !ok || order.Price.LessThan(match.Price) {
			break
		}
		a.heap.Drop()
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
		a.heap.Push(r)
	}

	return &dto, err
}

func (a *AskHeap) DropBack() {
	var qt *QueryTrigger = nil
	last, ok := a.heap.PeekBack()
	if ok {
		qt = NewQueryTrigger(last)
		a.heap.DropBack()
	}
	if lessQt(qt, a.qt) {
		a.qt = qt
	}
}

func (a *AskHeap) Fill() error {
	top, ok := a.heap.Peek()
	if !ok {
		orders, err := a.orderRepo.GetAsks(a.symbol, a.heap.Cap())
		if err != nil {
			return err
		}
		for _, o := range orders {
			a.heap.Push(o)
		}
		return nil
	}
	if a.qt == nil {
		return nil
	}
	if lessQt(a.qt, NewQueryTrigger(top)) {
		qf := dto.NewQueryFill(top.ID, top.Price, top.CreatedAt, a.symbol, a.heap.Cap())
		orders, err := a.orderRepo.GetAsksLT(qf)
		if err != nil {
			return err
		}
		for _, o := range orders {
			a.heap.Push(o)
		}
		a.qt = nil
	}
	return nil
}
