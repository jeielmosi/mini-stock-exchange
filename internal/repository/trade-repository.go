package repository

import (
	"database/sql"
	"fmt"

	"mini-stock-exchange/internal/config"
	"mini-stock-exchange/internal/entity"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type TradeRepository interface {
	Stop() error
	GetByID(id uuid.UUID) (entity.Trade, error)
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
	return r.db.Close()
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
