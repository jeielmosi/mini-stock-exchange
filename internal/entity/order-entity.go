package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"math/big"
)

type OrderType string

const (
	Bid OrderType = "BID"
	Ask OrderType = "ASK"
)

type OrderStatus string

const (
	Pending OrderStatus = "PENDING"
	Partial OrderStatus = "PARTIAL"
	Filled  OrderStatus = "FILLED"
	Expired OrderStatus = "EXPIRED"
)

type Order struct {
	ID                uuid.UUID
	BrokerID          uuid.UUID
	BrokerName        string
	OwnerDoc          string
	Symbol            string
	Type              OrderType
	Price             *big.Rat
	Quantity          int
	RemainingQuantity int
	ValidUntil        time.Time
	Status            OrderStatus
	CreatedAt         time.Time
}

func NewOrder(
	brokerID uuid.UUID,
	ownerDoc, symbol string,
	orderType OrderType, price *big.Rat,
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
