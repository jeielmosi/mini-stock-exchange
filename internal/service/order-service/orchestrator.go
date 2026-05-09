package order_service

import (
	"context"
	"sync"

	"mini-stock-exchange/internal/domain"
	"mini-stock-exchange/internal/repository"
)

type Orchestrator interface {
	RouteOrder(ctx context.Context, order domain.Order) error
}

type orchestrator struct {
	executors map[string]Executor
	orderRepo repository.OrderRepository
	tradeRepo repository.TradeRepository
	mu        sync.RWMutex
}

func NewOrchestrator(orderRepo repository.OrderRepository, tradeRepo repository.TradeRepository) Orchestrator {
	return &orchestrator{
		executors: make(map[string]Executor),
		orderRepo: orderRepo,
		tradeRepo: tradeRepo,
		mu:        sync.RWMutex{},
	}
}

func (o *orchestrator) RouteOrder(ctx context.Context, order domain.Order) error {
	o.mu.RLock()
	executor, ok := o.executors[order.Symbol]
	o.mu.RUnlock()

	if !ok {
		o.mu.Lock()
		// Double-check after acquiring write lock
		executor, ok = o.executors[order.Symbol]
		if !ok {
				var err error
				executor, err = NewExecutor(order.Symbol, o.orderRepo, o.tradeRepo)
				if err != nil {
					o.mu.Unlock()
					return err
				}
				o.executors[order.Symbol] = executor
		}
		o.mu.Unlock()
	}

	return executor.ProcessOrder(ctx, order)
}
