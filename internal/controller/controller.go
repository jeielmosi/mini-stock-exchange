package controller

import (
	"database/sql"
	"mini-stock-exchange/internal/repository"
	order_service "mini-stock-exchange/internal/service/order-service"

	"github.com/go-chi/chi/v5"
)

type Controller interface {
	RegisterRoutes(r chi.Router)
}

func RegisterRoutes(r chi.Router, db *sql.DB) {
	orderRepo := repository.NewOrderRepository(db)
	tradeRepo := repository.NewTradeRepository(db)
	orchestrator := order_service.NewOrchestrator(orderRepo, tradeRepo)
	orderService := order_service.NewOrderService(orderRepo, tradeRepo, orchestrator)

	orderController := NewOrderController(orderService)
	orderController.RegisterRoutes(r)

	healthController := NewHealthController()
	healthController.RegisterRoutes(r)

	metricController := NewMetricController()
	metricController.RegisterRoutes(r)
}
