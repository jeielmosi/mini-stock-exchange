package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	dto_helper "mini-stock-exchange/internal/dto/helper"
	"mini-stock-exchange/internal/entity"

	"github.com/go-playground/validator/v10"
)

type CreateBrokerRequest struct {
	Name string `json:"name" validate:"required"`
}

func (r *CreateBrokerRequest) validate() error {
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

func (r *CreateBrokerRequest) ToBroker() (entity.Broker, error) {
	return entity.NewBroker(r.Name)
}

func NewCreateBrokerRequest(body io.ReadCloser) (CreateBrokerRequest, error) {
	var dto CreateBrokerRequest
	if err := json.NewDecoder(body).Decode(&dto); err != nil {
		return dto, err
	}
	if err := dto.validate(); err != nil {
		return dto, err
	}
	return dto, nil
}

type CreateBrokerResponse struct {
	ID string `json:"id,omitempty"`
}

func NewCreateBrokerResponse(broker entity.Broker) (CreateBrokerResponse, error) {
	id, err := dto_helper.EncodeUUID(broker.ID)
	if err != nil {
		return CreateBrokerResponse{}, err
	}
	return CreateBrokerResponse{
		ID: id,
	}, nil
}
