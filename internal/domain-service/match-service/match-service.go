package match_service

import (
	"context"
	"fmt"

	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/observability"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("match-service")

type MatchService interface {
	SubmitOrder(ctx context.Context, req dto.CreateOrderRequest) (dto.CreateOrderResponse, error)
}

type matchService struct {
	orchestrator Orchestrator
}

func NewMatchService(orchestrator Orchestrator) MatchService {
	return &matchService{
		orchestrator: orchestrator,
	}
}

// TODO should recieve dto and not order
func (s *matchService) SubmitOrder(ctx context.Context, req dto.CreateOrderRequest) (dto.CreateOrderResponse, error) {
	ctx, span := tracer.Start(ctx, "SubmitOrder")
	defer span.End()

	order, err := req.ToOrder()
	if err != nil {
		return dto.CreateOrderResponse{}, err
	}

	observability.OrdersSubmitted.WithLabelValues(string(order.Type), order.Symbol).Inc()
	observability.ActiveOrders.Inc()

	if err := s.orchestrator.RouteOrder(ctx, order); err != nil {
		return dto.CreateOrderResponse{}, fmt.Errorf("failed to route order: %w", err)
	}

	return dto.NewCreateOrderResponse(order.ID)
}
