package main

import (
	"config-service"
	"net"
	"order/v1/internal/app"
	"order/v1/internal/repository"
	"order/v1/internal/service"
	database "package/Database"
	"package/rabbitmq"
	"package/rabbitmq/consumer"
	"package/rabbitmq/publisher"

	"order/v1/proto/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	cfg, err := config.SetUpEnv()
	if err != nil {
		panic(err)
	}

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

	publisher := publisher.NewPublisher(conn)

	orderRepo := repository.NewOrderRepository(db.Gorm, db.Sqlx)
	orderService := service.NewOrderServer(orderRepo, publisher)

	consumer := consumer.NewConsumer(conn)

	app := app.NewAppServer(orderService)

	go consumer.StartConsumer(app.Worker)

	s := grpc.NewServer()

	// ✅ Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	listener, err := net.Listen("tcp", ":1028")
	if err != nil {
		panic(err)
	}

	order.RegisterOrderServiceServer(s, orderService)

	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
