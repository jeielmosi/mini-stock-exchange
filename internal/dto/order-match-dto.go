package dto

import (
	"math/big"
	"mini-stock-exchange/internal/entity"
	"time"
)

type OrderMatch struct {
	Ask        *entity.Order
	Bid        *entity.Order
	Quantity   int
	Price      *big.Rat
	ExecutedAt time.Time
}

func NewOrderMatch(ask *entity.Order, bid *entity.Order, quantity int, price *big.Rat) OrderMatch {
	now := time.Now()
	return OrderMatch{
		Ask:        ask,
		Bid:        bid,
		Quantity:   quantity,
		Price:      price,
		ExecutedAt: now,
	}
}
