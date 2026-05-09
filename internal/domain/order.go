package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID                uuid.UUID
	BrokerID          string
	OwnerDoc          string
	Type              OrderType
	Symbol            string
	Price             decimal.Decimal
	Quantity          int
	RemainingQuantity int
	ValidUntil        time.Time
	Status            OrderStatus
	CreatedAt         time.Time
}
