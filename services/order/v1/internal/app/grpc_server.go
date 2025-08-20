package app

import (
	"net"
	"order/v1/proto/order"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func StartgRPCServer(orderServiceRPC order.OrderServiceServer) {
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// ✅ Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	listener, err := net.Listen("tcp", ":1028")
	if err != nil {
		panic(err)
	}

	order.RegisterOrderServiceServer(s, orderServiceRPC)

	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
