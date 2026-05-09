package order_heaps

import (
	"mini-stock-exchange/internal/domain"
	"mini-stock-exchange/internal/repository"
)

func greater(lhs *domain.Order, rhs *domain.Order) bool {
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

func NewBidHeap(symbol string, orderRepo repository.OrderRepository) (OrderHeap, error) {
	h := orderHeap{
		heap: []domain.Order{},
		cmp:  greater,
	}
	orders, err := orderRepo.GetBids(symbol)
	if err != nil {
		return nil, err
	}
	for _, o := range orders {
		h.Push(o)
	}
	return &h, nil
}
