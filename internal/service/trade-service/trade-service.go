package trade_service

import (
	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/repository"
)

type TradeService interface {
	GetTrade(req dto.GetTradeRequest) (dto.GetTradeResponse, error)
}

type tradeService struct {
	tradeRepo repository.TradeRepository
}

func NewTradeService(tradeRepo repository.TradeRepository) TradeService {
	return &tradeService{
		tradeRepo: tradeRepo,
	}
}

func (t *tradeService) GetTrade(req dto.GetTradeRequest) (dto.GetTradeResponse, error) {
	id := req.ID
	trade, err := t.tradeRepo.GetByID(id)
	if err != nil {
		return dto.GetTradeResponse{}, err
	}
	res, err := dto.NewGetTradeResponse(trade)
	if err != nil {
		return dto.GetTradeResponse{}, err
	}
	return res, nil
}
