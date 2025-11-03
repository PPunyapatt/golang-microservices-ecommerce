package server

import (
	"context"
	"log/slog"
	"order/v1/internal/app"
	"order/v1/internal/repository"
	"order/v1/internal/service"
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
	shutdown := tracer.InitTracer("order-service")
	defer func() { _ = shutdown(ctx) }()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	prometheusMetrics := metrics.NewMetrics()

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

	// RabbitMQ Connection
	rb, err := rabbitmq.NewRabbitMQConnection(ctx, s.cfg.RabbitMQUrl)
	if err != nil {
		slog.Error("NewRabbitMQConnection Error", "error", err.Error())
		return err
	}

	// ✅ Repository & Publisher
	orderRepo := repository.NewOrderRepository(db.Gorm, db.Sqlx)
	orderPublisher, err := publisher.NewPublisher(rb.Conn)
	if err != nil {
		slog.Error("OrderPublisher Error", "error", err.Error())
		return err
	}
	orderPublisher.Configure(
		publisher.TopicType("topic"),
	)
	orderService, orderServiceRPC := service.NewOrderServer(orderRepo, orderPublisher, otel.Tracer("inventory-service"), prometheusMetrics)
	newInitConsumer := app.NewInitConsumer(orderService, rb.Conn)
	newInitConsumer.InitConsumerWithReconnection(ctx)

	var wg sync.WaitGroup
	wg.Add(3)
	go rb.HandleGracefulShutdown(ctx, &wg)
	go prometheusMetrics.PrometheusHttp(ctx, &wg)
	go app.StartgRPCServer(ctx, orderServiceRPC, &wg, prometheusMetrics)
	wg.Wait()

	return nil
}
