package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

func NewDatabase() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("no DATABASE_URL found")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}

	return db, nil
}
