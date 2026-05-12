package usecase

import (
	"fmt"
	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/entity"

	"github.com/google/uuid"
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

	id, err := uuid.NewV7()
	if err != nil {
		return entity.Trade{}, err
	}

	trade := entity.NewTrade(id, dto.Bid.ID, dto.Ask.ID, dto.Ask.Symbol, dto.Price, dto.Quantity, dto.ExecutedAt)
	return trade, nil
}
