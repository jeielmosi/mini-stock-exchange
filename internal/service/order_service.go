package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"mini-stock-exchange/internal/domain"
	"mini-stock-exchange/internal/observability"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("order-service")

type OrderService struct {
	orderRepo domain.OrderRepository
	tradeRepo domain.TradeRepository
	orderChan chan *domain.Order
}

func NewOrderService(orderRepo domain.OrderRepository, tradeRepo domain.TradeRepository) *OrderService {
	s := &OrderService{
		orderRepo: orderRepo,
		tradeRepo: tradeRepo,
		orderChan: make(chan *domain.Order, 100),
	}

	go s.processOrders()

	return s
}

func (s *OrderService) SubmitOrder(ctx context.Context, order *domain.Order) error {
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

	if err := s.orderRepo.Create(order); err != nil {
		slog.Error("failed to create order", "error", err, "order_id", order.ID)
		return fmt.Errorf("failed to create order: %w", err)
	}

	observability.OrdersSubmitted.WithLabelValues(string(order.Type), order.Symbol).Inc()
	observability.ActiveOrders.Inc()

	s.orderChan <- order

	return nil
}

func (s *OrderService) processOrders() {
	for order := range s.orderChan {
		if err := s.matchOrder(context.Background(), order); err != nil {
			slog.Error("failed to match order in background", "error", err, "order_id", order.ID)
		}
	}
}

func (s *OrderService) matchOrder(ctx context.Context, order *domain.Order) error {
	ctx, span := tracer.Start(ctx, "matchOrder")
	defer span.End()

	start := time.Now()
	defer func() {
		observability.MatchingLatency.Observe(time.Since(start).Seconds())
	}()

	for order.RemainingQuantity > 0 {
		matches, err := s.orderRepo.FindMatches(*order)
		if err != nil {
			slog.Error("failed to find matches", "error", err, "symbol", order.Symbol)
			return fmt.Errorf("failed to find matches: %w", err)
		}

		if len(matches) == 0 {
			break
		}

		match := matches[0]
		tradeQty := s.calculateTradeQuantity(order.RemainingQuantity, match.RemainingQuantity)

		var executionPrice decimal.Decimal
		if order.Type == domain.Bid {
			executionPrice = match.Price
		} else {
			executionPrice = order.Price
		}

		tradeID, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("failed to generate UUID v7 for trade: %w", err)
		}

		trade := &domain.Trade{
			ID:         tradeID,
			Symbol:     order.Symbol,
			Price:      executionPrice,
			Quantity:   tradeQty,
			ExecutedAt: time.Now(),
		}

		if order.Type == domain.Bid {
			trade.BuyOrderID = order.ID
			trade.SellOrderID = match.ID
		} else {
			trade.BuyOrderID = match.ID
			trade.SellOrderID = order.ID
		}

		if err := s.tradeRepo.Create(trade); err != nil {
			slog.Error("failed to create trade", "error", err, "trade_id", trade.ID)
			return fmt.Errorf("failed to create trade: %w", err)
		}

		observability.TradesExecuted.Inc()

		order.RemainingQuantity -= tradeQty
		if order.RemainingQuantity == 0 {
			order.Status = domain.Filled
			observability.ActiveOrders.Dec()
		} else {
			order.Status = domain.Partial
		}
		if err := s.orderRepo.Update(order); err != nil {
			return fmt.Errorf("failed to update order: %w", err)
		}

		match.RemainingQuantity -= tradeQty
		if match.RemainingQuantity == 0 {
			match.Status = domain.Filled
			observability.ActiveOrders.Dec()
		} else {
			match.Status = domain.Partial
		}
		if err := s.orderRepo.Update(&match); err != nil {
			return fmt.Errorf("failed to update match: %w", err)
		}
	}

	return nil
}

func (s *OrderService) calculateTradeQuantity(q1, q2 int) int {
	if q1 < q2 {
		return q1
	}
	return q2
}

func (s *OrderService) GetOrderStatus(id uuid.UUID) (*domain.Order, error) {
	return s.orderRepo.GetByID(id)
}
