package server

import (
	"context"
	"gateway/v1/internal/api/handler"
	"gateway/v1/internal/helper"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"package/metrics"
	"sync"
	"syscall"

	"package/tracer"

	"github.com/joho/godotenv"
)

func Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Initialize structured JSON logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// Load environment variables from .env file (optional)
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, skipping...")
	}

	// Initialize middleware
	prometheusMetrics := metrics.NewMetrics()

	// Initialize OpenTelemetry tracer
	shutdown := tracer.InitTracer("gateway")
	defer func() { _ = shutdown(context.Background()) }()

	// Connect to all gRPC services
	conn := helper.NewClientsGRPC()
	log.Println("Connected to all gRPC server")

	// Initialize service handler
	service := handler.ServiceNew(conn, prometheusMetrics)

	var wg sync.WaitGroup
	wg.Add(2)
	// Start HTTP routes in a separate goroutine
	go MapRoutes(ctx, service, prometheusMetrics, &wg)

	// Start Prometheus endpoint
	go prometheusMetrics.PrometheusHttp(ctx, &wg)

	<-ctx.Done()
	wg.Wait()

	return nil
}
