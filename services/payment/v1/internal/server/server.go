package server

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	database "package/Database"
	"package/config"
	"package/metrics"
	"package/rabbitmq"
	"package/rabbitmq/publisher"
	"package/tracer"
	"payment/v1/internal/app"
	"payment/v1/internal/repository"
	"payment/v1/internal/service"
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
	shutdown := tracer.InitTracer("payment-service")
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
		slog.Error("Failed to connect to RabbitMQ", "error", err)
		return err
	}

	// ✅ Repository & Publisher
	paymentPublisher, err := publisher.NewPublisher(rb.Conn)
	if err != nil {
		slog.Error("Failed to create payment publisher", "error", err)
		return err
	}
	paymentPublisher.Configure(publisher.TopicType("topic"))

	paymentRepo := repository.NewPaymentRepository(db.Gorm, db.Sqlx)
	paymentService, paymentServiceRPC := service.NewPaymentService(s.cfg.StripeKey, paymentRepo, paymentPublisher, otel.Tracer("inventory-service"), prometheusMetrics)

	newInitConsumer := app.NewInitConsumer(paymentService, rb.Conn)
	newInitConsumer.InitConsumerWithReconnection(ctx)

	var wg sync.WaitGroup
	wg.Add(3)
	go rb.HandleGracefulShutdown(ctx, &wg)
	go prometheusMetrics.PrometheusHttp(ctx, &wg)
	go app.StartgRPCServer(ctx, paymentServiceRPC, &wg, prometheusMetrics)
	wg.Wait()

	return nil
}
