package usecase

import (
	"fmt"
	"mini-stock-exchange/internal/entity"
	"time"

	"github.com/google/uuid"
)

func NewTrade(order1 *entity.Order, order2 *entity.Order) (entity.Trade, error) {
	if (order1 == nil) || (order2 == nil) {
		return entity.Trade{}, fmt.Errorf("nil order")
	}

	if order1.Symbol != order2.Symbol {
		return entity.Trade{}, fmt.Errorf("different symbols for matching")
	}

	ask := order1
	bid := order2

	if (order1.Type == entity.Bid) && (order2.Type == entity.Ask) {
		ask = order2
		bid = order1
	}

	if ask.Type != entity.Ask || bid.Type != entity.Bid {
		return entity.Trade{}, fmt.Errorf("invalid order types for matching")
	}

	if bid.Price.LessThan(ask.Price) {
		return entity.Trade{}, fmt.Errorf("unmacthed price")
	}

	tradeQty := min(bid.RemainingQuantity, ask.RemainingQuantity)
	if tradeQty == 0 {
		return entity.Trade{}, fmt.Errorf("not enough quantity for matching")
	}

	bid.RemainingQuantity -= tradeQty
	if bid.RemainingQuantity == 0 {
		bid.Status = entity.Filled
	} else {
		bid.Status = entity.Partial
	}

	ask.RemainingQuantity -= tradeQty
	if ask.RemainingQuantity == 0 {
		ask.Status = entity.Filled
	} else {
		ask.Status = entity.Partial
	}

	tradeID, err := uuid.NewV7()
	if err != nil {
		return entity.Trade{}, fmt.Errorf("failed to generate UUID v7 for trade: %w", err)
	}

	trade := entity.Trade{
		ID:          tradeID,
		Symbol:      ask.Symbol,
		Price:       ask.Price,
		Quantity:    tradeQty,
		ExecutedAt:  time.Now(),
		BuyOrderID:  bid.ID,
		SellOrderID: ask.ID,
	}
	return trade, nil
}
