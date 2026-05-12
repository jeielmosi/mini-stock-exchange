package usecase

import (
	"fmt"
	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/entity"
)

type OrderMatchUsecase interface {
	MatchOrder(bid *entity.Order, ask *entity.Order) (dto.OrderMatch, error)
	UnmatchOrder(dto dto.OrderMatch)
}

type orderMatchUsecase struct{}

func NewOrderMatchUsecase() OrderMatchUsecase {
	return &orderMatchUsecase{}
}

func (u *orderMatchUsecase) MatchOrder(bid *entity.Order, ask *entity.Order) (dto.OrderMatch, error) {
	if (ask == nil) || (bid == nil) {
		return dto.OrderMatch{}, fmt.Errorf("nil pointers")
	}

	if (ask.Type == entity.Bid) && (bid.Type == entity.Ask) {
		return dto.OrderMatch{}, fmt.Errorf("can not match orders")
	}

	if ask.Type != entity.Ask || bid.Type != entity.Bid {
		return dto.OrderMatch{}, fmt.Errorf("could not match orders, needed ask and bid")
	}
	tradeQty := min(bid.RemainingQuantity, ask.RemainingQuantity)
	if tradeQty == 0 {
		return dto.OrderMatch{}, fmt.Errorf("no quantity to trade")
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

	return dto.NewOrderMatch(ask, bid, tradeQty, ask.Price), nil
}

func (u *orderMatchUsecase) UnmatchOrder(match dto.OrderMatch) {
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
