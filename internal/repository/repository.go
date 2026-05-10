package repository

import (
	"database/sql"
	"fmt"
	"mini-stock-exchange/internal/config"
)

func NewDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", config.ENV.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}

	return db, nil
}
