package main

import (
	"context"
	"log"
	"net/http"

	"mini-stock-exchange/internal/config"
	"mini-stock-exchange/internal/controller"
	"mini-stock-exchange/internal/observability"
	"mini-stock-exchange/internal/repository"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
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

	config.Load()

	db, err := repository.NewPostgres(nil)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	orderRepo, err := repository.NewOrderRepository(db)
	if err != nil {
		log.Fatalf("failed to create order repository: %v", err)
	}
	defer orderRepo.Stop()

	tradeRepo, err := repository.NewTradeRepository(db)
	if err != nil {
		log.Fatalf("failed to create tarde repository: %v", err)
	}
	defer tradeRepo.Stop()

	r := chi.NewRouter()
	ctrl := controller.NewController(r, orderRepo, tradeRepo)
	ctrl.RegisterRoutes(r)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
