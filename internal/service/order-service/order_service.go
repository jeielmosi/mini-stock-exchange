package order_service

import (
	"context"
	"fmt"
	"time"

	"mini-stock-exchange/internal/domain"
	"mini-stock-exchange/internal/observability"
	"mini-stock-exchange/internal/repository"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("order-service")

type OrderService interface {
	SubmitOrder(ctx context.Context, order domain.Order) error
	GetOrder(id uuid.UUID) (domain.Order, error)
}

type orderService struct {
	orderRepo   repository.OrderRepository
	tradeRepo   repository.TradeRepository
	centralizer Orchestrator
}

func NewOrderService(orderRepo repository.OrderRepository, tradeRepo repository.TradeRepository, centralizer Orchestrator) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		tradeRepo:   tradeRepo,
		centralizer: centralizer,
	}
}

func (s *orderService) SubmitOrder(ctx context.Context, order domain.Order) error {
	ctx, span := tracer.Start(ctx, "SubmitOrder")
	defer span.End()

	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate UUID v7: %w", err)
	}
	order.ID = id
	order.CreatedAt = time.Now()
	if order.RemainingQuantity == 0 {
		order.RemainingQuantity = order.Quantity
	}
	order.Status = domain.Pending

	observability.OrdersSubmitted.WithLabelValues(string(order.Type), order.Symbol).Inc()
	observability.ActiveOrders.Inc()

	if err := s.centralizer.RouteOrder(ctx, order); err != nil {
		return fmt.Errorf("failed to route order: %w", err)
	}

	return nil
}

func (s *orderService) GetOrder(id uuid.UUID) (domain.Order, error) {
	return s.orderRepo.GetByID(id)
}
