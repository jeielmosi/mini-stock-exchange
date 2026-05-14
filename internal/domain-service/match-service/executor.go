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
	"mini-stock-exchange/internal/utils"

	"github.com/google/uuid"
)

type Executor interface {
	ProcessOrder(ctx context.Context, order *entity.Order) error
	Stop()
}

type executor struct {
	symbol       string
	bidHeap      *utils.PriorityQueue[*entity.Order]
	askHeap      *utils.PriorityQueue[*entity.Order]
	orderRepo    repository.OrderRepository
	matchUsecase usecase.OrderMatchUsecase
	createTrade  usecase.CreateTradeUsecase
	mu           sync.Mutex
	orderChan    chan *entity.Order
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewExecutor(
	symbol string, orderRepo repository.OrderRepository,
	matchUsecase usecase.OrderMatchUsecase, createTrade usecase.CreateTradeUsecase,
) Executor {
	bidHeap := order_heap.NewBidHeap(config.ENV.ExecutorCapacity)
	askHeap := order_heap.NewAskHeap(config.ENV.ExecutorCapacity)

	ctx, cancel := context.WithCancel(context.Background())

	e := &executor{
		symbol:       symbol,
		bidHeap:      bidHeap,
		askHeap:      askHeap,
		orderRepo:    orderRepo,
		matchUsecase: matchUsecase,
		createTrade:  createTrade,
		orderChan:    make(chan *entity.Order, config.ENV.ExecutorCapacity),
		ctx:          ctx,
		cancel:       cancel,
	}

	go e.run()

	return e
}

func (e *executor) ProcessOrder(ctx context.Context, order *entity.Order) error {
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

func (e *executor) matchMake(order *entity.Order) {
	e.mu.Lock()
	defer e.mu.Unlock()

	start := time.Now()
	defer func() {
		observability.MatchingLatency.Observe(time.Since(start).Seconds())
	}()

	var OrderHeap *utils.PriorityQueue[*entity.Order]
	var matchHeap *utils.PriorityQueue[*entity.Order]
	switch order.Type {
	case entity.Bid:
		OrderHeap = e.bidHeap
		matchHeap = e.askHeap
	case entity.Ask:
		OrderHeap = e.askHeap
		matchHeap = e.bidHeap
	default:
		return
	}

	retryArr := []*entity.Order{}
	expiredArr := []uuid.UUID{}

	for 0 < order.RemainingQuantity {
		match, ok := matchHeap.Peek()
		if !ok {
			var err error
			var orders []entity.Order
			if order.Type == entity.Bid {
				orders, err = e.orderRepo.GetAsks(e.symbol, e.askHeap.Cap())
			}
			if order.Type == entity.Ask {
				orders, err = e.orderRepo.GetBids(e.symbol, e.bidHeap.Cap())
			}
			if err != nil {
				slog.Error("failed to get orders from database", "error", err)
				break
			}
			if len(orders) == 0 {
				break
			}
			for _, order := range orders {
				matchHeap.Push(&order)
			}
			continue
		}

		dto, err := e.matchUsecase.ValidateMatch(order, match)
		if err != nil {
			slog.Error("failed to validate match", "error", err)
			break
		}
		matchHeap.Drop()
		if dto.Retry != nil {
			retryArr = append(retryArr, dto.Retry)
		}
		if dto.Expired != nil {
			expiredArr = append(expiredArr, dto.Expired.ID)
		}
		if dto.Ask == nil || dto.Bid == nil {
			continue
		}

		if err := e.match(order, match); err != nil {
			slog.Error("failed to match order", "error", err, "order_id", order.ID, "match_id", match.ID)
			break
		}
		if 0 < match.RemainingQuantity {
			matchHeap.Push(match)
			break
		}
	}

	if err := e.orderRepo.Expire(expiredArr); err != nil {
		slog.Error("failed to update expired order", "error", err)
	}

	if 0 < order.RemainingQuantity {
		OrderHeap.Push(order)
	}
}

func (e *executor) match(order *entity.Order, match *entity.Order) error {
	dto, err := e.matchUsecase.MatchOrder(order, match)
	if err != nil {
		return fmt.Errorf("failed to match order: %w", err)
	}
	trade, err := e.createTrade.CreateTrade(*dto)
	if err != nil {
		return fmt.Errorf("failed to create trade: %w", err)
	}
	defer func() {
		if err != nil {
			e.matchUsecase.UnmatchOrder(*dto)
		}
	}()

	err = e.orderRepo.Match(e.ctx, repository.MatchDTO{
		Ask:   *dto.Ask,
		Bid:   *dto.Bid,
		Trade: trade,
	})
	if err != nil {
		return fmt.Errorf("failed to save match to database: %w", err)
	}

	observability.TradesExecuted.Inc()

	return nil
}
