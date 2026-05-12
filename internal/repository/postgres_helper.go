package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"mini-stock-exchange/internal/config"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewMockPostgres(ctx context.Context) (*sql.DB, func() error, error) {
	pgContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
		),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			pgContainer.Terminate(ctx)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
		return nil, nil, err
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open db: %s", err)
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			db.Close()
		}
	}()

	// TODO Apply all schemas
	_, filename, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(filename), "migrations")
	schemaPath := filepath.Join(migrationsDir, "000001_init_schema.up.sql")

	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("failed to read schema at %s: %s", schemaPath, err)
		db.Close()
		return nil, nil, err
	}

	if _, err := db.Exec(string(schema)); err != nil {
		log.Fatalf("failed to apply schema: %s", err)
		return nil, nil, err
	}

	return db, func() error {
		err := db.Close()
		if err != nil {
			return errors.New("failed to close db: " + err.Error())
		}
		return pgContainer.Terminate(ctx)
	}, nil
}

func NewPostgres(db *sql.DB) (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	var err error
	db, err = sql.Open("postgres", config.ENV.DatabaseURL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}
	return db, nil
}
