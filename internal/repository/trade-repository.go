package repository

import (
	"database/sql"
	"fmt"
	"log/slog"

	"mini-stock-exchange/internal/config"
	"mini-stock-exchange/internal/entity"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type TradeRepository interface {
	Stop() error
	GetByID(id uuid.UUID) (entity.Trade, error)
	GetByOrderID(orderId uuid.UUID) ([]uuid.UUID, error)
}

type tradeRepository struct {
	db *sql.DB
}

func NewTradeRepository(db *sql.DB) (TradeRepository, error) {
	var err error
	if db == nil {
		db, err = sql.Open("postgres", config.ENV.DatabaseURL)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to connect to db: %v", err)
		}
	}
	return &tradeRepository{db: db}, nil
}

func (r *tradeRepository) Stop() error {
	err := r.db.Close()
	slog.Error("TradeRepository", "error", err)
	return err
}

func (r *tradeRepository) GetByID(id uuid.UUID) (entity.Trade, error) {
	trade := entity.Trade{}
	query := `SELECT id, symbol, price, quantity, executed_at, buy_order_id, sell_order_id FROM trades WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&trade.ID, &trade.Symbol, &trade.Price, &trade.Quantity, &trade.ExecutedAt, &trade.BuyOrderID, &trade.SellOrderID)
	if err != nil {
		return trade, err
	}
	return trade, nil
}

func (r *tradeRepository) GetByOrderID(orderId uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT id
	FROM trades 
	WHERE buy_order_id = $1 OR sell_order_id = $1`
	args := []interface{}{orderId}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tradeIds := []uuid.UUID{}
	for rows.Next() {
		tradeId := uuid.UUID{}
		err := rows.Scan(&tradeId)
		if err != nil {
			return nil, err
		}
		tradeIds = append(tradeIds, tradeId)
	}
	return tradeIds, nil
}
