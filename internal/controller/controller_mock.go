package controller

import (
	"mini-stock-exchange/internal/repository"

	match_service "mini-stock-exchange/internal/domain-service/match-service"
	order_service "mini-stock-exchange/internal/service/order-service"
	trade_service "mini-stock-exchange/internal/service/trade-service"

	"github.com/go-chi/chi/v5"
)

func NewMockController(r chi.Router, orderRepo repository.OrderRepository, tradeRepo repository.TradeRepository) Controller {
	orchestrator := match_service.NewMockOrchestrator(orderRepo)
	orderService := order_service.NewOrderService(orderRepo)
	tradeService := trade_service.NewTradeService(tradeRepo)
	matchService := match_service.NewMatchService(orchestrator)

	orderController := NewOrderController(matchService, orderService)
	tradeController := NewTradeController(tradeService)
	healthController := NewHealthController()
	metricController := NewMetricController()

	ctrl := controller{
		r:                r,
		orderController:  orderController,
		tradeController:  tradeController,
		healthController: healthController,
		metricController: metricController,
	}
	return &ctrl
}
