package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"mini-stock-exchange/internal/domain"
	order_service "mini-stock-exchange/internal/service/order-service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderHandler struct {
	service order_service.OrderService
}

func NewOrderHandler(service order_service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

type CreateOrderRequest struct {
	BrokerID   string           `json:"broker_id"`
	OwnerDoc   string           `json:"owner_doc"`
	Type       domain.OrderType `json:"type"`
	Symbol     string           `json:"symbol"`
	Price      decimal.Decimal  `json:"price"`
	Quantity   int              `json:"quantity"`
	ValidUntil string           `json:"valid_until"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validUntil, err := time.Parse(time.RFC3339, req.ValidUntil)
	if err != nil {
		http.Error(w, "invalid valid_until format, use RFC3339", http.StatusBadRequest)
		return
	}

	order := domain.Order{
		BrokerID:   req.BrokerID,
		OwnerDoc:   req.OwnerDoc,
		Type:       req.Type,
		Symbol:     req.Symbol,
		Price:      req.Price,
		Quantity:   req.Quantity,
		ValidUntil: validUntil,
	}

	if err := h.service.SubmitOrder(r.Context(), order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": order.ID.String()})
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(id)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) RegisterRoutes(r chi.Router) {
	r.Post("/orders", h.CreateOrder)
	r.Get("/orders/{id}", h.GetOrder)
}
