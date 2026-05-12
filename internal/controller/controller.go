package controller

import (
	"mini-stock-exchange/internal/repository"

	match_service "mini-stock-exchange/internal/domain-service/match-service"
	order_service "mini-stock-exchange/internal/service/order-service"
	trade_service "mini-stock-exchange/internal/service/trade-service"

	"github.com/go-chi/chi/v5"
)

type Controller interface {
	RegisterRoutes(r chi.Router)
}

type controller struct {
	r chi.Router

	orderController  Controller
	tradeController  Controller
	healthController Controller
	metricController Controller
}

func NewController(r chi.Router, orderRepo repository.OrderRepository, tradeRepo repository.TradeRepository) Controller {
	orchestrator := match_service.NewOrchestrator(orderRepo)
	orderService := order_service.NewOrderService(orderRepo, tradeRepo)
	tradeService := trade_service.NewTradeService(tradeRepo)
	matchService := match_service.NewMatchService(orchestrator, orderService)

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

func (c *controller) RegisterRoutes(r chi.Router) {
	c.orderController.RegisterRoutes(r)
	c.tradeController.RegisterRoutes(r)
	c.healthController.RegisterRoutes(r)
	c.metricController.RegisterRoutes(r)
}
