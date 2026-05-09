package order_service

import (
	"context"
	"fmt"
	"time"

	"mini-stock-exchange/internal/domain"
	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/observability"
	"mini-stock-exchange/internal/repository"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("order-service")

type OrderService interface {
	SubmitOrder(ctx context.Context, order domain.Order) (dto.CreateOrderResponse, error)
	GetOrder(req dto.GetOrderRequest) (dto.GetOrderResponse, error)
}

type orderService struct {
	orderRepo    repository.OrderRepository
	tradeRepo    repository.TradeRepository
	orchestrator Orchestrator
}

func NewOrderService(orderRepo repository.OrderRepository, tradeRepo repository.TradeRepository, orchestrator Orchestrator) OrderService {
	return &orderService{
		orderRepo:    orderRepo,
		tradeRepo:    tradeRepo,
		orchestrator: orchestrator,
	}
}

// TODO should recieve dto and not order
func (s *orderService) SubmitOrder(ctx context.Context, order domain.Order) (dto.CreateOrderResponse, error) {
	ctx, span := tracer.Start(ctx, "SubmitOrder")
	defer span.End()

	id, err := uuid.NewV7()
	if err != nil {
		return dto.CreateOrderResponse{}, fmt.Errorf("failed to generate UUIDv7: %w", err)
	}
	order.ID = id
	order.CreatedAt = time.Now()
	if order.RemainingQuantity == 0 {
		order.RemainingQuantity = order.Quantity
	}
	order.Status = domain.Pending

	observability.OrdersSubmitted.WithLabelValues(string(order.Type), order.Symbol).Inc()
	observability.ActiveOrders.Inc()

	if err := s.orchestrator.RouteOrder(ctx, order); err != nil {
		return dto.CreateOrderResponse{}, fmt.Errorf("failed to route order: %w", err)
	}

	return dto.CreateOrderResponse{ID: order.ID.String()}, nil
}

func (s *orderService) GetOrder(req dto.GetOrderRequest) (dto.GetOrderResponse, error) {
	id, err := req.ToUUID()
	if err != nil {
		return dto.GetOrderResponse{}, err
	}
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return dto.GetOrderResponse{}, err
	}
	return dto.NewGetOrderResponse(order), nil
}
