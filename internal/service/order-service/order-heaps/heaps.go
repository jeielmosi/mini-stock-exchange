package order_heaps

import (
	"mini-stock-exchange/internal/domain"
	"time"
)

type MatchDTO struct {
	Matches []domain.Order
	Expired []domain.Order
}

type OrderHeap interface {
	Pop()
	Peek() (domain.Order, bool)
	Push(order domain.Order)
	MatchMake(order domain.Order) MatchDTO
}

type orderHeap struct {
	heap []domain.Order
	cmp  func(a *domain.Order, b *domain.Order) bool
}

func (oh *orderHeap) Push(order domain.Order) {
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

func (oh *orderHeap) Peek() (domain.Order, bool) {
	if len(oh.heap) == 0 {
		return domain.Order{}, false
	}
	return oh.heap[0], true
}

func (oh *orderHeap) Pop() {
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

func couldMatch(bid *domain.Order, ask *domain.Order) bool {
	if bid == nil || ask == nil {
		return false
	}

	if (bid.Type == domain.Ask) && (ask.Type == domain.Bid) {
		return couldMatch(ask, bid)
	}

	if (bid.Type != domain.Bid) || (ask.Type != domain.Ask) {
		return false
	}

	return bid.Price.GreaterThanOrEqual(ask.Price)
}

func (oh *orderHeap) MatchMake(order domain.Order) MatchDTO {
	now := time.Now()
	retry := []domain.Order{}
	matches := []domain.Order{}
	expired := []domain.Order{}

	quantity := order.Quantity
	for 0 < quantity {
		match, ok := oh.Peek()
		if !ok || !couldMatch(&order, &match) {
			break
		}
		oh.Pop()
		if order.OwnerDoc == match.OwnerDoc {
			retry = append(retry, match)
			continue
		}
		if match.ValidUntil.Before(now) {
			match.Status = domain.Expired
			expired = append(expired, match)
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
