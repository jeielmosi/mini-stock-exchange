package match_service

import (
	"context"
	"testing"
	"time"

	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"
	order_service "mini-stock-exchange/internal/service/order-service"
	broker_service "mini-stock-exchange/internal/service/broker-service"
	trade_service "mini-stock-exchange/internal/service/trade-service"
	dto_helper "mini-stock-exchange/internal/dto/helper"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	brokerRepo, err := repository.NewBrokerRepository(db)
	assert.NoError(t, err)
	brokerID := mustNewV7()
	err = brokerRepo.Insert(entity.Broker{ID: brokerID, Name: "broker1"})
	assert.NoError(t, err)

	orch := NewMockOrchestrator(orderRepo)
	brokerSvc := broker_service.NewBrokerService(brokerRepo)
	tradeSvc := trade_service.NewTradeService(tradeRepo)
	orderSvc := order_service.NewOrderService(orderRepo, brokerSvc, tradeSvc)

	svc := NewMatchService(orch, orderSvc)

	symbol := "AAPL"
	bidPrice := float64(150)

	bidOrderEncodedID, err := dto_helper.EncodeUUID(brokerID)
	require.NoError(t, err)

	bidOrder := dto.CreateOrderRequest{
		BrokerID:   bidOrderEncodedID,
		Symbol:     symbol,
		Price:      bidPrice,
		Quantity:   10,
		Type:       entity.Bid,
		ValidUntil: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}

	resp, err := svc.SubmitOrder(context.Background(), bidOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.ID)

	decodedID, err := dto_helper.DecodeUUID(resp.ID)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, decodedID)
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

	brokerRepo, err := repository.NewBrokerRepository(db)
	assert.NoError(t, err)
	brokerID := mustNewV7()
	err = brokerRepo.Insert(entity.Broker{ID: brokerID, Name: "broker1"})
	assert.NoError(t, err)

	brokerSvc := broker_service.NewBrokerService(brokerRepo)
	tradeSvc := trade_service.NewTradeService(tradeRepo)
	orderSvc := order_service.NewOrderService(orderRepo, brokerSvc, tradeSvc)
	orch := NewMockOrchestrator(orderRepo)
	svc := NewMatchService(orch, orderSvc)

	symbol := "AAPL"
	bidPrice := float64(150)

	bidOrderEncodedID, err := dto_helper.EncodeUUID(brokerID)
	require.NoError(t, err)

	bidOrder := dto.CreateOrderRequest{
		BrokerID:   bidOrderEncodedID,
		Symbol:     symbol,
		Price:      bidPrice,
		Quantity:   10,
		Type:       entity.Bid,
		ValidUntil: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}

	resp, err := svc.SubmitOrder(context.Background(), bidOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.ID)

	decodedID, err := dto_helper.DecodeUUID(resp.ID)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, decodedID)
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

	brokerRepo, err := repository.NewBrokerRepository(db)
	assert.NoError(t, err)
	brokerID := mustNewV7()
	err = brokerRepo.Insert(entity.Broker{ID: brokerID, Name: "broker1"})
	assert.NoError(t, err)

	orch := NewMockOrchestrator(orderRepo)
	brokerSvc := broker_service.NewBrokerService(brokerRepo)
	tradeSvc := trade_service.NewTradeService(tradeRepo)
	orderSvc := order_service.NewOrderService(orderRepo, brokerSvc, tradeSvc)

	svc := NewMatchService(orch, orderSvc)

	symbol := "AAPL"
	bidPrice := float64(100)

	bidOrderEncodedID, err := dto_helper.EncodeUUID(brokerID)
	require.NoError(t, err)

	bidOrder := dto.CreateOrderRequest{
		BrokerID:   bidOrderEncodedID,
		Symbol:     symbol,
		Price:      bidPrice,
		Quantity:   10,
		Type:       entity.Bid,
		ValidUntil: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}

	resp, err := svc.SubmitOrder(context.Background(), bidOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.ID)

	decodedID, err := dto_helper.DecodeUUID(resp.ID)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, decodedID)
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

	brokerRepo, err := repository.NewBrokerRepository(db)
	assert.NoError(t, err)
	brokerID := mustNewV7()
	err = brokerRepo.Insert(entity.Broker{ID: brokerID, Name: "broker1"})
	assert.NoError(t, err)

	orch := NewMockOrchestrator(orderRepo)
	brokerSvc := broker_service.NewBrokerService(brokerRepo)
	tradeSvc := trade_service.NewTradeService(tradeRepo)
	orderSvc := order_service.NewOrderService(orderRepo, brokerSvc, tradeSvc)

	svc := NewMatchService(orch, orderSvc)

	symbol := "AAPL"
	bidPrice := float64(150)

	bidOrderEncodedID, err := dto_helper.EncodeUUID(brokerID)
	require.NoError(t, err)

	bidOrder := dto.CreateOrderRequest{
		BrokerID:   bidOrderEncodedID,
		Symbol:     symbol,
		Price:      bidPrice,
		Quantity:   10,
		Type:       entity.Bid,
		ValidUntil: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}

	resp, err := svc.SubmitOrder(context.Background(), bidOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.ID)

	decodedID, err := dto_helper.DecodeUUID(resp.ID)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, decodedID)
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

	brokerRepo, err := repository.NewBrokerRepository(db)
	assert.NoError(t, err)
	brokerID := mustNewV7()
	err = brokerRepo.Insert(entity.Broker{ID: brokerID, Name: "broker1"})
	assert.NoError(t, err)

	orch := NewMockOrchestrator(orderRepo)
	brokerSvc := broker_service.NewBrokerService(brokerRepo)
	tradeSvc := trade_service.NewTradeService(tradeRepo)
	orderSvc := order_service.NewOrderService(orderRepo, brokerSvc, tradeSvc)

	svc := NewMatchService(orch, orderSvc)

	symbol := "AAPL"
	askPrice := float64(130)

	askOrderEncodedID, err := dto_helper.EncodeUUID(brokerID)
	require.NoError(t, err)

	askOrder := dto.CreateOrderRequest{
		BrokerID:   askOrderEncodedID,
		Symbol:     symbol,
		Price:      askPrice,
		Quantity:   10,
		Type:       entity.Ask,
		ValidUntil: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}

	resp, err := svc.SubmitOrder(context.Background(), askOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.ID)

	decodedID, err := dto_helper.DecodeUUID(resp.ID)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, decodedID)
}
