package dto

import (
	dto_helper "mini-stock-exchange/internal/dto/helper"
	"mini-stock-exchange/internal/entity"

	"github.com/google/uuid"
)

type GetOrderRequest struct {
	ID uuid.UUID
}

func NewGetOrderRequest(id64 string) (GetOrderRequest, error) {
	id, err := dto_helper.DecodeUUIDv7(id64)
	if err != nil {
		return GetOrderRequest{}, err
	}
	err = dto_helper.IsValidUUIDv7(id)
	if err != nil {
		return GetOrderRequest{}, err
	}

	return GetOrderRequest{
		ID: id,
	}, nil
}

type GetOrderResponse struct {
	ID              string   `json:"id,omitempty"`
	Type            string   `json:"type,omitempty"`
	Symbol          string   `json:"symbol,omitempty"`
	Price           float64  `json:"price,omitempty"`
	Quantity        int      `json:"quantity,omitempty"`
	PendingQuantity int      `json:"pending_quantity,omitempty"`
	FilledQuantity  int      `json:"filled_quantity,omitempty"`
	Status          string   `json:"status,omitempty"`
	CreatedAt       string   `json:"created_at,omitempty"`
	ValidUntil      string   `json:"valid_until,omitempty"`
	Trades          []string `json:"trades,omitempty"`
}

func NewGetOrderResponse(order entity.Order, tradeIds []uuid.UUID) (GetOrderResponse, error) {
	id, err := dto_helper.EncodeUUID(order.ID)
	if err != nil {
		return GetOrderResponse{}, err
	}

	trades := make([]string, len(tradeIds))
	for t, trade := range tradeIds {
		trades[t], err = dto_helper.EncodeUUID(trade)
		if err != nil {
			return GetOrderResponse{}, err
		}
	}

	return GetOrderResponse{
		ID:              id,
		Type:            string(order.Type),
		Symbol:          order.Symbol,
		Price:           order.Price.InexactFloat64(),
		Quantity:        order.Quantity,
		PendingQuantity: order.RemainingQuantity,
		FilledQuantity:  order.Quantity - order.RemainingQuantity,
		Status:          string(order.Status),
		Trades:          trades,
	}, nil
}
