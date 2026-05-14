package order_heaps

import (
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/utils"
)

func bidCmp(lhs *entity.Order, rhs *entity.Order) bool {
	if lhs == nil {
		return false
	}
	if rhs == nil {
		return true
	}

	if lhs.Price.Cmp(rhs.Price) == 0 {
		return lhs.CreatedAt.Before(rhs.CreatedAt)
	}
	return 0 < lhs.Price.Cmp(rhs.Price)
}

func NewBidHeap(capacity int) *utils.PriorityQueue[*entity.Order] {
	return utils.NewPriorityQueue(bidCmp, capacity)
}
