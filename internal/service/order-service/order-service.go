package order_service

import (
	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"
	broker_service "mini-stock-exchange/internal/service/broker-service"
	trade_service "mini-stock-exchange/internal/service/trade-service"
	"time"

	"github.com/google/uuid"
)

type OrderService interface {
	GetOrder(req dto.GetOrderRequest) (dto.GetOrderResponse, error)
	PostOrder(order entity.Order) error
}

type orderService struct {
	orderRepo repository.OrderRepository
	brokerSvc broker_service.BrokerService
	tradeSvc  trade_service.TradeService
}

func NewOrderService(orderRepo repository.OrderRepository, brokerSvc broker_service.BrokerService, tradeSvc trade_service.TradeService) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		brokerSvc: brokerSvc,
		tradeSvc:  tradeSvc,
	}
}

func (o *orderService) GetOrder(req dto.GetOrderRequest) (dto.GetOrderResponse, error) {
	now := time.Now()
	order, err := o.orderRepo.GetByID(req.ID)
	if err != nil {
		return dto.GetOrderResponse{}, err
	}
	if order.ValidUntil.Before(now) {
		o.orderRepo.Expire([]uuid.UUID{order.ID})
	}

	broker, err := o.brokerSvc.GetBrokerByID(order.BrokerID)
	if err != nil {
		return dto.GetOrderResponse{}, err
	}
	order.BrokerName = broker.Name

	trades, err := o.tradeSvc.GetTradesByOrder(req.ID)
	if err != nil {
		return dto.GetOrderResponse{}, err
	}
	res, err := dto.NewGetOrderResponse(order, trades)
	if err != nil {
		return dto.GetOrderResponse{}, err
	}
	return res, nil
}

func (o *orderService) PostOrder(order entity.Order) error {
	return o.orderRepo.Insert(order)
}
