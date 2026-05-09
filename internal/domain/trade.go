package domain

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
