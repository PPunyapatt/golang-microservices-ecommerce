package main

import (
	"config-service"
	"context"
	"order/v1/internal/app"
	"order/v1/internal/repository"
	"order/v1/internal/service"
	database "package/Database"
	"package/rabbitmq"
	"package/rabbitmq/publisher"
	"package/tracer"

	"go.opentelemetry.io/otel"
)

func main() {
	cfg, err := config.SetUpEnv()
	if err != nil {
		panic(err)
	}

	// ✅ Init tracer
	shutdown := tracer.InitTracer("order-service")
	defer func() { _ = shutdown(context.Background()) }()

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

	// ✅ Repository & Publisher
	orderRepo := repository.NewOrderRepository(db.Gorm, db.Sqlx)
	orderPublisher, err := publisher.NewPublisher(conn)
	if err != nil {
		panic(err)
	}
	orderPublisher.Configure(
		publisher.TopicType("topic"),
	)

	orderService, orderServiceRPC := service.NewOrderServer(orderRepo, orderPublisher, otel.Tracer("inventory-service"))

	app.InitConsumer(orderService, conn)

	app.StartgRPCServer(orderServiceRPC)
}
