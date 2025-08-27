package main

import (
	"config-service"
	"context"
	"inventories/v1/internal/app"
	"inventories/v1/internal/repository"
	"inventories/v1/internal/services"
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
	shutdown := tracer.InitTracer("inventory-service")
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

	inventoryRepo := repository.NewInventoryRepository(db.Gorm, db.Sqlx)

	// RabbitMQ Connection
	conn, err := rabbitmq.NewRabbitMQConnection(cfg.RabbitMQUrl)
	if err != nil {
		panic(err)
	}

	inventoryPublisher, err := publisher.NewPublisher(conn)
	if err != nil {
		panic(err)
	}

	inventoryPublisher.Configure(publisher.TopicType("topic"))

	inventoryServiceRPC, inventoryService := services.NewInventoryServer(inventoryRepo, inventoryPublisher, otel.Tracer("inventory-service"))

	newInitConsumer := app.NewInitConsumer(inventoryService, conn)
	newInitConsumer.InitConsumerWithReconnection()

	app.StartgRPCServer(inventoryServiceRPC)
}
