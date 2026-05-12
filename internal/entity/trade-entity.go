package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Trade struct {
	ID          uuid.UUID
	BuyOrderID  uuid.UUID
	SellOrderID uuid.UUID
	Symbol      string
	Price       decimal.Decimal
	Quantity    int
	ExecutedAt  time.Time
}

func NewTrade(id uuid.UUID, buyOrderID uuid.UUID, sellOrderID uuid.UUID, symbol string, price decimal.Decimal, quantity int, executedAt time.Time) Trade {
	return Trade{
		ID:          id,
		BuyOrderID:  buyOrderID,
		SellOrderID: sellOrderID,
		Symbol:      symbol,
		Price:       price,
		Quantity:    quantity,
		ExecutedAt:  executedAt,
	}
}
