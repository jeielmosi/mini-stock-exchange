package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID                uuid.UUID
	BrokerID          string
	OwnerDoc          string
	Symbol            string
	Type              OrderType
	Price             decimal.Decimal
	Quantity          int
	RemainingQuantity int
	ValidUntil        time.Time
	Status            OrderStatus
	CreatedAt         time.Time
}

func NewOrder(
	brokerID, ownerDoc, symbol string,
	orderType OrderType, price decimal.Decimal,
	quantity int, validUntil time.Time,
) (Order, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return Order{}, err
	}

	if validUntil.Before(time.Now()) {
		return Order{}, fmt.Errorf("order already expirated")
	}
	return Order{
		ID:                uuid,
		BrokerID:          brokerID,
		OwnerDoc:          ownerDoc,
		Type:              orderType,
		Symbol:            symbol,
		Price:             price,
		Quantity:          quantity,
		RemainingQuantity: quantity,
		ValidUntil:        validUntil,
		Status:            Pending,
		CreatedAt:         time.Now(),
	}, nil
}
