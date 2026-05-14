package usecase

import (
	"fmt"
	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/entity"
	"time"
)

type OrderMatchUsecase interface {
	CouldMatch(order1 *entity.Order, order2 *entity.Order) bool
	MatchOrder(order1 *entity.Order, order2 *entity.Order) (dto.OrderMatch, error)
	UnmatchOrder(dto dto.OrderMatch)
}

type orderMatchUsecase struct{}

func NewOrderMatchUsecase() OrderMatchUsecase {
	return &orderMatchUsecase{}
}

func (u *orderMatchUsecase) CouldMatch(o1 *entity.Order, o2 *entity.Order) bool {
	if (o1 == nil) || (o2 == nil) {
		return false
	}
	if o1.OwnerDoc == o2.OwnerDoc {
		return false
	}

	now := time.Now()
	if o1.ValidUntil.Before(now) || o2.ValidUntil.Before(now) {
		return false
	}

	ask, bid := o1, o2
	if ask.Type == entity.Bid && bid.Type == entity.Ask {
		ask, bid = bid, ask
	}
	if ask.Type != entity.Ask || bid.Type != entity.Bid {
		return false
	}

	if ask.Type != entity.Ask || bid.Type != entity.Bid {
		return false
	}

	if bid.Price.Cmp(ask.Price) < 0 {
		return false
	}
	return true
}

func (u *orderMatchUsecase) MatchOrder(o1 *entity.Order, o2 *entity.Order) (dto.OrderMatch, error) {
	if !u.CouldMatch(o1, o2) {
		return dto.OrderMatch{}, fmt.Errorf("could not match orders")
	}

	ask, bid := o1, o2
	if ask.Type == entity.Bid && bid.Type == entity.Ask {
		ask, bid = bid, ask
	}

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
