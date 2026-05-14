package entity

import (
	"time"

	"github.com/google/uuid"
	"math/big"
)

type Trade struct {
	ID          uuid.UUID
	BuyOrderID  uuid.UUID
	SellOrderID uuid.UUID
	Symbol      string
	Price       *big.Rat
	Quantity    int
	ExecutedAt  time.Time
}

func NewTrade(id uuid.UUID, buyOrderID uuid.UUID, sellOrderID uuid.UUID, symbol string, price *big.Rat, quantity int, executedAt time.Time) Trade {
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
