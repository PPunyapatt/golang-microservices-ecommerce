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
	shutdown := tracer.InitTracer("auth-service")
	defer func() { _ = shutdown(ctx) }()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	prometheusMetrics := metrics.NewMetrics()

	db, err := database.InitDatabase(s.cfg)
	if err != nil {
		return err
	}

	// RabbitMQ Connection
	conn, err := rabbitmq.NewRabbitMQConnection(s.cfg.RabbitMQUrl)
	if err != nil {
		panic(err)
	}

	cartRepo := repository.NewRepository(db.Mongo, logger)

	cartService, cartServerRPC := service.NewCartServer(cartRepo, otel.Tracer("cart-service"), logger)

	newInitConsumer := app.NewInitConsumer(cartService, conn)
	newInitConsumer.InitConsumerWithReconnection()

	app.StartgRPCServer(cartServerRPC)
}
