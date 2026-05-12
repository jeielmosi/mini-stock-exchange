package match_service

import (
	"context"
	"testing"
	"time"

	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"
	order_service "mini-stock-exchange/internal/service/order-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func mustNewV7() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return id
}

func TestMatchService_SubmitOrder_Match(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := repository.SetupTestDB(ctx)
	assert.NoError(t, err)
	defer cleanup()
	orderRepo, err := repository.NewOrderRepository(db)
	assert.NoError(t, err)
	defer orderRepo.Stop()

	tradeRepo, err := repository.NewTradeRepository(db)
	assert.NoError(t, err)
	defer tradeRepo.Stop()

	orch := NewMockOrchestrator(orderRepo)
	orderSvc := order_service.NewOrderService(orderRepo, tradeRepo)

	svc := NewMatchService(orch, orderSvc)

	symbol := "AAPL"
	bidPrice := float64(150)

	bidOrder := dto.CreateOrderRequest{
		Symbol:     symbol,
		Price:      bidPrice,
		Quantity:   10,
		Type:       entity.Bid,
		ValidUntil: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}

	dto, err := svc.SubmitOrder(context.Background(), bidOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, dto.ID)
}

func TestMatchService_SubmitOrder_PartialMatch(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := repository.SetupTestDB(ctx)
	assert.NoError(t, err)
	defer cleanup()
	orderRepo, err := repository.NewOrderRepository(db)
	assert.NoError(t, err)
	defer orderRepo.Stop()

	tradeRepo, err := repository.NewTradeRepository(db)
	assert.NoError(t, err)
	defer tradeRepo.Stop()

	orderSvc := order_service.NewOrderService(orderRepo, tradeRepo)
	orch := NewMockOrchestrator(orderRepo)
	svc := NewMatchService(orch, orderSvc)

	symbol := "AAPL"
	bidPrice := float64(150)

	bidOrder := dto.CreateOrderRequest{
		Symbol:     symbol,
		Price:      bidPrice,
		Quantity:   10,
		Type:       entity.Bid,
		ValidUntil: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}

	dto, err := svc.SubmitOrder(context.Background(), bidOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, dto.ID)
}

func TestMatchService_SubmitOrder_NoMatch(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := repository.SetupTestDB(ctx)
	assert.NoError(t, err)
	defer cleanup()
	orderRepo, err := repository.NewOrderRepository(db)
	assert.NoError(t, err)
	defer orderRepo.Stop()
	tradeRepo, err := repository.NewTradeRepository(db)
	assert.NoError(t, err)
	defer tradeRepo.Stop()

	orch := NewMockOrchestrator(orderRepo)
	orderSvc := order_service.NewOrderService(orderRepo, tradeRepo)

	svc := NewMatchService(orch, orderSvc)

	symbol := "AAPL"
	bidPrice := float64(100)

	bidOrder := dto.CreateOrderRequest{
		Symbol:     symbol,
		Price:      bidPrice,
		Quantity:   10,
		Type:       entity.Bid,
		ValidUntil: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}

	dto, err := svc.SubmitOrder(context.Background(), bidOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, dto.ID)
}

func TestMatchService_SubmitOrder_MultipleMatches(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := repository.SetupTestDB(ctx)
	assert.NoError(t, err)
	defer cleanup()
	orderRepo, err := repository.NewOrderRepository(db)
	assert.NoError(t, err)
	defer orderRepo.Stop()
	tradeRepo, err := repository.NewTradeRepository(db)
	assert.NoError(t, err)
	defer tradeRepo.Stop()

	orch := NewMockOrchestrator(orderRepo)
	orderSvc := order_service.NewOrderService(orderRepo, tradeRepo)

	svc := NewMatchService(orch, orderSvc)

	symbol := "AAPL"
	bidPrice := float64(150)

	bidOrder := dto.CreateOrderRequest{
		Symbol:     symbol,
		Price:      bidPrice,
		Quantity:   10,
		Type:       entity.Bid,
		ValidUntil: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}

	dto, err := svc.SubmitOrder(context.Background(), bidOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, dto.ID)
}

func TestMatchService_SubmitOrder_AskMatch(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := repository.SetupTestDB(ctx)
	assert.NoError(t, err)
	defer cleanup()
	orderRepo, err := repository.NewOrderRepository(db)
	assert.NoError(t, err)
	defer orderRepo.Stop()

	tradeRepo, err := repository.NewTradeRepository(db)
	assert.NoError(t, err)
	defer tradeRepo.Stop()

	orch := NewMockOrchestrator(orderRepo)
	orderSvc := order_service.NewOrderService(orderRepo, tradeRepo)

	svc := NewMatchService(orch, orderSvc)

	symbol := "AAPL"
	askPrice := float64(130)

	askOrder := dto.CreateOrderRequest{
		Symbol:     symbol,
		Price:      askPrice,
		Quantity:   10,
		Type:       entity.Ask,
		ValidUntil: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}

	dto, err := svc.SubmitOrder(context.Background(), askOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, dto.ID)
}
