package controller

import (
	"encoding/json"
	"net/http"

	"mini-stock-exchange/internal/dto"
	broker_service "mini-stock-exchange/internal/service/broker-service"

	"github.com/go-chi/chi/v5"
)

type brokerController struct {
	service broker_service.BrokerService
}

func NewBrokerController(service broker_service.BrokerService) Controller {
	return &brokerController{
		service: service,
	}
}

func (h *brokerController) GetBroker(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req, err := dto.NewGetBrokerRequest(id)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	broker, err := h.service.GetBroker(req)
	if err != nil {
		sendError(w, "broker not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(broker)
}

func (h *brokerController) RegisterRoutes(r chi.Router) {
	r.Get("/brokers/{id}", h.GetBroker)
	r.Post("/brokers", h.CreateBroker)
}

func (h *brokerController) CreateBroker(w http.ResponseWriter, r *http.Request) {
	req, err := dto.NewCreateBrokerRequest(r.Body)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	broker, err := h.service.Create(req)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(broker)
}