package server

import (
	"cart/v1/internal/app"
	"cart/v1/internal/repository"
	"cart/v1/internal/service"
	"context"
	"log/slog"
	"os"
	"os/signal"
	database "package/Database"
	"package/config"
	"package/metrics"
	"package/tracer"
	"sync"
	"syscall"

	"package/rabbitmq"

	"go.opentelemetry.io/otel"
)

type server struct {
	cfg *config.AppConfig
}

func NewServer(cfg *config.AppConfig) *server {
	return &server{cfg}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// ✅ Init tracer
	shutdown := tracer.InitTracer("cart-service")
	defer func() { _ = shutdown(ctx) }()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	promMetrics := metrics.NewMetrics()
	promMetrics.RegisterMetrics(
		promMetrics.Grpc.ErrorRequests,
		promMetrics.Grpc.SuccessRequests,
		promMetrics.Grpc.RequestsTotal,
		promMetrics.Grpc.RequestDuration,
	)

	db, err := database.InitDatabase(s.cfg)
	if err != nil {
		return err
	}

	// RabbitMQ Connection
	rb, err := rabbitmq.NewRabbitMQConnection(ctx, s.cfg.RabbitMQUrl)
	if err != nil {
		slog.Error("rabbitMQ", "error", err.Error())
		return err
	}

	cartRepo := repository.NewRepository(db.Mongo, logger)

	cartService, cartServerRPC := service.NewCartServer(cartRepo, otel.Tracer("cart-service"))

	newInitConsumer := app.NewInitConsumer(cartService, rb.Conn)
	newInitConsumer.InitConsumerWithReconnection(ctx)

	var wg sync.WaitGroup
	wg.Add(3)
	go rb.HandleGracefulShutdown(ctx, &wg)
	go promMetrics.PrometheusHttp(ctx, &wg)
	go app.StartgRPCServer(ctx, cartServerRPC, &wg, promMetrics)
	wg.Wait()

	return nil
}
