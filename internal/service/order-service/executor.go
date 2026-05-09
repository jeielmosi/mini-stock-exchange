package order_service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"mini-stock-exchange/internal/domain"
	"mini-stock-exchange/internal/observability"
	"mini-stock-exchange/internal/repository"
	order_heaps "mini-stock-exchange/internal/service/order-service/order-heaps"

	"github.com/google/uuid"
)

type Executor interface {
	ProcessOrder(ctx context.Context, order domain.Order) error
	Stop()
}

type executor struct {
	symbol    string
	bidHeap   order_heaps.OrderHeap
	askHeap   order_heaps.OrderHeap
	orderRepo repository.OrderRepository
	tradeRepo repository.TradeRepository
	mu        sync.Mutex
	orderChan chan domain.Order
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewExecutor(symbol string, orderRepo repository.OrderRepository, tradeRepo repository.TradeRepository) (Executor, error) {
	bidHeap, err := order_heaps.NewBidHeap(symbol, orderRepo)
	if err != nil {
		return nil, err
	}
	askHeap, err := order_heaps.NewAskHeap(symbol, orderRepo)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	e := &executor{
		symbol:    symbol,
		bidHeap:   bidHeap,
		askHeap:   askHeap,
		orderRepo: orderRepo,
		tradeRepo: tradeRepo,
		orderChan: make(chan domain.Order, 1024),
		ctx:       ctx,
		cancel:    cancel,
	}

	go e.run()

	return e, nil
}

func (e *executor) ProcessOrder(ctx context.Context, order domain.Order) error {
	if err := e.orderRepo.Insert(order); err != nil {
		return fmt.Errorf("failed to save order to database: %w", err)
	}

	select {
	case e.orderChan <- order:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-e.ctx.Done():
		return fmt.Errorf("executor for symbol %s is stopped", e.symbol)
	}
}

func (e *executor) Stop() {
	e.cancel()
}

func (e *executor) run() {
	for {
		select {
		case <-e.ctx.Done():
			return
		case order := <-e.orderChan:
			e.matchMake(order)
		}
	}
}

func (e *executor) matchMake(order domain.Order) {
	e.mu.Lock()
	defer e.mu.Unlock()

	start := time.Now()
	defer func() {
		observability.MatchingLatency.Observe(time.Since(start).Seconds())
	}()

	var orderHeap order_heaps.OrderHeap
	var matchHeap order_heaps.OrderHeap
	switch order.Type {
	case domain.Bid:
		orderHeap = e.bidHeap
		matchHeap = e.askHeap
	case domain.Ask:
		orderHeap = e.askHeap
		matchHeap = e.bidHeap
	default:
		return
	}

	dto := matchHeap.MatchMake(order)

	for _, exp := range dto.Expired {
		if err := e.orderRepo.Update(exp); err != nil {
			slog.Error("failed to update expired order", "error", err, "order_id", exp.ID)
			panic(err)
		}
	}

	for _, match := range dto.Matches {
		if err := e.match(&order, &match); err != nil {
			slog.Error("failed to match order", "error", err, "order_id", order.ID, "match_id", match.ID)
		}
		if 0 < match.RemainingQuantity {
			matchHeap.Push(match)
		}
	}

	if 0 < order.RemainingQuantity {
		orderHeap.Push(order)
	}
}

func (e *executor) match(bid *domain.Order, ask *domain.Order) error {
	if bid == nil || ask == nil {
		return fmt.Errorf("nil order")
	}
	isBid := (bid.Type == domain.Bid)
	isAsk := (ask.Type == domain.Ask)
	if !isBid && !isAsk {
		return e.match(ask, bid)
	}
	if !isAsk || !isBid {
		return fmt.Errorf("invalid order types for matching")
	}
	if bid.Price.LessThan(ask.Price) {
		return nil
	}

	tradeQty := e.calculateTradeQuantity(bid.RemainingQuantity, ask.RemainingQuantity)

	tradeID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate UUID v7 for trade: %w", err)
	}

	trade := domain.Trade{
		ID:          tradeID,
		Symbol:      e.symbol,
		Price:       ask.Price,
		Quantity:    tradeQty,
		ExecutedAt:  time.Now(),
		BuyOrderID:  bid.ID,
		SellOrderID: ask.ID,
	}

	if err := e.tradeRepo.Create(trade); err != nil {
		slog.Error("failed to create trade", "error", err, "trade_id", trade.ID)
		return fmt.Errorf("failed to create trade: %w", err)
	}

	observability.TradesExecuted.Inc()

	for _, order := range []*domain.Order{bid, ask} {
		order.RemainingQuantity -= tradeQty
		if order.RemainingQuantity == 0 {
			order.Status = domain.Filled
			observability.ActiveOrders.Dec()
		} else {
			order.Status = domain.Partial
		}
		if err := e.orderRepo.Update(*order); err != nil {
			return fmt.Errorf("failed to update order: %w", err)
		}
	}

	return nil
}

func (e *executor) calculateTradeQuantity(q1, q2 int) int {
	if q1 < q2 {
		return q1
	}
	return q2
}
