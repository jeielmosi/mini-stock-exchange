package repository

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupTestDB(ctx context.Context) (*sql.DB, func()) {
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
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open db: %s", err)
	}

	// Apply schema
	_, filename, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(filename), "migrations")
	schemaPath := filepath.Join(migrationsDir, "000001_init_schema.up.sql")

	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("failed to read schema at %s: %s", schemaPath, err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		log.Fatalf("failed to apply schema: %s", err)
	}

	return db, func() {
		db.Close()
		pgContainer.Terminate(ctx)
	}
}