package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type QueryFill struct {
	Price     decimal.Decimal
	CreatedAt time.Time
	Symbol    string
	Limit     int
	ID        uuid.UUID
}

func NewQueryFill(id uuid.UUID, price decimal.Decimal, createdAt time.Time, symbol string, limit int) QueryFill {
	return QueryFill{
		Price:     price,
		CreatedAt: createdAt,
		Symbol:    symbol,
		Limit:     limit,
		ID:        id,
	}
}
