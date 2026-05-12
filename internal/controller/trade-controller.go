package controller

import (
	"encoding/json"
	"net/http"

	"mini-stock-exchange/internal/dto"
	trade_service "mini-stock-exchange/internal/service/trade-service"

	"github.com/go-chi/chi/v5"
)

type tradeController struct {
	service trade_service.TradeService
}

func NewTradeController(service trade_service.TradeService) Controller {
	return &tradeController{
		service: service,
	}
}

func (h *tradeController) GetTrade(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req, err := dto.NewGetTradeRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	trade, err := h.service.GetTrade(req)
	if err != nil {
		http.Error(w, "trade not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(trade)
}

func (h *tradeController) RegisterRoutes(r chi.Router) {
	r.Get("/trades/{id}", h.GetTrade)
}
