package server

import (
	"context"
	"inventories/v1/internal/app"
	"inventories/v1/internal/repository"
	"inventories/v1/internal/services"
	"log/slog"
	"os"
	"os/signal"
	database "package/Database"
	"package/config"
	"package/metrics"
	"package/rabbitmq"
	"package/rabbitmq/publisher"
	"package/tracer"
	"sync"
	"syscall"

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
	shutdown := tracer.InitTracer("inventory-service")
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

	// database connection
	db, err := database.InitDatabase(s.cfg)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.Gorm.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()
	defer db.Sqlx.Close()

	inventoryRepo := repository.NewInventoryRepository(db.Gorm, db.Sqlx)

	// RabbitMQ Connection
	rb, err := rabbitmq.NewRabbitMQConnection(ctx, s.cfg.RabbitMQUrl)
	if err != nil {
		slog.Error("NewRabbitMQConnection Error", "error", err.Error())
		return err
	}

	inventoryPublisher, err := publisher.NewPublisher(rb.Conn)
	if err != nil {
		slog.Error("InventoryPublisher Error", "error", err.Error())
		return err
	}

	inventoryPublisher.Configure(publisher.TopicType("topic"))

	inventoryServiceRPC, inventoryService := services.NewInventoryServer(inventoryRepo, inventoryPublisher, otel.Tracer("inventory-service"))

	newInitConsumer := app.NewInitConsumer(inventoryService, rb.Conn)
	newInitConsumer.InitConsumerWithReconnection(ctx)

	var wg sync.WaitGroup
	wg.Add(3)
	go rb.HandleGracefulShutdown(ctx, &wg)
	go promMetrics.PrometheusHttp(ctx, &wg)
	go app.StartgRPCServer(ctx, inventoryServiceRPC, &wg, promMetrics)
	wg.Wait()

	return nil
}
