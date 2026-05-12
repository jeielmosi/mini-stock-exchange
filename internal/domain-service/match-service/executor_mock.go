package match_service

import (
	"mini-stock-exchange/internal/config"
	"mini-stock-exchange/internal/repository"
	"mini-stock-exchange/internal/usecase"
)

func NewMockExecutor(symbol string, capacity int, orderRepo repository.OrderRepository) Executor {
	config.LoadTest(capacity)
	return NewExecutor(symbol, orderRepo, usecase.NewOrderMatchUsecase(), usecase.NewCreateTradeUsecase())
}
