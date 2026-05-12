package order_service

import (
	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"
)

// TODO get the order by ID, with trades
type OrderService interface {
	GetOrder(req dto.GetOrderRequest) (dto.GetOrderResponse, error)
	PostOrder(order entity.Order) error
}

type orderService struct {
	orderRepo repository.OrderRepository
	tradeRepo repository.TradeRepository
}

func NewOrderService(orderRepo repository.OrderRepository, tradeRepo repository.TradeRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		tradeRepo: tradeRepo,
	}
}

func (o *orderService) GetOrder(req dto.GetOrderRequest) (dto.GetOrderResponse, error) {
	order, err := o.orderRepo.GetByID(req.ID)
	if err != nil {
		return dto.GetOrderResponse{}, err
	}
	trades, err := o.tradeRepo.GetByOrderID(req.ID)
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
