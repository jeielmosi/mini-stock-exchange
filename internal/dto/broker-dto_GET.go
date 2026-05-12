package dto

import (
	dto_helper "mini-stock-exchange/internal/dto/helper"
	"mini-stock-exchange/internal/entity"

	"github.com/google/uuid"
)

type GetBrokerRequest struct {
	ID uuid.UUID
}

func NewGetBrokerRequest(id64 string) (GetBrokerRequest, error) {
	id, err := dto_helper.DecodeUUIDv7(id64)
	if err != nil {
		return GetBrokerRequest{}, err
	}

	return GetBrokerRequest{
		ID: id,
	}, nil
}

type GetBrokerResponse struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func NewGetBrokerResponse(broker entity.Broker) (GetBrokerResponse, error) {
	id, err := dto_helper.EncodeUUID(broker.ID)
	if err != nil {
		return GetBrokerResponse{}, err
	}

	return GetBrokerResponse{
		ID:   id,
		Name: broker.Name,
	}, nil
}
