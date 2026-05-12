package controller

import (
	"mini-stock-exchange/internal/repository"

	match_service "mini-stock-exchange/internal/domain-service/match-service"
	order_service "mini-stock-exchange/internal/service/order-service"
	trade_service "mini-stock-exchange/internal/service/trade-service"
	broker_service "mini-stock-exchange/internal/service/broker-service"

	"github.com/go-chi/chi/v5"
)

func NewMockController(r chi.Router, orderRepo repository.OrderRepository, tradeRepo repository.TradeRepository, brokerRepo repository.BrokerRepository) Controller {
	orchestrator := match_service.NewMockOrchestrator(orderRepo)
	tradeService := trade_service.NewTradeService(tradeRepo)
	brokerService := broker_service.NewBrokerService(brokerRepo)
	orderService := order_service.NewOrderService(orderRepo, brokerService, tradeService)
	matchService := match_service.NewMatchService(orchestrator, orderService)

	orderController := NewOrderController(matchService, orderService)
	tradeController := NewTradeController(tradeService)
	healthController := NewHealthController()
	metricController := NewMetricController()
	brokerController := NewBrokerController(brokerService)

	ctrl := controller{
		r:                r,
		orderController:  orderController,
		tradeController:  tradeController,
		healthController: healthController,
		metricController: metricController,
		brokerController: brokerController,
	}
	return &ctrl
}
