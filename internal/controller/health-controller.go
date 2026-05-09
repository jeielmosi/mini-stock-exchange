package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HealthController struct{}

func NewHealthController() Controller {
	return &HealthController{}
}

func (h *HealthController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *HealthController) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.HealthCheck)
}
