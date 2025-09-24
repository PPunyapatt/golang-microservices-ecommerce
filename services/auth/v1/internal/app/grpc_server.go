package app

import (
	"auth-service/v1/proto/auth"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func StartgRPCServer(authService auth.AuthServiceServer) {
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// ✅ Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	listener, err := net.Listen("tcp", ":1024")
	if err != nil {
		panic(err)
	}

	auth.RegisterAuthServiceServer(s, authService)

	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
