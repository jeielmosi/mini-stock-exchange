package repository

import (
	"database/sql"

	"mini-stock-exchange/internal/domain"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type TradeRepository interface {
	Create(trade domain.Trade) error
	GetByOrderID(orderID uuid.UUID) ([]domain.Trade, error)
}

type postgresTradeRepository struct {
	db *sql.DB
}

func NewTradeRepository(db *sql.DB) TradeRepository {
	return &postgresTradeRepository{db: db}
}

// TODO IT SHOULD BE  A TRANSACTION, not just a single request
func (r *postgresTradeRepository) Create(trade domain.Trade) error {
	query := `INSERT INTO trades (id, buy_order_id, sell_order_id, symbol, price, quantity, executed_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, trade.ID, trade.BuyOrderID, trade.SellOrderID, trade.Symbol, trade.Price, trade.Quantity, trade.ExecutedAt)
	return err
}

func (r *postgresTradeRepository) GetByOrderID(orderID uuid.UUID) ([]domain.Trade, error) {
	query := `SELECT id, buy_order_id, sell_order_id, symbol, price, quantity, executed_at FROM trades WHERE buy_order_id = $1 OR sell_order_id = $1`
	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []domain.Trade
	for rows.Next() {
		trade := domain.Trade{}
		if err := rows.Scan(&trade.ID, &trade.BuyOrderID, &trade.SellOrderID, &trade.Symbol, &trade.Price, &trade.Quantity, &trade.ExecutedAt); err != nil {
			return nil, err
		}
		trades = append(trades, trade)
	}
	return trades, nil
}
