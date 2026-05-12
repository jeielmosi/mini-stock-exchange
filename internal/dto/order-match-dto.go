package dto

import (
	"mini-stock-exchange/internal/entity"
	"time"

	"github.com/shopspring/decimal"
)

type OrderMatch struct {
	Ask        *entity.Order
	Bid        *entity.Order
	Quantity   int
	Price      decimal.Decimal
	ExecutedAt time.Time
}

func NewOrderMatch(ask *entity.Order, bid *entity.Order, quantity int, price decimal.Decimal) OrderMatch {
	now := time.Now().UTC()
	return OrderMatch{
		Ask:        ask,
		Bid:        bid,
		Quantity:   quantity,
		Price:      price,
		ExecutedAt: now,
	}
}
