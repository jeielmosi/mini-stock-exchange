package repository

import (
	"database/sql"

	"mini-stock-exchange/internal/domain"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type OrderRepository interface {
	Insert(order domain.Order) error
	GetByID(id uuid.UUID) (domain.Order, error)
	Update(order domain.Order) error
	GetBids(symbol string) ([]domain.Order, error)
	GetAsks(symbol string) ([]domain.Order, error)
}

type postgresOrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &postgresOrderRepository{db: db}
}

func (r *postgresOrderRepository) Insert(order domain.Order) error {
	query := `INSERT INTO orders (id, broker_id, owner_doc, type, symbol, price, quantity, remaining_quantity, valid_until, status, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.Exec(query, order.ID, order.BrokerID, order.OwnerDoc, order.Type, order.Symbol, order.Price, order.Quantity, order.RemainingQuantity, order.ValidUntil, order.Status, order.CreatedAt)
	return err
}

func (r *postgresOrderRepository) GetByID(id uuid.UUID) (domain.Order, error) {
	order := domain.Order{}
	query := `SELECT id, broker_id, owner_doc, type, symbol, price, quantity, remaining_quantity, valid_until, status, created_at FROM orders WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&order.ID, &order.BrokerID, &order.OwnerDoc, &order.Type, &order.Symbol, &order.Price, &order.Quantity, &order.RemainingQuantity, &order.ValidUntil, &order.Status, &order.CreatedAt)
	if err != nil {
		return order, err
	}
	return order, nil
}

func (r *postgresOrderRepository) Update(order domain.Order) error {
	query := `UPDATE orders SET remaining_quantity = $1, status = $2 WHERE id = $3`
	_, err := r.db.Exec(query, order.RemainingQuantity, order.Status, order.ID)
	return err
}

func (r *postgresOrderRepository) GetBids(symbol string) ([]domain.Order, error) {
	query := `SELECT id, broker_id, owner_doc, type, symbol, price,
		quantity, remaining_quantity, valid_until, status, created_at
	FROM orders 
	WHERE symbol = $1 AND type = 'BID' AND status IN ('PENDING', 'PARTIAL')
	`
	args := []interface{}{symbol}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		order := domain.Order{}
		err := rows.Scan(
			&order.ID, &order.BrokerID, &order.OwnerDoc,
			&order.Type, &order.Symbol, &order.Price,
			&order.Quantity, &order.RemainingQuantity, &order.ValidUntil,
			&order.Status, &order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *postgresOrderRepository) GetAsks(symbol string) ([]domain.Order, error) {
	query := `SELECT id, broker_id, owner_doc, type, symbol, price,
		quantity, remaining_quantity, valid_until, status, created_at
	FROM orders 
	WHERE symbol = $1 AND type = 'ASK' AND status IN ('PENDING', 'PARTIAL')
	`
	args := []interface{}{symbol}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		order := domain.Order{}
		err := rows.Scan(
			&order.ID, &order.BrokerID, &order.OwnerDoc,
			&order.Type, &order.Symbol, &order.Price,
			&order.Quantity, &order.RemainingQuantity, &order.ValidUntil,
			&order.Status, &order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
