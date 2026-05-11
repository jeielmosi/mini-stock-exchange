package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type QueryFill struct {
	Price     decimal.Decimal
	CreatedAt time.Time
	Symbol    string
	Limit     int
}

func NewQueryFill(price decimal.Decimal, createdAt time.Time, symbol string, limit int) QueryFill {
	return QueryFill{
		Price:     price,
		CreatedAt: createdAt,
		Symbol:    symbol,
		Limit:     limit,
	}
}
