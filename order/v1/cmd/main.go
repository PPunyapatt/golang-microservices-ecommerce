package main

import (
	"net"
	"order/v1/config"
	"order/v1/internal/repository"
	"order/v1/internal/service"
	database "order/v1/pkg/Database"
	"order/v1/proto/order"

	"github.com/streadway/amqp"
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
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	orderRepo := repository.NewOrderRepository(dbConn.Gorm, dbConn.Sqlx)
	s := grpc.NewServer()

	// ✅ Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	listener, err := net.Listen("tcp", ":1028")
	if err != nil {
		panic(err)
	}

	order.RegisterOrderServiceServer(s, service.NewOrderServer(orderRepo))

	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
