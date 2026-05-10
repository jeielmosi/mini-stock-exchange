package controller

import (
	"database/sql"
	match_service "mini-stock-exchange/internal/domain-service/match-service"
	"mini-stock-exchange/internal/repository"

	"github.com/go-chi/chi/v5"
)

type Controller interface {
	RegisterRoutes(r chi.Router)
}

func RegisterRoutes(r chi.Router, db *sql.DB) {
	orderRepo := repository.NewOrderRepository(db)
	orchestrator := match_service.NewOrchestrator(orderRepo)
	matchService := match_service.NewMatchService(orderRepo, orchestrator)

	orderController := NewOrderController(matchService)
	orderController.RegisterRoutes(r)

	healthController := NewHealthController()
	healthController.RegisterRoutes(r)

	metricController := NewMetricController()
	metricController.RegisterRoutes(r)
}
