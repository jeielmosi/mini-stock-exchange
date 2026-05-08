package observability

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var (
	OrdersSubmitted = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "orders_submitted_total",
		Help: "Total number of orders submitted",
	}, []string{"type", "symbol"})

	TradesExecuted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "trades_executed_total",
		Help: "Total number of trades executed",
	})

	MatchingLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "matching_engine_latency_seconds",
		Help:    "Latency of the matching engine",
		Buckets: prometheus.DefBuckets,
	})

	ActiveOrders = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_orders_count",
		Help: "Number of active orders",
	})
)

func InitLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func InitTracer(ctx context.Context, serviceName string) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tracerProvider)
	return tracerProvider, nil
}
