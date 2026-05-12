package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	dto_helper "mini-stock-exchange/internal/dto/helper"
	"mini-stock-exchange/internal/entity"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateOrderRequest struct {
	BrokerID   string           `json:"broker_id" validate:"required"`
	OwnerDoc   string           `json:"owner_doc" validate:"required"`
	Type       entity.OrderType `json:"type" validate:"required"`
	Symbol     string           `json:"symbol" validate:"required"`
	Price      float64          `json:"price" validate:"required"`
	Quantity   int              `json:"quantity" validate:"required"`
	ValidUntil string           `json:"valid_until" validate:"required"`
}

func (r *CreateOrderRequest) ToOrder() (entity.Order, error) {
	validUntil, err := dto_helper.DateToEndOfDay(r.ValidUntil)
	if err != nil {
		return entity.Order{}, fmt.Errorf("invalid valid_until format, use Date Only format")
	}

	return entity.NewOrder(
		r.BrokerID, r.OwnerDoc, r.Symbol,
		r.Type, decimal.NewFromFloat(r.Price),
		r.Quantity, validUntil,
	)
}

func (r *CreateOrderRequest) validate() error {
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
	return err
}

func NewCreateOrderRequest(body io.ReadCloser) (CreateOrderRequest, error) {
	var dto CreateOrderRequest
	if err := json.NewDecoder(body).Decode(&dto); err != nil {
		return dto, err
	}
	if err := dto.validate(); err != nil {
		return dto, err
	}
	return dto, nil
}

type CreateOrderResponse struct {
	ID string `json:"id,omitempty"`
}

func NewCreateOrderResponse(id uuid.UUID) (CreateOrderResponse, error) {
	id64, err := dto_helper.EncodeUUID(id)
	if err != nil {
		return CreateOrderResponse{}, err
	}
	return CreateOrderResponse{ID: id64}, nil
}
