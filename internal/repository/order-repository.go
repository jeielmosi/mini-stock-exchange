package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"math/big"
	"mini-stock-exchange/internal/entity"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type MatchDTO struct {
	Ask   entity.Order
	Bid   entity.Order
	Trade entity.Trade
}

type OrderRepository interface {
	Stop() error
	Insert(order entity.Order) error
	GetByID(id uuid.UUID) (entity.Order, error)
	Match(ctx context.Context, match MatchDTO) error
	Expire(ids []uuid.UUID) error
	GetBids(symbol string, limit int) ([]entity.Order, error)
	GetAsks(symbol string, limit int) ([]entity.Order, error)
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) (OrderRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &orderRepository{db: db}, nil
}

func (r *orderRepository) Stop() error {
	err := r.db.Close()
	slog.Error("OrderRepository", "error", err)
	return err
}

func (r *orderRepository) Insert(order entity.Order) error {
	query := `INSERT INTO orders (id, broker_id, owner_doc, type, symbol, price, quantity, remaining_quantity, valid_until, status, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.Exec(query, order.ID, order.BrokerID, order.OwnerDoc, order.Type, order.Symbol, order.Price.FloatString(8), order.Quantity, order.RemainingQuantity, order.ValidUntil, order.Status, order.CreatedAt)
	return err
}

func (r *orderRepository) GetByID(id uuid.UUID) (entity.Order, error) {
	order := entity.Order{}
	var priceStr string
	query := `SELECT id, broker_id, owner_doc, type, symbol, price, quantity, remaining_quantity, valid_until, status, created_at FROM orders WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&order.ID, &order.BrokerID, &order.OwnerDoc, &order.Type, &order.Symbol, &priceStr, &order.Quantity, &order.RemainingQuantity, &order.ValidUntil, &order.Status, &order.CreatedAt)
	if err != nil {
		return order, err
	}

	price, ok := new(big.Rat).SetString(priceStr)
	if !ok {
		return order, fmt.Errorf("failed to parse price: %s", priceStr)
	}
	order.Price = price

	return order, nil
}

func (r *orderRepository) Match(ctx context.Context, match MatchDTO) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//Update Ask
	_, err = tx.ExecContext(ctx,
		`UPDATE orders SET remaining_quantity = $1, status = $2 WHERE id = $3 AND type = 'ASK' AND status IN ('PENDING', 'PARTIAL') AND symbol = $4`,
		match.Ask.RemainingQuantity, match.Ask.Status, match.Ask.ID, match.Ask.Symbol,
	)
	if err != nil {
		return err
	}

	//Update Bid
	_, err = tx.ExecContext(ctx,
		`UPDATE orders SET remaining_quantity = $1, status = $2 WHERE id = $3 AND type = 'BID' AND status IN ('PENDING', 'PARTIAL') AND symbol = $4`,
		match.Bid.RemainingQuantity, match.Bid.Status, match.Bid.ID, match.Ask.Symbol,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO trades (id, symbol, price, quantity, executed_at, buy_order_id, sell_order_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		match.Trade.ID, match.Trade.Symbol, match.Trade.Price.FloatString(8), match.Trade.Quantity,
		match.Trade.ExecutedAt, match.Trade.BuyOrderID, match.Trade.SellOrderID,
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) Expire(ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	query := `UPDATE orders
	SET status = 'EXPIRED'
	WHERE id = ANY($1::uuid[])
	AND status IN ('PENDING', 'PARTIAL', 'EXPIRED')
	`
	_, err := r.db.Exec(query, ids)

	return err
}

func (r *orderRepository) GetBids(symbol string, limit int) ([]entity.Order, error) {
	query := `SELECT id, broker_id, owner_doc, type, symbol, price,
		quantity, remaining_quantity, valid_until, status, created_at
	FROM orders 
	WHERE symbol = $1 AND type = 'BID' AND status IN ('PENDING', 'PARTIAL') AND NOW() < valid_until
	ORDER BY price DESC, created_at ASC
	LIMIT $2
	`
	args := []interface{}{symbol, limit}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToOrders(rows)
}

func (r *orderRepository) GetAsks(symbol string, limit int) ([]entity.Order, error) {
	query := `SELECT id, broker_id, owner_doc, type, symbol, price,
		quantity, remaining_quantity, valid_until, status, created_at
	FROM orders 
	WHERE symbol = $1 AND type = 'ASK' AND status IN ('PENDING', 'PARTIAL') AND NOW() < valid_until
	ORDER BY price ASC, created_at ASC
	LIMIT $2
	`
	args := []interface{}{symbol, limit}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToOrders(rows)
}

func rowsToOrders(rows *sql.Rows) ([]entity.Order, error) {
	var orders []entity.Order
	for rows.Next() {
		order := entity.Order{}
		var priceStr string
		err := rows.Scan(
			&order.ID, &order.BrokerID, &order.OwnerDoc,
			&order.Type, &order.Symbol, &priceStr,
			&order.Quantity, &order.RemainingQuantity, &order.ValidUntil,
			&order.Status, &order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		price, ok := new(big.Rat).SetString(priceStr)
		if !ok {
			return nil, fmt.Errorf("failed to parse price: %s", priceStr)
		}
		order.Price = price

		orders = append(orders, order)
	}
	return orders, nil
}
