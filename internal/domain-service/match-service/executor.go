package match_service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"mini-stock-exchange/internal/config"
	order_heap "mini-stock-exchange/internal/domain-service/match-service/order-heap"
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/observability"
	"mini-stock-exchange/internal/repository"
	"mini-stock-exchange/internal/usecase"
)

type Executor interface {
	ProcessOrder(ctx context.Context, order entity.Order) error
	Stop()
}

type executor struct {
	symbol    string
	bidHeap   order_heap.OrderHeap
	askHeap   order_heap.OrderHeap
	orderRepo repository.OrderRepository
	mu        sync.Mutex
	orderChan chan entity.Order
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewExecutor(symbol string, orderRepo repository.OrderRepository) Executor {
	bidHeap := order_heap.NewBidHeap(symbol, config.ENV.ExecutorCapacity, orderRepo)
	askHeap := order_heap.NewAskHeap(symbol, config.ENV.ExecutorCapacity, orderRepo)

	ctx, cancel := context.WithCancel(context.Background())

	e := &executor{
		symbol:    symbol,
		bidHeap:   bidHeap,
		askHeap:   askHeap,
		orderRepo: orderRepo,
		orderChan: make(chan entity.Order, config.ENV.ExecutorCapacity),
		ctx:       ctx,
		cancel:    cancel,
	}

	go e.run()

	return e
}

func (e *executor) ProcessOrder(ctx context.Context, order entity.Order) error {
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

func (e *executor) matchMake(order entity.Order) {
	e.mu.Lock()
	defer e.mu.Unlock()

	start := time.Now()
	defer func() {
		observability.MatchingLatency.Observe(time.Since(start).Seconds())
	}()

	var orderHeap order_heap.OrderHeap
	var matchHeap order_heap.OrderHeap
	switch order.Type {
	case entity.Bid:
		orderHeap = e.bidHeap
		matchHeap = e.askHeap
	case entity.Ask:
		orderHeap = e.askHeap
		matchHeap = e.bidHeap
	default:
		return
	}

	dto, err := matchHeap.Pop(order)
	if err != nil {
		slog.Error("failed to pop order from heap", "error", err, "order_id", order.ID)
		return
	}

	if err := e.orderRepo.Expire(dto.Expired); err != nil {
		slog.Error("failed to update expired order", "error", err)
		return
	}

	for _, match := range dto.Matches {
		if err := e.match(&order, &match); err != nil {
			slog.Error("failed to match order", "error", err, "order_id", order.ID, "match_id", match.ID)
			panic(err)
		}
		if 0 < match.RemainingQuantity {
			matchHeap.Push(match)
		}
	}

	if 0 < order.RemainingQuantity {
		orderHeap.Push(order)
	}
}

func (e *executor) match(bid *entity.Order, ask *entity.Order) error {
	trade, err := usecase.NewTrade(bid, ask)
	if err != nil {
		return fmt.Errorf("failed to create trade: %w", err)
	}

	err = e.orderRepo.Match(e.ctx, repository.MatchDTO{
		Ask:   *ask,
		Bid:   *bid,
		Trade: trade,
	})
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	observability.TradesExecuted.Inc()

	return nil
}
