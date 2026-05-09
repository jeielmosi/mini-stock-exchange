package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mini-stock-exchange/internal/domain"
	"mini-stock-exchange/internal/handler"
	"mini-stock-exchange/internal/repository"
	order_service "mini-stock-exchange/internal/service/order-service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer() (*httptest.Server, repository.OrderRepository, repository.TradeRepository, func()) {
	ctx := context.Background()
	db, cleanup := repository.SetupTestDB(ctx)

	orderRepo := repository.NewOrderRepository(db)
	tradeRepo := repository.NewTradeRepository(db)
	orchestrator := order_service.NewOrchestrator(orderRepo, tradeRepo)
	orderService := order_service.NewOrderService(orderRepo, tradeRepo, orchestrator)
	orderHandler := handler.NewOrderHandler(orderService)

	r := chi.NewRouter()
	orderHandler.RegisterRoutes(r)

	server := httptest.NewServer(r)

	return server, orderRepo, tradeRepo, func() {
		server.Close()
		cleanup()
	}
}

func TestOrderFlow(t *testing.T) {
	server, orderRepo, tradeRepo, cleanup := setupTestServer()
	defer cleanup()

	symbol := "AAPL"
	bidPrice := decimal.NewFromInt(150)
	askPrice := decimal.NewFromInt(140)

	// 1. Submit an Ask order
	askRequest := map[string]interface{}{
		"broker_id":   "broker1",
		"owner_doc":   "doc1",
		"type":        "ASK",
		"symbol":      symbol,
		"price":       askPrice.String(),
		"quantity":    10,
		"valid_until": time.Now().Add(time.Hour).Format(time.RFC3339),
	}
	body, _ := json.Marshal(askRequest)
	resp, err := http.Post(server.URL+"/orders", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var askResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&askResp)
	require.NoError(t, err)
	askID, ok := askResp["id"]
	require.True(t, ok)

	// Verify Ask order is in DB
	askOrder, err := orderRepo.GetByID(uuid.MustParse(askID))
	require.NoError(t, err)
	assert.Equal(t, symbol, askOrder.Symbol)
	assert.Equal(t, domain.Ask, askOrder.Type)
	assert.Equal(t, 10, askOrder.Quantity)

	// 2. Submit a matching Bid order
	bidRequest := map[string]interface{}{
		"broker_id":   "broker2",
		"owner_doc":   "doc2",
		"type":        "BID",
		"symbol":      symbol,
		"price":       bidPrice.String(),
		"quantity":    10,
		"valid_until": time.Now().Add(time.Hour).Format(time.RFC3339),
	}
	body, _ = json.Marshal(bidRequest)
	resp, err = http.Post(server.URL+"/orders", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var bidResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&bidResp)
	require.NoError(t, err)
	bidID, ok := bidResp["id"]
	require.True(t, ok)

	// Give some time for background matching
	time.Sleep(100 * time.Millisecond)

	// 3. Verify matching results
	// Check Bid order status
	bidOrder, err := orderRepo.GetByID(uuid.MustParse(bidID))
	require.NoError(t, err)
	assert.Equal(t, domain.Filled, bidOrder.Status)
	assert.Equal(t, 0, bidOrder.RemainingQuantity)

	// Check Ask order status
	askOrder, err = orderRepo.GetByID(uuid.MustParse(askID))
	require.NoError(t, err)
	assert.Equal(t, domain.Filled, askOrder.Status)
	assert.Equal(t, 0, askOrder.RemainingQuantity)

	// Check trade exists
	trades, err := tradeRepo.GetByOrderID(uuid.MustParse(bidID))
	require.NoError(t, err)
	assert.Len(t, trades, 1)
	assert.True(t, trades[0].Price.Equal(askPrice))
}

func TestOrderNotFound(t *testing.T) {
	server, _, _, cleanup := setupTestServer()
	defer cleanup()

	nonExistentID := uuid.New().String()
	resp, err := http.Get(server.URL + "/orders/" + nonExistentID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestOrderNoMatch(t *testing.T) {
	server, orderRepo, _, cleanup := setupTestServer()
	defer cleanup()

	symbol := "AAPL"
	bidPrice := decimal.NewFromInt(100)
	askPrice := decimal.NewFromInt(110)

	// 1. Submit an Ask order
	askRequest := map[string]interface{}{
		"broker_id":   "broker1",
		"owner_doc":   "doc1",
		"type":        "ASK",
		"symbol":      symbol,
		"price":       askPrice.String(),
		"quantity":    10,
		"valid_until": time.Now().Add(time.Hour).Format(time.RFC3339),
	}
	body, _ := json.Marshal(askRequest)
	resp, err := http.Post(server.URL+"/orders", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var askResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&askResp)
	require.NoError(t, err)
	askID, ok := askResp["id"]
	require.True(t, ok)

	// 2. Submit a Bid order that doesn't match (price too low)
	bidRequest := map[string]interface{}{
		"broker_id":   "broker2",
		"owner_doc":   "doc2",
		"type":        "BID",
		"symbol":      symbol,
		"price":       bidPrice.String(),
		"quantity":    10,
		"valid_until": time.Now().Add(time.Hour).Format(time.RFC3339),
	}
	body, _ = json.Marshal(bidRequest)
	resp, err = http.Post(server.URL+"/orders", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var bidResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&bidResp)
	require.NoError(t, err)
	bidID, ok := bidResp["id"]
	require.True(t, ok)

	// Give some time for background matching
	time.Sleep(100 * time.Millisecond)

	// 3. Verify no matching occurred
	// Check Bid order status
	bidOrder, err := orderRepo.GetByID(uuid.MustParse(bidID))
	require.NoError(t, err)
	assert.Equal(t, domain.Pending, bidOrder.Status)
	assert.Equal(t, 10, bidOrder.RemainingQuantity)

	// Check Ask order status
	askOrder, err := orderRepo.GetByID(uuid.MustParse(askID))
	require.NoError(t, err)
	assert.Equal(t, domain.Pending, askOrder.Status)
	assert.Equal(t, 10, askOrder.RemainingQuantity)
}

func TestOrderPartialFill(t *testing.T) {
	server, orderRepo, tradeRepo, cleanup := setupTestServer()
	defer cleanup()

	symbol := "AAPL"
	bidPrice := decimal.NewFromInt(150)
	askPrice := decimal.NewFromInt(140)

	// 1. Submit an Ask order for 10
	askRequest := map[string]interface{}{
		"broker_id":   "broker1",
		"owner_doc":   "doc1",
		"type":        "ASK",
		"symbol":      symbol,
		"price":       askPrice.String(),
		"quantity":    10,
		"valid_until": time.Now().Add(time.Hour).Format(time.RFC3339),
	}
	body, _ := json.Marshal(askRequest)
	resp, err := http.Post(server.URL+"/orders", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var askResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&askResp)
	require.NoError(t, err)
	askID, ok := askResp["id"]
	require.True(t, ok)

	// 2. Submit a Bid order for 5
	bidRequest := map[string]interface{}{
		"broker_id":   "broker2",
		"owner_doc":   "doc2",
		"type":        "BID",
		"symbol":      symbol,
		"price":       bidPrice.String(),
		"quantity":    5,
		"valid_until": time.Now().Add(time.Hour).Format(time.RFC3339),
	}
	body, _ = json.Marshal(bidRequest)
	resp, err = http.Post(server.URL+"/orders", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var bidResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&bidResp)
	require.NoError(t, err)
	bidID, ok := bidResp["id"]
	require.True(t, ok)

	// Give some time for background matching
	time.Sleep(100 * time.Millisecond)

	// 3. Verify partial fill results
	// Check Bid order status (should be FILLED)
	bidOrder, err := orderRepo.GetByID(uuid.MustParse(bidID))
	require.NoError(t, err)
	assert.Equal(t, domain.Filled, bidOrder.Status)
	assert.Equal(t, 0, bidOrder.RemainingQuantity)

	// Check Ask order status (should be PARTIAL)
	askOrder, err := orderRepo.GetByID(uuid.MustParse(askID))
	require.NoError(t, err)
	assert.Equal(t, domain.Partial, askOrder.Status)
	assert.Equal(t, 5, askOrder.RemainingQuantity)

	// Check trade exists
	trades, err := tradeRepo.GetByOrderID(uuid.MustParse(bidID))
	require.NoError(t, err)
	assert.Len(t, trades, 1)
	assert.Equal(t, 5, trades[0].Quantity)
	assert.True(t, trades[0].Price.Equal(askPrice))
}
