package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"mini-stock-exchange/internal/handler"
	"mini-stock-exchange/internal/observability"
	"mini-stock-exchange/internal/repository"
	"mini-stock-exchange/internal/service"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	observability.InitLogger()
	ctx := context.Background()
	tp, err := observability.InitTracer(ctx, "stock-exchange-api")
	if err != nil {
		log.Fatalf("failed to init tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("error shutting down tracer provider: %v", err)
		}
	}()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/stockexchange?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	orderRepo := repository.NewOrderRepository(db)
	tradeRepo := repository.NewTradeRepository(db)
	orderService := service.NewOrderService(orderRepo, tradeRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	r := chi.NewRouter()
	orderHandler.RegisterRoutes(r)
	r.Handle("/metrics", promhttp.Handler())

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
