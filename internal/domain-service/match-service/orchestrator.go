package match_service

import (
	"context"
	"sync"

	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"
)

type Orchestrator interface {
	RouteOrder(ctx context.Context, order entity.Order) error
}

var (
	orch Orchestrator = nil
	once sync.Once
)

type orchestrator struct {
	executors map[string]Executor
	orderRepo repository.OrderRepository
	mu        sync.RWMutex
}

func NewOrchestrator(orderRepo repository.OrderRepository) Orchestrator {
	once.Do(func() {
		orch = &orchestrator{
			executors: make(map[string]Executor),
			orderRepo: orderRepo,
			mu:        sync.RWMutex{},
		}
	})
	return orch
}

func (o *orchestrator) RouteOrder(ctx context.Context, order entity.Order) error {
	o.mu.RLock()
	executor, ok := o.executors[order.Symbol]
	o.mu.RUnlock()

	if !ok {
		o.mu.Lock()
		// Double-check after acquiring write lock
		executor, ok = o.executors[order.Symbol]
		if !ok {
			executor = NewExecutor(order.Symbol, o.orderRepo)
			o.executors[order.Symbol] = executor
		}
		o.mu.Unlock()
	}

	return executor.ProcessOrder(ctx, order)
}
