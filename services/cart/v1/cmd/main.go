package main

import (
	"cart/v1/config"
	"cart/v1/internal/repository"
	"cart/v1/internal/service"
	"cart/v1/pkg/Database/gorm"
	"cart/v1/pkg/Database/postgres"
	"cart/v1/proto/cart"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// s := grpc.NewServer()

	// listener, err := net.Listen("tcp", ":50051")
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }

	cfg, err := config.SetUpEnv()
	if err != nil {
		panic(err)
	}

	// database connection
	// dbGorm, err := gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{})
	// if err != nil {
	// 	panic(err)
	// }

	// dbSqlx, err := sqlx.Connect("pgx", cfg.Dsn)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	postgresDB, err := postgres.NewPostgresDB(cfg.Dsn)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	gormDB, err := gorm.NewGormConnection(cfg.Dsn)
	if err != nil {
		log.Fatalf("failed to connect to gorm: %v", err)
	}

	cartRepo := repository.NewRepository(postgresDB, gormDB)
	s := grpc.NewServer()

	// ✅ Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	listener, err := net.Listen("tcp", ":1027")
	if err != nil {
		panic(err)
	}

	cart.RegisterCartServiceServer(s, service.NewCartServer(cartRepo))

	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
