package match_service

import (
	"sync"

	"mini-stock-exchange/internal/repository"
	"mini-stock-exchange/internal/usecase"
)

func NewMockOrchestrator(orderRepo repository.OrderRepository) Orchestrator {
	return &orchestrator{
		executors:    make(map[string]Executor),
		orderRepo:    orderRepo,
		mu:           sync.RWMutex{},
		matchUsecase: usecase.NewOrderMatchUsecase(),
		createTrade:  usecase.NewCreateTradeUsecase(),
	}
}
