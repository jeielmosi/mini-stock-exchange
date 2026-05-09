package order_heaps

import (
	"mini-stock-exchange/internal/domain"
	"mini-stock-exchange/internal/repository"
)

func less(lhs *domain.Order, rhs *domain.Order) bool {
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

func NewAskHeap(symbol string, orderRepo repository.OrderRepository) (OrderHeap, error) {
	h := orderHeap{
		heap: []domain.Order{},
		cmp:  less,
	}
	orders, err := orderRepo.GetAsks(symbol)
	if err != nil {
		return nil, err
	}
	for _, o := range orders {
		h.Push(o)
	}
	return &h, nil
}
