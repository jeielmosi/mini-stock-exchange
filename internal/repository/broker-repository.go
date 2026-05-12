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

type BrokerRepository interface {
	Stop() error
	Insert(broker entity.Broker) error
	GetByID(id uuid.UUID) (entity.Broker, error)
}

type brokerRepository struct {
	db *sql.DB
}

func NewBrokerRepository(db *sql.DB) (BrokerRepository, error) {
	var err error
	if db == nil {
		db, err = sql.Open("postgres", config.ENV.DatabaseURL)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to connect to db: %v", err)
		}
	}
	return &brokerRepository{db: db}, nil
}

func (r *brokerRepository) Stop() error {
	err := r.db.Close()
	slog.Error("BrokerRepository", "error", err)
	return err
}

func (r *brokerRepository) Insert(broker entity.Broker) error {
	query := `INSERT INTO brokers (id, name) VALUES ($1, $2)`
	_, err := r.db.Exec(query, broker.ID, broker.Name)
	return err
}
func (r *brokerRepository) GetByID(id uuid.UUID) (entity.Broker, error) {
	broker := entity.Broker{}
	query := `SELECT id, name FROM brokers WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&broker.ID, &broker.Name)
	if err != nil {
		return broker, err
	}
	return broker, nil
}
