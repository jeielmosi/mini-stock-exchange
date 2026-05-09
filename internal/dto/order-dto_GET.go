package dto

import (
	"errors"
	"fmt"
	"time"

	"mini-stock-exchange/internal/domain"
	dto_helper "mini-stock-exchange/internal/dto/helper"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type GetOrderRequest struct {
	ID string `json:"id" validate:"required"`
}

func (r *GetOrderRequest) validate() error {
	var validate = validator.New()
	err := validate.Struct(r)
	if err != nil {
		msg := "Validation error: ["
		for _, e := range err.(validator.ValidationErrors) {
			msg += fmt.Sprintf("%s,", e.Field())
		}
		msg = msg[:len(msg)-1] + "]"
		return errors.New(msg)
	}

	err = dto_helper.IsValidUUIDv7(r.ID)
	return err
}

func (r *GetOrderRequest) ToUUID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

func NewGetOrderRequest(id string) (GetOrderRequest, error) {
	dto := GetOrderRequest{ID: id}
	if err := dto.validate(); err != nil {
		return dto, err
	}
	return dto, nil
}

type GetOrderResponse struct {
	ID                string  `json:"id" validate:"required"`
	Type              string  `json:"type" validate:"required"`
	Symbol            string  `json:"symbol" validate:"required"`
	Price             float64 `json:"price" validate:"required"`
	Quantity          int     `json:"quantity" validate:"required"`
	RemainingQuantity int     `json:"remaining_quantity" validate:"required"`
	Status            string  `json:"status" validate:"required"`
	CreatedAt         string  `json:"created_at" validate:"required"`
	ValidUntil        string  `json:"valid_until" validate:"required"`
}

func NewGetOrderResponse(order domain.Order) GetOrderResponse {
	return GetOrderResponse{
		ID:                order.ID.String(),
		Type:              string(order.Type),
		Symbol:            order.Symbol,
		Price:             order.Price.InexactFloat64(),
		Quantity:          order.Quantity,
		RemainingQuantity: order.RemainingQuantity,
		Status:            string(order.Status),
		CreatedAt:         order.CreatedAt.Format(time.DateOnly),
		ValidUntil:        order.ValidUntil.Format(time.DateOnly),
	}
}
