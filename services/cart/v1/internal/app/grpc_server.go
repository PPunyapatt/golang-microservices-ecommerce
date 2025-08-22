package app

import (
	"cart/v1/proto/cart"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func StartgRPCServer(cartServerRPC cart.CartServiceServer) {
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// ✅ Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	listener, err := net.Listen("tcp", ":1027")
	if err != nil {
		panic(err)
	}

	cart.RegisterCartServiceServer(s, cartServerRPC)

	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
