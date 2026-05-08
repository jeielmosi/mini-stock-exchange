package repository

import (
	"database/sql"
	"time"

	"mini-stock-exchange/internal/domain"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type postgresOrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) domain.OrderRepository {
	return &postgresOrderRepository{db: db}
}

func (r *postgresOrderRepository) Create(order *domain.Order) error {
	query := `INSERT INTO orders (id, broker_id, owner_doc, type, symbol, price, quantity, remaining_quantity, valid_until, status, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.Exec(query, order.ID, order.BrokerID, order.OwnerDoc, order.Type, order.Symbol, order.Price, order.Quantity, order.RemainingQuantity, order.ValidUntil, order.Status, order.CreatedAt)
	return err
}

func (r *postgresOrderRepository) GetByID(id uuid.UUID) (*domain.Order, error) {
	order := &domain.Order{}
	query := `SELECT id, broker_id, owner_doc, type, symbol, price, quantity, remaining_quantity, valid_until, status, created_at FROM orders WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&order.ID, &order.BrokerID, &order.OwnerDoc, &order.Type, &order.Symbol, &order.Price, &order.Quantity, &order.RemainingQuantity, &order.ValidUntil, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *postgresOrderRepository) Update(order *domain.Order) error {
	query := `UPDATE orders SET remaining_quantity = $1, status = $2 WHERE id = $3`
	_, err := r.db.Exec(query, order.RemainingQuantity, order.Status, order.ID)
	return err
}

// symbol string, orderType domain.OrderType, price decimal.Decimal, quantity int
// ([]*domain.Order, error)
func (r *postgresOrderRepository) FindMatches(order domain.Order) ([]domain.Order, error) {
	var query string
	var args []interface{}

	switch order.Type {
	case domain.Bid:
		query, args = createAskQuery(order)
	case domain.Ask:
		query, args = createBidQuery(order)
	}

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

func (r *postgresOrderRepository) UpdateRemainingQuantity(id uuid.UUID, quantity int, status domain.OrderStatus) error {
	query := `UPDATE orders SET remaining_quantity = $1, status = $2 WHERE id = $3`
	_, err := r.db.Exec(query, quantity, status, id)
	return err
}

func createAskQuery(order domain.Order) (query string, args []interface{}) {
	now := time.Now()
	query = `WITH ranked_asks AS (
			SELECT id, broker_id, owner_doc, type, symbol, price, quantity, remaining_quantity, valid_until, status, created_at,
				SUM(remaining_quantity) OVER (ORDER BY price ASC, created_at ASC) as running_total
			FROM orders 
			WHERE symbol = $1 AND type = 'ASK' AND status IN ('PENDING', 'PARTIAL') AND price <= $2 AND valid_until > $3
		)
		SELECT id, broker_id, owner_doc, type, symbol, price, quantity, remaining_quantity, valid_until, status, created_at 
			FROM ranked_asks o
			WHERE o.running_total - o.remaining_quantity < $4`
	args = []interface{}{order.Symbol, order.Price, now, order.Quantity}
	return
}

func createBidQuery(order domain.Order) (query string, args []interface{}) {
	now := time.Now()
	query = `WITH ranked_bids AS (
		SELECT id, broker_id, owner_doc, type, symbol, price, quantity, remaining_quantity, valid_until, status, created_at,
			SUM(remaining_quantity) OVER (ORDER BY price DESC, created_at ASC) as running_total
		FROM orders 
		WHERE symbol = $1 AND type = 'BID' AND status IN ('PENDING', 'PARTIAL') AND price >= $2 AND valid_until > $3
	)
	SELECT id, broker_id, owner_doc, type, symbol, price, quantity, remaining_quantity, valid_until, status, created_at 
		FROM ranked_bids o
		WHERE o.running_total - o.remaining_quantity < $4`
	args = []interface{}{order.Symbol, order.Price, now, order.Quantity}
	return
}
