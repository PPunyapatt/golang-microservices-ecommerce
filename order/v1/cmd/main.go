package main

import (
	"net"
	"order/v1/config"
	"order/v1/internal/repository"
	"order/v1/internal/service"
	database "order/v1/pkg/Database"
	"order/v1/pkg/rabbitmq"
	"order/v1/pkg/rabbitmq/consumer"
	"order/v1/pkg/rabbitmq/publisher"
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
	dbConn, err := database.InitDatabase(cfg)
	if err != nil {
		panic(err)
	}

	sqlDB, err := dbConn.Gorm.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()
	defer dbConn.Sqlx.Close()

	// RabbitMQ Connection
	conn, err := rabbitmq.NewRabbitMQConnection(cfg.RabbitMQUrl)
	if err != nil {
		panic(err)
	}

	consumer := consumer.NewConsumer(conn)
	publisher := publisher.NewPublisher(conn)

	consumer.Configure()

	publisher.Configure()

	orderRepo := repository.NewOrderRepository(dbConn.Gorm, dbConn.Sqlx)
	s := grpc.NewServer()

	// ✅ Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	listener, err := net.Listen("tcp", ":1028")
	if err != nil {
		panic(err)
	}

	order.RegisterOrderServiceServer(s, service.NewOrderServer(orderRepo, publisher))

	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
