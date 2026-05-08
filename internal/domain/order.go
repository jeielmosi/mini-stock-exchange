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

type OrderRepository interface {
	Create(order *Order) error
	GetByID(id uuid.UUID) (*Order, error)
	Update(order *Order) error
	FindMatches(order Order) ([]Order, error)
	UpdateRemainingQuantity(id uuid.UUID, quantity int, status OrderStatus) error
}
