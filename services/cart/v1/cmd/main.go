package main

import (
	"cart/v1/internal/app"
	"cart/v1/internal/repository"
	"cart/v1/internal/service"
	"config-service"
	"context"
	"log/slog"
	"os"
	database "package/Database"
	"package/rabbitmq"
	"package/tracer"

	"go.opentelemetry.io/otel"
)

func main() {
	cfg, err := config.SetUpEnv()
	if err != nil {
		panic(err)
	}

	// ✅ Init tracer
	shutdown := tracer.InitTracer("cart-service")
	defer func() { _ = shutdown(context.Background()) }()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // 👈 show file and line
	}))

	// database connection
	db, err := database.InitDatabase(cfg)
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
	conn, err := rabbitmq.NewRabbitMQConnection(cfg.RabbitMQUrl)
	if err != nil {
		panic(err)
	}

	cartRepo := repository.NewRepository(db.Sqlx, db.Gorm, db.Mongo, logger)

	cartService, cartServerRPC := service.NewCartServer(cartRepo, otel.Tracer("cart-service"), logger)

	newInitConsumer := app.NewInitConsumer(cartService, conn)
	newInitConsumer.InitConsumerWithReconnection()

	app.StartgRPCServer(cartServerRPC)
}
