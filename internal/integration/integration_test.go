package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mini-stock-exchange/internal/config"
	"mini-stock-exchange/internal/controller"
	"mini-stock-exchange/internal/dto"
	dto_helper "mini-stock-exchange/internal/dto/helper"
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer() (*httptest.Server, repository.OrderRepository, func()) {
	config.LoadTest(10)
	ctx := context.Background()
	db, cleanup, err := repository.SetupTestDB(ctx)
	if err != nil {
		panic(err)
	}

	orderRepo, err := repository.NewOrderRepository(db)
	if err != nil {
		panic(err)
	}
	tradeRepo, err := repository.NewTradeRepository(db)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	ctrl := controller.NewMockController(r, orderRepo, tradeRepo)
	ctrl.RegisterRoutes(r)

	server := httptest.NewServer(r)

	return server, orderRepo, func() {
		server.Close()
		cleanup()
	}
}

func TestOrderFlow(t *testing.T) {
	server, orderRepo, cleanup := setupTestServer()
	defer cleanup()

	symbol := "AAPL"
	bidPrice := float64(150)
	askPrice := float64(140)

	// 1. Submit an Ask order
	askRequest := map[string]interface{}{
		"broker_id":   "broker1",
		"owner_doc":   "doc1",
		"type":        entity.Ask,
		"symbol":      symbol,
		"price":       askPrice,
		"quantity":    10,
		"valid_until": time.Now().Format(time.DateOnly),
	}

	body, err := json.Marshal(askRequest)
	require.NoError(t, err)
	resp, err := http.Post(server.URL+"/orders", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Log(err.Error())
	}
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var respDTO dto.CreateOrderResponse
	err = json.NewDecoder(resp.Body).Decode(&respDTO)
	assert.NoError(t, err)
	assert.NotEmpty(t, respDTO.ID)

	askID, err := dto_helper.DecodeUUID(respDTO.ID)
	require.NoError(t, err)
	askIDStr := askID.String()
	require.NotEmpty(t, askIDStr)

	// Verify Ask order is in DB
	askOrder, err := orderRepo.GetByID(uuid.MustParse(askIDStr))
	require.NoError(t, err)
	assert.Equal(t, symbol, askOrder.Symbol)
	assert.Equal(t, entity.Ask, askOrder.Type)
	assert.Equal(t, 10, askOrder.Quantity)

	// 2. Submit a matching Bid order
	bidRequest := map[string]interface{}{
		"broker_id":   "broker2",
		"owner_doc":   "doc2",
		"type":        "BID",
		"symbol":      symbol,
		"price":       bidPrice,
		"quantity":    10,
		"valid_until": time.Now().Format(time.DateOnly),
	}

	body, err = json.Marshal(bidRequest)
	require.NoError(t, err)
	resp, err = http.Post(server.URL+"/orders", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var bidResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&bidResp)
	require.NoError(t, err)
	bidIdEncoded, ok := bidResp["id"]
	require.True(t, ok)

	bidID, err := dto_helper.DecodeUUID(bidIdEncoded)
	require.NoError(t, err)
	bidIdStr := bidID.String()
	require.NotEmpty(t, bidIdStr)

	// Give some time for background matching
	time.Sleep(100 * time.Millisecond)

	// 3. Verify matching results
	// Check Bid order status
	bidOrder, err := orderRepo.GetByID(uuid.MustParse(bidIdStr))
	require.NoError(t, err)
	assert.Equal(t, entity.Filled, bidOrder.Status)
	assert.Equal(t, 0, bidOrder.RemainingQuantity)

	// Check Ask order status
	askOrder, err = orderRepo.GetByID(uuid.MustParse(askIDStr))
	require.NoError(t, err)
	assert.Equal(t, entity.Filled, askOrder.Status)
	assert.Equal(t, 0, askOrder.RemainingQuantity)
}

func TestOrderNotFound(t *testing.T) {
	server, _, cleanup := setupTestServer()
	defer cleanup()

	id, err := uuid.NewV7()
	require.NoError(t, err)
	nonExistentID, err := dto_helper.EncodeUUID(id)
	require.NoError(t, err)
	resp, err := http.Get(server.URL + "/orders/" + nonExistentID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestOrderNoMatch(t *testing.T) {
	server, orderRepo, cleanup := setupTestServer()
	defer cleanup()

	symbol := "AAPL"
	bidPrice := float64(100)
	askPrice := float64(110)

	// 1. Submit an Ask order
	askRequest := map[string]interface{}{
		"broker_id":   "broker1",
		"owner_doc":   "doc1",
		"type":        "ASK",
		"symbol":      symbol,
		"price":       askPrice,
		"quantity":    10,
		"valid_until": time.Now().Format(time.DateOnly),
	}
	t.Log(askRequest)
	body, err := json.Marshal(askRequest)
	require.NoError(t, err)
	resp, err := http.Post(server.URL+"/orders", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var askResp dto.CreateOrderResponse
	err = json.NewDecoder(resp.Body).Decode(&askResp)
	require.NoError(t, err)
	assert.NotEmpty(t, askResp.ID)
	askID, err := dto_helper.DecodeUUID(askResp.ID)
	require.NoError(t, err)
	askIdStr := askID.String()
	require.NotEmpty(t, askIdStr)

	// 2. Submit a Bid order that doesn't match (price too low)
	bidRequest := map[string]interface{}{
		"broker_id":   "broker2",
		"owner_doc":   "doc2",
		"type":        "BID",
		"symbol":      symbol,
		"price":       bidPrice,
		"quantity":    10,
		"valid_until": time.Now().Format(time.DateOnly),
	}
	t.Log(bidPrice)
	body, err = json.Marshal(bidRequest)
	require.NoError(t, err)
	resp, err = http.Post(server.URL+"/orders", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var bidResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&bidResp)
	require.NoError(t, err)
	bidIdEncoded, ok := bidResp["id"]
	require.True(t, ok)

	bidID, err := dto_helper.DecodeUUID(bidIdEncoded)
	require.NoError(t, err)
	bidIdStr := bidID.String()
	require.NotEmpty(t, bidIdStr)

	// Give some time for background matching
	time.Sleep(100 * time.Millisecond)

	// 3. Verify no matching occurred
	// Check Bid order status
	bidOrder, err := orderRepo.GetByID(uuid.MustParse(bidIdStr))
	require.NoError(t, err)
	assert.Equal(t, entity.Pending, bidOrder.Status)
	assert.Equal(t, 10, bidOrder.RemainingQuantity)

	// Check Ask order status
	askOrder, err := orderRepo.GetByID(uuid.MustParse(askIdStr))
	require.NoError(t, err)
	assert.Equal(t, entity.Pending, askOrder.Status)
	assert.Equal(t, 10, askOrder.RemainingQuantity)
}

func TestOrderPartialFill(t *testing.T) {
	server, orderRepo, cleanup := setupTestServer()
	defer cleanup()
	config.LoadTest(10)

	symbol := "AAPL"
	bidPrice := float64(150)
	askPrice := float64(140)

	// 1. Submit an Ask order
	askRequest := map[string]interface{}{
		"broker_id":   "broker1",
		"owner_doc":   "doc1",
		"type":        "ASK",
		"symbol":      symbol,
		"price":       askPrice,
		"quantity":    10,
		"valid_until": time.Now().Format(time.DateOnly),
	}
	body, err := json.Marshal(askRequest)
	require.NoError(t, err)
	resp, err := http.Post(server.URL+"/orders", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var askResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&askResp)
	require.NoError(t, err)
	askIdEncoded, ok := askResp["id"]
	require.True(t, ok)
	askId, err := dto_helper.DecodeUUID(askIdEncoded)
	require.NoError(t, err)
	askIdStr := askId.String()
	require.NotEmpty(t, askIdStr)

	// 2. Submit a Bid order
	bidRequest := map[string]interface{}{
		"broker_id":   "broker2",
		"owner_doc":   "doc2",
		"type":        "BID",
		"symbol":      symbol,
		"price":       bidPrice,
		"quantity":    5,
		"valid_until": time.Now().Format(time.DateOnly),
	}

	body, err = json.Marshal(bidRequest)
	require.NoError(t, err)
	resp, err = http.Post(server.URL+"/orders", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var bidResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&bidResp)
	require.NoError(t, err)
	bidIdEncoded, ok := bidResp["id"]
	require.True(t, ok)
	bidId, err := dto_helper.DecodeUUID(bidIdEncoded)
	require.NoError(t, err)
	bidIdStr := bidId.String()
	require.NotEmpty(t, bidIdStr)

	// Give some time for background matching
	time.Sleep(100 * time.Millisecond)

	// 3. Verify partial fill results
	// Check Bid order status (should be FILLED)
	bidOrder, err := orderRepo.GetByID(uuid.MustParse(bidIdStr))
	require.NoError(t, err)
	assert.Equal(t, entity.Filled, bidOrder.Status)
	assert.Equal(t, 0, bidOrder.RemainingQuantity)

	// Check Ask order status (should be PARTIAL)
	askOrder, err := orderRepo.GetByID(uuid.MustParse(askIdStr))
	require.NoError(t, err)
	assert.Equal(t, entity.Partial, askOrder.Status)
	assert.Equal(t, 5, askOrder.RemainingQuantity)
}
