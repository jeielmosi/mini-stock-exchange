package usecase

import (
	"mini-stock-exchange/internal/entity"

	"github.com/google/uuid"
)

type CreateTradeUsecase interface {
	CreateTrade(dto OrderMatch) (entity.Trade, error)
}

type createTradeUsecase struct{}

func NewCreateTradeUsecase() CreateTradeUsecase {
	return &createTradeUsecase{}
}

func (c *createTradeUsecase) CreateTrade(dto OrderMatch) (entity.Trade, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return entity.Trade{}, err
	}

	trade := entity.NewTrade(id, dto.Bid.ID, dto.Ask.ID, dto.Ask.Symbol, dto.Price, dto.Quantity, dto.ExecutedAt)
	return trade, nil
}
