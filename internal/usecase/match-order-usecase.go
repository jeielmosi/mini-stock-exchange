package usecase

import (
	"fmt"
	"math/big"
	"mini-stock-exchange/internal/entity"
	"time"
)

type OrderMatchUsecase interface {
	ValidateMatch(order *entity.Order, match *entity.Order) (*ValidateMatch, error)
	MatchOrder(order *entity.Order, match *entity.Order) (*OrderMatch, error)
	UnmatchOrder(dto OrderMatch)
}

type orderMatchUsecase struct{}

func NewOrderMatchUsecase() OrderMatchUsecase {
	return &orderMatchUsecase{}
}

type ValidateMatch struct {
	Ask     *entity.Order
	Bid     *entity.Order
	Retry   *entity.Order
	Expired *entity.Order
}

func (u *orderMatchUsecase) ValidateMatch(order *entity.Order, match *entity.Order) (*ValidateMatch, error) {
	now := time.Now()
	if (order == nil) || (match == nil) {
		return nil, fmt.Errorf("order or match is nil")
	}
	if order.ValidUntil.Before(now) {
		return nil, fmt.Errorf("order is expired")
	}
	var validateMatch ValidateMatch
	if match.ValidUntil.Before(now) {
		validateMatch.Expired = match
		return &validateMatch, nil
	}
	if order.OwnerDoc == match.OwnerDoc {
		validateMatch.Retry = order
		return &validateMatch, nil
	}

	ask, bid := order, match
	if ask.Type == entity.Bid && bid.Type == entity.Ask {
		ask, bid = bid, ask
	}
	if ask.Type != entity.Ask || bid.Type != entity.Bid {
		return nil, fmt.Errorf("could not match, should be an ask and a bid")
	}
	if bid.Price.Cmp(ask.Price) < 0 {
		return &validateMatch, nil
	}

	tradeQty := min(ask.Quantity, bid.Quantity)
	if tradeQty <= 0 {
		return nil, fmt.Errorf("trade quantity is invalid")
	}

	validateMatch.Ask = ask
	validateMatch.Bid = bid
	return &validateMatch, nil
}

type OrderMatch struct {
	Ask        *entity.Order
	Bid        *entity.Order
	Quantity   int
	Price      *big.Rat
	ExecutedAt time.Time
}

func NewOrderMatch(ask *entity.Order, bid *entity.Order, quantity int, price *big.Rat) *OrderMatch {
	now := time.Now()
	return &OrderMatch{
		Ask:        ask,
		Bid:        bid,
		Quantity:   quantity,
		Price:      price,
		ExecutedAt: now,
	}
}

func (u *orderMatchUsecase) MatchOrder(order *entity.Order, match *entity.Order) (*OrderMatch, error) {

	dto, err := u.ValidateMatch(order, match)
	if err != nil {
		return nil, err
	}

	if dto.Ask == nil || dto.Bid == nil {
		return nil, fmt.Errorf("should not match by business logic")
	}

	ask, bid := dto.Ask, dto.Bid

	tradeQty := min(ask.Quantity, bid.Quantity)

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

	return NewOrderMatch(ask, bid, tradeQty, ask.Price), nil
}

func (u *orderMatchUsecase) UnmatchOrder(match OrderMatch) {
	match.Bid.RemainingQuantity += match.Quantity
	if match.Bid.RemainingQuantity == match.Bid.Quantity {
		match.Bid.Status = entity.Pending
	} else {
		match.Bid.Status = entity.Partial
	}

	match.Ask.RemainingQuantity += match.Quantity
	if match.Ask.RemainingQuantity == match.Ask.Quantity {
		match.Ask.Status = entity.Pending
	} else {
		match.Ask.Status = entity.Partial
	}
}
