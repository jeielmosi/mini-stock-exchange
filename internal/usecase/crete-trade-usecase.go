package usecase

import (
	"fmt"
	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/entity"
	//"mini-stock-exchange/internal/repository"
)

type CreateTradeUsecase interface {
	CreateTrade(dto dto.OrderMatch) (entity.Trade, error)
}

type createTradeUsecase struct {
	//trade repo
}

func NewCreateTradeUsecase() CreateTradeUsecase {
	return &createTradeUsecase{}
}

func (c *createTradeUsecase) CreateTrade(dto dto.OrderMatch) (entity.Trade, error) {
	if dto.Ask == nil || dto.Ask.Type != entity.Ask {
		return entity.Trade{}, fmt.Errorf("invalid ask order")
	}
	if dto.Bid == nil || dto.Bid.Type != entity.Bid {
		return entity.Trade{}, fmt.Errorf("invalid bid order")
	}
	if dto.Ask.Quantity < dto.Quantity || dto.Bid.Quantity < dto.Quantity {
		return entity.Trade{}, fmt.Errorf("quantity exceeds order quantity")
	}
	if dto.Price.IsZero() {
		return entity.Trade{}, fmt.Errorf("invalid price")
	}
	if dto.Ask.OwnerDoc == dto.Bid.OwnerDoc {
		return entity.Trade{}, fmt.Errorf("same owner")
	}

	trade := entity.Trade{
		BuyOrderID:  dto.Bid.ID,
		SellOrderID: dto.Ask.ID,
		Symbol:      dto.Ask.Symbol,
		Price:       dto.Price,
		Quantity:    dto.Quantity,
		ExecutedAt:  dto.ExecutedAt,
	}
	return trade, nil
}
