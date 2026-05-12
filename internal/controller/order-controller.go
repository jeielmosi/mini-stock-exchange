package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	match_service "mini-stock-exchange/internal/domain-service/match-service"
	"mini-stock-exchange/internal/dto"
	order_service "mini-stock-exchange/internal/service/order-service"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type orderController struct {
	match match_service.MatchService
	order order_service.OrderService
}

func NewOrderController(match match_service.MatchService, order order_service.OrderService) Controller {
	return &orderController{
		match: match,
		order: order,
	}
}

func (h *orderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	req, err := dto.NewCreateOrderRequest(r.Body)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		msg := "Validation error: ["
		for _, e := range err.(validator.ValidationErrors) {
			msg += fmt.Sprintf("%s,", e.Field())
		}
		msg = msg[:len(msg)-1] + "]"
		sendError(w, msg, http.StatusBadRequest)
		return
	}

	dto, err := h.match.SubmitOrder(r.Context(), req)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto)
}

func (h *orderController) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req, err := dto.NewGetOrderRequest(id)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.order.GetOrder(req)
	if err != nil {
		sendError(w, "order not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *orderController) RegisterRoutes(r chi.Router) {
	r.Post("/orders", h.CreateOrder)
	r.Get("/orders/{id}", h.GetOrder)
}
