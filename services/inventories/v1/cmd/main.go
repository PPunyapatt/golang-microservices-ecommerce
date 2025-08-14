package main

import (
	"config-service"
	"context"
	"inventories/v1/internal/repository"
	"inventories/v1/internal/services"
	"inventories/v1/proto/Inventory"
	"net"
	database "package/Database"
	"package/rabbitmq"
	"package/rabbitmq/consumer"
	"package/rabbitmq/publisher"
	"package/tracer"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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

	inventoryPublisher := publisher.NewPublisher(conn)
	inventoryPublisher.Configure(
		publisher.TopicType("topic"),
	)

	inventoryConsumer := consumer.NewConsumer(conn)
	inventoryConsumer.Configure(
		consumer.ExchangeName([]string{"order.excahnge"}),
		consumer.RoutingKeys([]string{"payment.*", "order.*"}),
		consumer.WorkerPoolSize(1),
		consumer.TopicType("topic"),
	)

	inventoryService := services.NewInventoryServer(inventoryRepo, inventoryPublisher, otel.Tracer("inventory-service"))

	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// ✅ Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	listener, err := net.Listen("tcp", ":1026")
	if err != nil {
		panic(err)
	}

	Inventory.RegisterInventoryServiceServer(s, inventoryService)

	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
