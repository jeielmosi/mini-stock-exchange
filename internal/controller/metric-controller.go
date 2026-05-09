package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricController struct{}

func NewMetricController() Controller {
	return &MetricController{}
}

func (h *MetricController) RegisterRoutes(r chi.Router) {
	r.Handle("/metrics", promhttp.Handler())
}
