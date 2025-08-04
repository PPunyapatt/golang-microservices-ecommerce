package main

import (
	"auth-service/v1/config"
	"auth-service/v1/internal/helper"
	"auth-service/v1/internal/repository"
	"auth-service/v1/internal/service"
	"auth-service/v1/pkg/database"
	"auth-service/v1/proto/auth"
	"context"
	"log"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
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
	shutdown := helper.InitTracer("auth-service")
	defer func() { _ = shutdown(context.Background()) }()

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

	authRepo := repository.NewAuthRepository(dbConn.Gorm, dbConn.Sqlx)

	// s := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	listener, err := net.Listen("tcp", ":1024")
	if err != nil {
		panic(err)
	}

	auth.RegisterAuthServiceServer(s, service.NewAuthServer(authRepo))

	// ✅ Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}

	log.Println("Auth service is running on port 1024")
}
