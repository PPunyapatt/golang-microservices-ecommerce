package main

import (
	"context"
	database "package/Database"
	"package/config"
	"package/rabbitmq"
	"package/rabbitmq/publisher"
	"package/tracer"
	"payment/v1/internal/app"
	"payment/v1/internal/repository"
	"payment/v1/internal/service"

	"go.opentelemetry.io/otel"
)

func main() {
	cfg, err := config.SetUpEnv()
	if err != nil {
		panic(err)
	}

	shutdown := tracer.InitTracer("payment-service")
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
	paymentPublisher, err := publisher.NewPublisher(conn)
	if err != nil {
		panic(err)
	}
	paymentPublisher.Configure(publisher.TopicType("topic"))

	paymentRepo := repository.NewPaymentRepository(db.Gorm, db.Sqlx)
	paymentService, paymentServiceRPC := service.NewPaymentService(cfg.StripeKey, paymentRepo, paymentPublisher, otel.Tracer("inventory-service"))

	newInitConsumer := app.NewInitConsumer(paymentService, conn)
	newInitConsumer.InitConsumerWithReconnection()

	app.StartgRPCServer(paymentServiceRPC)
}
