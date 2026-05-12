package dto

import (
	dto_helper "mini-stock-exchange/internal/dto/helper"
	"mini-stock-exchange/internal/entity"

	"github.com/google/uuid"
)

type GetTradeRequest struct {
	ID uuid.UUID
}

func NewGetTradeRequest(id64 string) (GetTradeRequest, error) {
	id, err := dto_helper.DecodeUUIDv7(id64)
	if err != nil {
		return GetTradeRequest{}, err
	}
	err = dto_helper.IsValidUUIDv7(id)
	if err != nil {
		return GetTradeRequest{}, err
	}

	return GetTradeRequest{
		ID: id,
	}, nil
}

type GetTradeResponse struct {
	ID       string  `json:"id,omitempty"`
	Symbol   string  `json:"symbol,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Quantity int     `json:"quantity,omitempty"`
	AskID    string  `json:"ask_id,omitempty"`
	BidID    string  `json:"bid_id,omitempty"`
}

func NewGetTradeResponse(trade entity.Trade) (GetTradeResponse, error) {
	id, err := dto_helper.EncodeUUID(trade.ID)
	askID, err := dto_helper.EncodeUUID(trade.SellOrderID)
	buyID, err := dto_helper.EncodeUUID(trade.BuyOrderID)

	if err != nil {
		return GetTradeResponse{}, err
	}

	return GetTradeResponse{
		ID:       id,
		Symbol:   trade.Symbol,
		Price:    trade.Price.InexactFloat64(),
		Quantity: trade.Quantity,
		AskID:    askID,
		BidID:    buyID,
	}, nil
}
