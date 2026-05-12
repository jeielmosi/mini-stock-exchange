package order_service

import (
	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/repository"
)

// TODO get the order by ID, with trades
type OrderService interface {
	GetOrder(req dto.GetOrderRequest) (dto.GetOrderResponse, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
	//tradeRepo repository.TradeRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		//tradeRepo: tradeRepo,
	}
}

// TODO
func (o *orderService) GetOrder(req dto.GetOrderRequest) (dto.GetOrderResponse, error) {
	order, err := o.orderRepo.GetByID(req.ID)
	if err != nil {
		return dto.GetOrderResponse{}, err
	}
	res, err := dto.NewGetOrderResponse(order)
	if err != nil {
		return dto.GetOrderResponse{}, err
	}
	return res, nil
}
